package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	TEMP_DIR_PATH = "/tmp/hurlui"
)

type HurlResult struct {
	OutputString string     `json:"outputString"`
	Report       HurlReport `json:"report,omitempty"`
}

type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

type FileExplorerState struct {
	CurrentDir   FileInfo `json:"currentDir"`
	SelectedFile FileInfo `json:"selectedFile"`
}

type EnvConfig struct {
	Global       map[string]string            `json:"global"`
	Environments map[string]map[string]string `json:"environments"`
}

type ReturnValue struct {
	FileContent  string            `json:"fileContent,omitempty"`
	FileExplorer FileExplorerState `json:"fileExplorer"`
	Files        []FileInfo        `json:"files"`
	Error        string            `json:"error,omitempty"`
	HurlReport   HurlReport        `json:"hurlReport,omitempty"`
	Envs         []string          `json:"envs,omitempty"`
	EnvFilePath  string            `json:"envFilePath,omitempty"`
}

type App struct {
	ctx           context.Context
	explorerState FileExplorerState
	cacheDB       *bolt.DB
}

func NewApp() *App {
	homeDir, _ := os.UserHomeDir()
	app := &App{
		explorerState: FileExplorerState{
			CurrentDir:   FileInfo{Name: "Home", Path: homeDir, IsDir: true},
			SelectedFile: FileInfo{},
		},
	}

	if err := app.initCache(); err != nil {
		fmt.Printf("Failed to initialize cache: %v\n", err)
	}

	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if err := a.CloseCache(); err != nil {
		fmt.Printf("Failed to close cache: %v\n", err)
	}
}

func (a *App) initCache() error {
	cacheDir := filepath.Join(os.TempDir(), "hurlui")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	dbPath := filepath.Join(cacheDir, "cache.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open cache database: %w", err)
	}

	a.cacheDB = db

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("hurl_cache"))
		if err != nil {
			return fmt.Errorf("failed to create hurl_cache bucket: %w", err)
		}
		return nil
	})
}

func createFileInfo(path string) (FileInfo, error) {

	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Name:     filepath.Base(path),
		Path:     path,
		IsDir:    info.IsDir(),
		Size:     info.Size(),
		Modified: info.ModTime().Format("2006-01-02 15:04:05"),
	}, nil
}

// tempOutputPathFor returns a path inside TEMP_DIR_PATH mirroring the given file path.
// It avoids filepath.Join swallowing TEMP_DIR_PATH when filePath is absolute by
// stripping leading separators and joining the remainder under the temp root.
func tempOutputPathFor(filePath string) string {
	cleaned := filepath.Clean(filePath)
	// Make absolute for consistency if possible
	if abs, err := filepath.Abs(cleaned); err == nil {
		cleaned = abs
	}
	// Remove volume (Windows) and leading separators to make it relative
	// Keep the rest of the path hierarchy under TEMP_DIR_PATH
	vol := filepath.VolumeName(cleaned)
	rel := strings.TrimPrefix(cleaned, vol)
	for strings.HasPrefix(rel, string(os.PathSeparator)) {
		rel = strings.TrimPrefix(rel, string(os.PathSeparator))
	}
	return filepath.Join(TEMP_DIR_PATH, rel)
}

func (a *App) selectedFileOutputPath() string {
	return tempOutputPathFor(a.explorerState.SelectedFile.Path)
	// return filepath.Join(TEMP_DIR_PATH, a.explorerState.SelectedFile.Path)
}

func (a *App) selectedFileReportPath() string {
	return filepath.Join(a.selectedFileOutputPath(), "report.json")

}

func (a *App) selectedFileStorePath() string {
	return filepath.Join(a.selectedFileOutputPath(), "store")
}

func (a *App) insertResponseData(h *HurlReport, filePath string) error {
	fmt.Println("Inserting response data from files in:", h)
	outputDir := a.selectedFileOutputPath()

	for i := range *h {
		session := &(*h)[i]
		for j := range session.Entries {
			entry := &session.Entries[j]
			for k := range entry.Calls {
				call := &entry.Calls[k]

				// Check if response body is a file reference
				if call.Response.BodyPath != "" {
					bodyFilePath := filepath.Join(outputDir, call.Response.BodyPath)

					// Read the response body file
					if bodyContent, err := os.ReadFile(bodyFilePath); err == nil {
						// Replace file path with actual content
						call.Response.Body = string(bodyContent)
					} else {
						fmt.Printf("Failed to read response body file %s: %v\n", bodyFilePath, err)
						// Keep the original file path if reading fails
					}
				}
			}
		}
	}

	return nil
}

func (a *App) getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "hurlui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func (a *App) getEnvFilePath() (string, error) {
	configDir, err := a.getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "env.json"), nil
}

func (a *App) loadEnvConfig() (*EnvConfig, error) {

	envConfigPath, err := a.getEnvFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		return &EnvConfig{
			Global:       make(map[string]string),
			Environments: make(map[string]map[string]string),
		}, nil
	}

	data, err := os.ReadFile(envConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read env config file: %w", err)
	}

	var config EnvConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse env config: %w", err)
	}

	if config.Global == nil {
		config.Global = make(map[string]string)
	}
	if config.Environments == nil {
		config.Environments = make(map[string]map[string]string)
	}

	return &config, nil
}

// func (a *App) GetAvailableEnvGroups() ReturnValue {
// 	config, err := a.loadEnvConfig()
// 	if err != nil {
// 		return ReturnValue{Error: err.Error()}
// 	}

// 	groups := make([]string, 0, len(config.Selectable))
// 	for groupName := range config.Selectable {
// 		groups = append(groups, groupName)
// 	}

// 	return ReturnValue{Files: []FileInfo{{Name: "env_groups"}}, EnvVars: map[string]string{"groups": fmt.Sprintf("%v", groups)}}
// }

// func (a *App) GetFileContentAndExecuteHurl(filePath string) (map[string]string, error) {
// 	content, err := a.GetFileContent(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read file: %w", err)
// 	}

// 	output, err := a.ExecuteHurl(filePath)
// 	result := map[string]string{
// 		"content": content,
// 		"output":  output,
// 	}

// 	if err != nil {
// 		result["error"] = err.Error()
// 	}

// 	return result, nil
// }
