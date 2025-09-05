package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

    // Persist last opened file if valid
    if file.Path != "" {
        a.preferences.LastOpenedFile = file.Path
        if err := a.savePreferences(); err != nil {
            // Non-fatal: log to console but don't interrupt flow
            fmt.Printf("failed to save preferences: %v\n", err)
        }
    }
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
    a.preferences.LastOpenedDir = fileInfo.Path
    if err := a.savePreferences(); err != nil {
        fmt.Printf("failed to save preferences after change dir: %v\n", err)
    }
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
		return ReturnValue{Error: fmt.Sprintf("failed to execute hurl: %s\n%s", err.Error(), string(bytes))}
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

// RenamePath renames a file or folder within its current directory.
// newName should be a base name (no path separators).
func (a *App) RenamePath(oldPath string, newName string) ReturnValue {
	if oldPath == "" {
		return ReturnValue{Error: "old path is empty"}
	}
	if newName == "" {
		return ReturnValue{Error: "new name is empty"}
	}
	if filepath.Base(newName) != newName {
		return ReturnValue{Error: "new name must not contain path separators"}
	}

	// Ensure the old path exists
	oldInfo, err := os.Stat(oldPath)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("path does not exist: %v", err)}
	}

	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	// Prevent overwriting existing paths
	if _, err := os.Stat(newPath); err == nil {
		return ReturnValue{Error: fmt.Sprintf("target already exists: %s", newPath)}
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return ReturnValue{Error: fmt.Sprintf("failed to rename: %v", err)}
	}

	// Update explorer state if needed
	if a.explorerState.SelectedFile.Path == oldPath {
		a.explorerState.SelectedFile.Path = newPath
		a.explorerState.SelectedFile.Name = filepath.Base(newPath)
		a.explorerState.SelectedFile.IsDir = oldInfo.IsDir()
		// Keep preferences in sync if the selected file was renamed
		if !oldInfo.IsDir() {
			a.preferences.LastOpenedFile = newPath
			if err := a.savePreferences(); err != nil {
				fmt.Printf("failed to save preferences after rename: %v\n", err)
			}
		}
	}
    if a.explorerState.CurrentDir.Path == oldPath && oldInfo.IsDir() {
        a.explorerState.CurrentDir.Path = newPath
        a.explorerState.CurrentDir.Name = filepath.Base(newPath)
        a.explorerState.CurrentDir.IsDir = true
        a.preferences.LastOpenedDir = newPath
        if err := a.savePreferences(); err != nil {
            fmt.Printf("failed to save preferences after dir rename: %v\n", err)
        }
    }

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

// DeletePath deletes a file or directory. Directories are removed recursively.
func (a *App) DeletePath(targetPath string) ReturnValue {
	if targetPath == "" {
		return ReturnValue{Error: "path is empty"}
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return ReturnValue{Error: fmt.Sprintf("path does not exist: %v", err)}
	}

	// Perform deletion
	if info.IsDir() {
		if err := os.RemoveAll(targetPath); err != nil {
			return ReturnValue{Error: fmt.Sprintf("failed to delete folder: %v", err)}
		}
	} else {
		if err := os.Remove(targetPath); err != nil {
			return ReturnValue{Error: fmt.Sprintf("failed to delete file: %v", err)}
		}
	}

	// Also delete any cached hurl results corresponding to the path.
	// If a file is deleted, remove its specific cache dir; if a folder is deleted,
	// remove the mirrored subtree under TEMP_DIR_PATH.
	cacheRoot := tempOutputPathFor(targetPath)
	if info.IsDir() {
		_ = os.RemoveAll(cacheRoot)
	} else {
		// Only .hurl files will have cached results, but removing a non-existent dir is safe
		_ = os.RemoveAll(cacheRoot)
	}

	// Update explorer state when selection or current dir are impacted
	// Clear selection if it was the deleted item or under it
	if a.explorerState.SelectedFile.Path == targetPath ||
		strings.HasPrefix(a.explorerState.SelectedFile.Path, targetPath+string(os.PathSeparator)) {
		a.ClearSelection()
	}

	// If current directory is deleted or was inside the deleted folder, move to parent
    if a.explorerState.CurrentDir.Path == targetPath ||
        strings.HasPrefix(a.explorerState.CurrentDir.Path, targetPath+string(os.PathSeparator)) {
        parent := filepath.Dir(targetPath)
        if stat, err := os.Stat(parent); err == nil && stat.IsDir() {
            if fi, err := createFileInfo(parent); err == nil {
                a.explorerState.CurrentDir = fi
                a.preferences.LastOpenedDir = fi.Path
                if err := a.savePreferences(); err != nil {
                    fmt.Printf("failed to save preferences after delete: %v\n", err)
                }
            }
        }
    }

	return ReturnValue{FileExplorer: a.explorerState}
}
