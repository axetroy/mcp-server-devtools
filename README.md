[![Build Status](https://github.com/axetroy/mcp-devtools/workflows/ci/badge.svg)](https://github.com/axetroy/mcp-devtools/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/mcp-devtools)](https://goreportcard.com/report/github.com/axetroy/mcp-devtools)
![Latest Version](https://img.shields.io/github/v/release/axetroy/mcp-devtools.svg)
![License](https://img.shields.io/github/license/axetroy/mcp-devtools.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/mcp-devtools.svg)

## MCP DevTools

> A Model Context Protocol (MCP) server that provides a comprehensive set of developer tools for local development.

MCP DevTools is an MCP server implementation that exposes useful development tools through the Model Context Protocol. It allows AI assistants and other MCP clients to interact with your local development environment through a standardized interface.

### Features

This MCP server provides the following categories of tools:

#### File Operations
- **read_file** - Read the contents of a file
- **write_file** - Write content to a file
- **list_files** - List files in a directory
- **file_exists** - Check if a file exists
- **delete_file** - Delete a file

#### Command Execution
- **execute_command** - Execute shell commands
- **get_environment** - Get environment variables
- **get_working_directory** - Get the current working directory

#### Git Operations
- **git_status** - Get repository status
- **git_diff** - Get diff of changes
- **git_log** - View commit history
- **git_branch** - List branches
- **git_add** - Stage files for commit
- **git_commit** - Create commits

#### System Information
- **get_system_info** - Get OS, architecture, and system details
- **get_hostname** - Get system hostname
- **get_disk_usage** - Get disk usage information
- **get_process_list** - List running processes
- **get_network_info** - Get network configuration

### Usage

The MCP DevTools server communicates via JSON-RPC over stdin/stdout. It follows the Model Context Protocol specification.

To use with an MCP client:

```bash
mcp-devtools
```

The server will start and wait for JSON-RPC requests on stdin, sending responses to stdout.

### Install

1. Shell (Mac/Linux)

   ```bash
   curl -fsSL https://github.com/release-lab/install/raw/v1/install.sh | bash -s -- -r=axetroy/mcp-devtools -e=mcp-devtools
   ```

2. PowerShell (Windows):

   ```powershell
   $r="axetroy/mcp-devtools";$e="mcp-devtools";iwr https://github.com/release-lab/install/raw/v1/install.ps1 -useb | iex
   ```

3. [Github release page](https://github.com/axetroy/mcp-devtools/releases) (All platforms)

   download the executable file and put the executable file to `$PATH`

4. Build and install from source using [Golang](https://golang.org) (All platforms)

   ```bash
   go install github.com/axetroy/mcp-devtools/cmd/mcp-devtools@latest
   ```

### Development

Build the project:

```bash
make build
```

Run tests:

```bash
make test
```

Format code:

```bash
make format
```

### License

The [MIT License](LICENSE)
