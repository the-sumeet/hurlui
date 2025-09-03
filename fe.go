package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func (a *App) ExecuteHurl(filePath string, envName string) ReturnValue {

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

	// Build hurl command with env variables
	// Load env config
	config, err := a.loadEnvConfig()
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to load env config: %v", err)}
	}

	// Merge globals with selected environment (env overrides globals)
	vars := map[string]string{}
	for k, v := range config.Global {
		vars[k] = v
	}
	if envName != "" {
		if envMap, ok := config.Environments[envName]; ok {
			for k, v := range envMap {
				vars[k] = v
			}
		}
	}

	command := []string{"hurl", "--report-json", outputDir}
	// Append variables as --variable key=value
	for k, v := range vars {
		command = append(command, "--variable", fmt.Sprintf("%s=%s", k, v))
	}
	// Finally the file path
	command = append(command, a.explorerState.SelectedFile.Path)

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

func (a *App) GetEnvVars() ReturnValue {
	config, err := a.loadEnvConfig()
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	envSelectables := []string{}
	for groupName := range config.Environments {
		envSelectables = append(envSelectables, groupName)
	}

	return ReturnValue{Envs: envSelectables}
}

func (a *App) GetEnvFilePath() ReturnValue {

	envConfigPath, err := a.getEnvFilePath()
	if err != nil {
		return ReturnValue{Error: err.Error()}
	}

	return ReturnValue{EnvFilePath: envConfigPath}
}
