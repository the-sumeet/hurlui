package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	Global     map[string]string            `json:"global"`
	Selectable map[string]map[string]string `json:"selectable"`
}

type ReturnValue struct {
	FileContent  string            `json:"fileContent,omitempty"`
	FileExplorer FileExplorerState `json:"fileExplorer"`
	Files        []FileInfo        `json:"files"`
	Error        string            `json:"error,omitempty"`
	HurlReport   HurlReport        `json:"hurlReport,omitempty"`
	EnvVars      map[string]string `json:"envVars,omitempty"`
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

func (a *App) GetCurrentDirectory() FileInfo {
	return a.explorerState.CurrentDir
}

func (a *App) GetFiles() ReturnValue {
	entries, err := os.ReadDir(a.explorerState.CurrentDir.Path)
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Name:     entry.Name(),
			Path:     filepath.Join(a.explorerState.CurrentDir.Path, entry.Name()),
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
			Modified: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		files = append(files, fileInfo)
	}

	return ReturnValue{Files: files, FileExplorer: a.explorerState}
}

func (a *App) SetCurrentFile(ctx context.Context, file FileInfo) {
	a.explorerState.SelectedFile = file
	runtime.WindowSetTitle(ctx, file.Path)
}

func (a *App) ChangeDirectory(path string) ReturnValue {
	fileInfo, err := createFileInfo(path)
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	if !fileInfo.IsDir {
		fmt.Println("Failed to change directory:", path)
		return ReturnValue{Error: fmt.Sprintf("path is not a directory: %s", path)}
	}

	a.explorerState.CurrentDir = fileInfo
	return ReturnValue{FileExplorer: a.explorerState}
}

func (a *App) NavigateUp() ReturnValue {
	parent := filepath.Dir(a.explorerState.CurrentDir.Path)
	if parent == a.explorerState.CurrentDir.Path {
		return ReturnValue{Error: "already at root directory"}
	}
	return a.ChangeDirectory(parent)
}

func (a *App) SelectFile(filePath string) ReturnValue {
	fileInfo, err := createFileInfo(filePath)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("file does not exist: %s", filePath)}
	}

	a.SetCurrentFile(a.ctx, fileInfo)
	// a.explorerState.SelectedFile = fileInfo
	return ReturnValue{
		FileExplorer: a.explorerState,
	}
}

func (a *App) GetSelectedFile() FileInfo {
	return a.explorerState.SelectedFile
}

func (a *App) ClearSelection() {
	a.SetCurrentFile(a.ctx, FileInfo{})
	// a.explorerState.SelectedFile = FileInfo{}
}

func (a *App) GetExplorerState() ReturnValue {
	return ReturnValue{
		FileExplorer: a.explorerState,
	}
}

func (a *App) GetFileContent(filePath string) ReturnValue {

	file, err := os.Open(filePath)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to open file: %v", err)}
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to read file: %v", err)}
	}

	return ReturnValue{FileContent: string(content)}
}

func (a *App) selectedFileOutputPath() string {
	return filepath.Join(TEMP_DIR_PATH, a.explorerState.SelectedFile.Path)
}

func (a *App) selectedFileReportPath() string {
	return filepath.Join(a.selectedFileOutputPath(), "report.json")

}

func (a *App) selectedFileStorePath() string {
	return filepath.Join(a.selectedFileOutputPath(), "store")
}

func (a *App) ExecuteHurl(filePath string) ReturnValue {

	// Create dir if not exists
	if err := os.MkdirAll(TEMP_DIR_PATH, 0755); err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to create temp dir: %w", err)}
	}

	if _, err := os.Stat(a.explorerState.SelectedFile.Path); err != nil {
		return ReturnValue{Error: fmt.Sprintf("file does not exist: %w", err)}
	}

	outputDir := a.selectedFileOutputPath()
	reportPath := a.selectedFileReportPath()
	outputBodyDir := a.selectedFileStorePath()
	// Delete the content inside the dir
	os.RemoveAll(outputDir)
	os.RemoveAll(outputBodyDir)
	os.Remove(reportPath)
	os.MkdirAll(outputDir, 0755)

	command := []string{"hurl", "--report-json", outputDir, a.explorerState.SelectedFile.Path}
	cmd := exec.Command(command[0], command[1:]...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to execute hurl: %w", string(bytes))}
	}

	// Read and parse JSON report from outputDir
	var report HurlReport
	if reportData, readErr := os.ReadFile(reportPath); readErr == nil {
		if parseErr := json.Unmarshal(reportData, &report); parseErr != nil {
			fmt.Printf("Failed to parse JSON report: %v\n", parseErr)
		}
	} else {
		fmt.Printf("Failed to read JSON report: %v\n", readErr)
	}

	marshalled, err := json.Marshal(report)
	fmt.Println(string(marshalled))

	a.insertResponseData(&report, outputBodyDir)

	return ReturnValue{HurlReport: report}
}

func (a *App) insertResponseData(h *HurlReport, filePath string) error {

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

func (a *App) GetHurlResult(filePath string) ReturnValue {

	reportPath := a.selectedFileReportPath()
	outputPath := a.selectedFileOutputPath()

	// Check if the dir exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return ReturnValue{}
	}

	var report HurlReport
	if reportData, readErr := os.ReadFile(reportPath); readErr == nil {
		if parseErr := json.Unmarshal(reportData, &report); parseErr != nil {
			fmt.Printf("Failed to parse JSON report: %v\n", parseErr)
		} else {
			// Insert response body data from files
			if err := a.insertResponseData(&report, outputPath); err != nil {
				fmt.Printf("Failed to insert response data: %v\n", err)
			}
		}
	} else {
		fmt.Printf("Failed to read JSON report: %v\n", readErr)
	}
	return ReturnValue{HurlReport: report}
}

func (a *App) CreateNewFile(fileName string, fileContent string) ReturnValue {
	// Create a new file in the current directory
	filePath := filepath.Join(a.explorerState.CurrentDir.Path, fileName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return ReturnValue{Error: fmt.Sprintf("file already exists: %s", fileName)}
	}

	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to create new file: %w", err)}
	}

	newFile, err := createFileInfo(filePath)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to create file info: %w", err)}
	}

	fmt.Println("New file created:", newFile.Name)

	a.SetCurrentFile(a.ctx, newFile)
	// a.explorerState.SelectedFile = newFile

	return ReturnValue{}
}

func (a *App) WriteToSelectedFile(content string) ReturnValue {

	filePath := a.explorerState.SelectedFile.Path

	if filePath == "" {
		return ReturnValue{Error: "no file selected"}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ReturnValue{Error: fmt.Sprintf("file does not exist: %s", filePath)}
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to write to file: %w", err)}
	}

	return ReturnValue{}
}

func (a *App) CreateFolder(folderName string) ReturnValue {
	// Create a new folder in the current directory
	folderPath := filepath.Join(a.explorerState.CurrentDir.Path, folderName)

	// Check if folder already exists
	if _, err := os.Stat(folderPath); err == nil {
		return ReturnValue{Error: fmt.Sprintf("folder already exists: %s", folderName)}
	}

	if err := os.Mkdir(folderPath, 0755); err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to create new folder: %w", err)}
	}

	// Do not change current directory; just acknowledge success.
	// Frontend will refresh listing and remain in the same directory.
	return ReturnValue{FileExplorer: a.explorerState}
}

// GetFileContentAndExecuteHurl reads a hurl file content and executes it
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

func (a *App) loadEnvConfig() (*EnvConfig, error) {
	configDir, err := a.getConfigDir()
	if err != nil {
		return nil, err
	}

	envConfigPath := filepath.Join(configDir, "env.json")

	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		return &EnvConfig{
			Global:     make(map[string]string),
			Selectable: make(map[string]map[string]string),
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
	if config.Selectable == nil {
		config.Selectable = make(map[string]map[string]string)
	}

	return &config, nil
}

func (a *App) GetEnvVars(selectedEnvGroup string) ReturnValue {
	config, err := a.loadEnvConfig()
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	envVars := make(map[string]string)

	for key, value := range config.Global {
		envVars[key] = value
	}

	if selectedEnvGroup != "" {
		if selectedGroup, exists := config.Selectable[selectedEnvGroup]; exists {
			for key, value := range selectedGroup {
				envVars[key] = value
			}
		}
	}

	return ReturnValue{EnvVars: envVars}
}

func (a *App) GetAvailableEnvGroups() ReturnValue {
	config, err := a.loadEnvConfig()
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	groups := make([]string, 0, len(config.Selectable))
	for groupName := range config.Selectable {
		groups = append(groups, groupName)
	}

	return ReturnValue{Files: []FileInfo{{Name: "env_groups"}}, EnvVars: map[string]string{"groups": fmt.Sprintf("%v", groups)}}
}

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
