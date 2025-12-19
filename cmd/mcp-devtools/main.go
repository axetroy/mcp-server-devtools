package main

import (
	"fmt"
	"os"

	"github.com/axetroy/mcp-devtools/internal/mcp"
	"github.com/axetroy/mcp-devtools/internal/tools"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Create MCP server
	server := mcp.NewServer("mcp-devtools", version)

	// Register file operation tools
	server.RegisterTool(mcp.Tool{
		Name:        "read_file",
		Description: "Read the contents of a file",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file to read",
				},
			},
			"required": []string{"path"},
		},
	}, tools.ReadFile)

	server.RegisterTool(mcp.Tool{
		Name:        "write_file",
		Description: "Write content to a file",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file to write",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "Content to write to the file",
				},
			},
			"required": []string{"path", "content"},
		},
	}, tools.WriteFile)

	server.RegisterTool(mcp.Tool{
		Name:        "list_files",
		Description: "List files in a directory",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the directory to list (defaults to current directory)",
				},
			},
		},
	}, tools.ListFiles)

	server.RegisterTool(mcp.Tool{
		Name:        "file_exists",
		Description: "Check if a file exists",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file to check",
				},
			},
			"required": []string{"path"},
		},
	}, tools.FileExists)

	server.RegisterTool(mcp.Tool{
		Name:        "delete_file",
		Description: "Delete a file",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file to delete",
				},
			},
			"required": []string{"path"},
		},
	}, tools.DeleteFile)

	// Register command execution tools
	server.RegisterTool(mcp.Tool{
		Name:        "execute_command",
		Description: "Execute a shell command",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{
					"type":        "string",
					"description": "Command to execute",
				},
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for the command (optional)",
				},
			},
			"required": []string{"command"},
		},
	}, tools.ExecuteCommand)

	server.RegisterTool(mcp.Tool{
		Name:        "get_environment",
		Description: "Get environment variables",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetEnvironment)

	server.RegisterTool(mcp.Tool{
		Name:        "get_working_directory",
		Description: "Get the current working directory",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetWorkingDirectory)

	// Register git tools
	server.RegisterTool(mcp.Tool{
		Name:        "git_status",
		Description: "Get git status",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
			},
		},
	}, tools.GitStatus)

	server.RegisterTool(mcp.Tool{
		Name:        "git_diff",
		Description: "Get git diff",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
				"file": map[string]interface{}{
					"type":        "string",
					"description": "Specific file to diff (optional)",
				},
			},
		},
	}, tools.GitDiff)

	server.RegisterTool(mcp.Tool{
		Name:        "git_log",
		Description: "Get git log",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
				"limit": map[string]interface{}{
					"type":        "string",
					"description": "Number of commits to show (default: 10)",
				},
			},
		},
	}, tools.GitLog)

	server.RegisterTool(mcp.Tool{
		Name:        "git_branch",
		Description: "List git branches",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
			},
		},
	}, tools.GitBranch)

	server.RegisterTool(mcp.Tool{
		Name:        "git_add",
		Description: "Stage files for commit",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"files": map[string]interface{}{
					"type":        "string",
					"description": "Files to stage (space-separated, default: .)",
				},
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
			},
		},
	}, tools.GitAdd)

	server.RegisterTool(mcp.Tool{
		Name:        "git_commit",
		Description: "Create a git commit",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Commit message",
				},
				"workdir": map[string]interface{}{
					"type":        "string",
					"description": "Working directory for git command (optional)",
				},
			},
			"required": []string{"message"},
		},
	}, tools.GitCommit)

	// Register system tools
	server.RegisterTool(mcp.Tool{
		Name:        "get_system_info",
		Description: "Get system information",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetSystemInfo)

	server.RegisterTool(mcp.Tool{
		Name:        "get_hostname",
		Description: "Get system hostname",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetHostname)

	server.RegisterTool(mcp.Tool{
		Name:        "get_disk_usage",
		Description: "Get disk usage information",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to check disk usage (default: current directory)",
				},
			},
		},
	}, tools.GetDiskUsage)

	server.RegisterTool(mcp.Tool{
		Name:        "get_process_list",
		Description: "Get list of running processes",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetProcessList)

	server.RegisterTool(mcp.Tool{
		Name:        "get_network_info",
		Description: "Get network information",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, tools.GetNetworkInfo)

	// Run the server
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
