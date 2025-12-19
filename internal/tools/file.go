package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads a file and returns its content
func ReadFile(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// WriteFile writes content to a file
func WriteFile(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}

	content, ok := args["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter is required")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote to %s", path), nil
}

// ListFiles lists files in a directory
func ListFiles(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		path = "."
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	result := ""
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		prefix := "f"
		if entry.IsDir() {
			prefix = "d"
		}

		result += fmt.Sprintf("[%s] %s (%d bytes)\n", prefix, entry.Name(), info.Size())
	}

	return result, nil
}

// FileExists checks if a file exists
func FileExists(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}

	_, err := os.Stat(path)
	if err == nil {
		return "File exists: true", nil
	}
	if os.IsNotExist(err) {
		return "File exists: false", nil
	}

	return nil, fmt.Errorf("error checking file: %w", err)
}

// DeleteFile deletes a file
func DeleteFile(args map[string]interface{}) (interface{}, error) {
	path, ok := args["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}

	if err := os.Remove(path); err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return fmt.Sprintf("Successfully deleted %s", path), nil
}
