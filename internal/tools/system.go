package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// GetSystemInfo returns system information
func GetSystemInfo(args map[string]interface{}) (interface{}, error) {
	result := fmt.Sprintf("OS: %s\nArchitecture: %s\nGo Version: %s\nCPU Count: %d\n",
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
		runtime.NumCPU(),
	)

	return result, nil
}

// GetHostname returns the system hostname
func GetHostname(args map[string]interface{}) (interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return hostname, nil
}

// GetDiskUsage returns disk usage information
func GetDiskUsage(args map[string]interface{}) (interface{}, error) {
	path, _ := args["path"].(string)
	if path == "" {
		path = "."
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption")
	} else {
		cmd = exec.Command("df", "-h", path)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage: %w", err)
	}

	return string(output), nil
}

// GetProcessList returns a list of running processes
func GetProcessList(args map[string]interface{}) (interface{}, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else {
		cmd = exec.Command("ps", "aux")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %w", err)
	}

	return string(output), nil
}

// GetNetworkInfo returns network information
func GetNetworkInfo(args map[string]interface{}) (interface{}, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ipconfig")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get network info: %w", err)
		}
		return string(output), nil
	}

	cmd = exec.Command("ifconfig")
	output, err := cmd.Output()
	if err != nil {
		// Try ip command if ifconfig is not available
		cmd = exec.Command("ip", "addr")
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get network info: %w", err)
		}
		return string(output), nil
	}
	return string(output), nil
}
