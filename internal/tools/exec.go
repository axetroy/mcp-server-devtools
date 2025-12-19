package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand executes a shell command
func ExecuteCommand(args map[string]interface{}) (interface{}, error) {
	command, ok := args["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}

	workDir, _ := args["workdir"].(string)

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := fmt.Sprintf("STDOUT:\n%s\n\nSTDERR:\n%s\n", stdout.String(), stderr.String())

	if err != nil {
		result += fmt.Sprintf("\nError: %v", err)
	}

	return result, nil
}

// GetEnvironment gets environment variables
func GetEnvironment(args map[string]interface{}) (interface{}, error) {
	cmd := exec.Command("env")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	return string(output), nil
}

// GetWorkingDirectory gets the current working directory
func GetWorkingDirectory(args map[string]interface{}) (interface{}, error) {
	cmd := exec.Command("pwd")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
