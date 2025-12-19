package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitStatus gets the git status
func GitStatus(args map[string]interface{}) (interface{}, error) {
	workDir, _ := args["workdir"].(string)

	cmd := exec.Command("git", "status")
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git status failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// GitDiff gets the git diff
func GitDiff(args map[string]interface{}) (interface{}, error) {
	workDir, _ := args["workdir"].(string)
	file, _ := args["file"].(string)

	cmdArgs := []string{"diff"}
	if file != "" {
		cmdArgs = append(cmdArgs, file)
	}

	cmd := exec.Command("git", cmdArgs...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// GitLog gets the git log
func GitLog(args map[string]interface{}) (interface{}, error) {
	workDir, _ := args["workdir"].(string)
	limitStr, _ := args["limit"].(string)

	limit := "10"
	if limitStr != "" {
		limit = limitStr
	}

	cmd := exec.Command("git", "log", "-n", limit, "--oneline")
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// GitBranch lists git branches
func GitBranch(args map[string]interface{}) (interface{}, error) {
	workDir, _ := args["workdir"].(string)

	cmd := exec.Command("git", "branch", "-a")
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git branch failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// GitCommit creates a git commit
func GitCommit(args map[string]interface{}) (interface{}, error) {
	message, ok := args["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required")
	}

	workDir, _ := args["workdir"].(string)

	cmd := exec.Command("git", "commit", "-m", message)
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git commit failed: %w\n%s", err, string(output))
	}

	return string(output), nil
}

// GitAdd stages files for commit
func GitAdd(args map[string]interface{}) (interface{}, error) {
	files, ok := args["files"].(string)
	if !ok {
		files = "."
	}

	workDir, _ := args["workdir"].(string)

	fileList := strings.Fields(files)
	cmdArgs := append([]string{"add"}, fileList...)

	cmd := exec.Command("git", cmdArgs...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git add failed: %w\n%s", err, string(output))
	}

	if len(output) == 0 {
		return "Files staged successfully", nil
	}

	return string(output), nil
}
