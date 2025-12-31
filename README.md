[![Build Status](https://github.com/axetroy/mcp-server-devtools/workflows/ci/badge.svg)](https://github.com/axetroy/mcp-server-devtools/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/mcp-server-devtools)](https://goreportcard.com/report/github.com/axetroy/mcp-server-devtools)
![Latest Version](https://img.shields.io/github/v/release/axetroy/mcp-server-devtools.svg)
![License](https://img.shields.io/github/license/axetroy/mcp-server-devtools.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/mcp-server-devtools.svg)

# MCP DevTools

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that provides useful developer tools for local development.

MCP DevTools is built with the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) and exposes a collection of tools through a standardized interface, allowing AI assistants and other MCP clients to interact with your local development environment.

## Features

This MCP server provides the following tools:

### 🎨 Color Utilities

- **`color_convert`** - Convert CSS color values between various formats
  - **Input:** CSS color value (e.g., `#ff5733`, `rgb(255, 87, 51)`, `hsl(9, 100%, 60%)`, or named colors like `red`)
  - **Output:** Hex, RGB, HSL, HSV, CMYK, LAB, XYZ, Linear RGB representations
  - **Additional:** Luminance value and light/dark classification

### 🌐 Network Utilities

- **`get_ip_address`** - Get the current computer's IP addresses
  - **Returns:** All active network interface IP addresses
  - **Identifies:** Primary IP address (first non-loopback IPv4)

### 🕐 Time Utilities

- **`get_current_time`** - Get the current server time
  - **Returns:** Current time in RFC1123 format

### 📦 Package Management Utilities

- **`npm_dependencies_analyze`** - Analyze npm package information and build complete dependency tree
  - **Input:** Package name (e.g., `express`, `react`, `@types/node`), optional version, and optional max depth
  - **Output:** Package metadata and complete dependency tree with nested dependencies
  - **Features:** 
    - Builds complete dependency tree recursively (not just direct dependencies)
    - Tracks circular dependencies
    - Configurable depth limit (default: 5, max: 10)
    - Shows total unique dependencies count
    - Fetches data from the official npm registry
  - **Supports:** Scoped packages with proper URL encoding

### 📁 File System Utilities

- **`list_old_downloads`** - Find old files in Downloads folder
  - **Returns:** List of files in the Downloads directory that haven't been modified in the last 3 months
  - **Includes:** File name, last modified time, and size

### 💻 System Utilities

- **`list_installed_apps`** - List installed applications
  - **Platform:** Currently supports macOS only
  - **Returns:** List of installed applications from `/Applications` directory

### 🌍 Browser Utilities

- **`open_in_browser`** - Open URLs in the default browser
  - **Input:** URL to open
  - **Platforms:** Windows, macOS, Linux
  - **Action:** Opens the specified URL in the system's default web browser

## Installation

### Option 1: Shell Script (Mac/Linux)

```bash
curl -fsSL https://github.com/release-lab/install/raw/v1/install.sh | bash -s -- -r=axetroy/mcp-server-devtools -e=mcp-server-devtools
```

### Option 2: PowerShell (Windows)

```powershell
$r="axetroy/mcp-server-devtools";$e="mcp-server-devtools";iwr https://github.com/release-lab/install/raw/v1/install.ps1 -useb | iex
```

### Option 3: Download Binary

Download the pre-built executable from the [GitHub Releases page](https://github.com/axetroy/mcp-server-devtools/releases) and add it to your `$PATH`.

### Option 4: Build from Source

Requires [Go](https://golang.org) 1.21 or later:

```bash
go install github.com/axetroy/mcp-server-devtools/cmd/mcp-server-devtools@latest
```

## Configuration

To use this server with an MCP client, add it to your client's configuration file.

### Claude Desktop

Add to your Claude Desktop configuration file:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "devtools": {
      "command": "mcp-server-devtools"
    }
  }
}
```

### Other MCP Clients

For other MCP clients, configure them to run:

```bash
mcp-server-devtools
```

The server communicates via stdin/stdout following the Model Context Protocol specification.

## Usage

Once configured, your MCP client can use the available tools. The server will:

1. Start when the MCP client initializes
2. Wait for tool call requests on stdin
3. Execute the requested tool
4. Return results via stdout

### Example Tool Calls

<details>
<summary><strong>Color Conversion</strong></summary>

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "color_convert",
    "arguments": {
      "color": "#ff5733"
    }
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "hex": "#ff5733",
    "rgb": "rgb(255, 87, 51)",
    "hsl": "hsl(9.0, 100.0%, 60.0%)",
    "hsv": "hsv(9.0, 80.0%, 100.0%)",
    "cmyk": "cmyk(0.0%, 65.9%, 80.0%, 0.0%)",
    "lab": "lab(61.57, 56.45, 51.48)",
    "xyz": "xyz(0.469, 0.305, 0.074)",
    "linear_rgb": "linear-rgb(1.000, 0.106, 0.030)",
    "luminance": 0.428,
    "is_light": false,
    "is_dark": true,
    "original": "#ff5733"
  }
}
```

</details>

<details>
<summary><strong>Get IP Address</strong></summary>

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "get_ip_address",
    "arguments": {}
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "addresses": ["192.168.1.100", "fe80::1"],
    "primary": "192.168.1.100"
  }
}
```

</details>

<details>
<summary><strong>NPM Package Analysis</strong></summary>

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "npm_dependencies_analyze",
    "arguments": {
      "package_name": "express",
      "max_depth": 3
    }
  }
}
```

Response (simplified for brevity):
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "name": "express",
    "version": "5.2.1",
    "description": "Fast, unopinionated, minimalist web framework",
    "license": "MIT",
    "homepage": "http://expressjs.com/",
    "repository": "https://github.com/expressjs/express",
    "author": "TJ Holowaychuk <tj@vision-media.ca>",
    "keywords": ["express", "framework", "sinatra", "web", "http", "rest", "restful", "router", "app", "api"],
    "latest_version": "5.2.1",
    "publish_time": "2024-12-25T14:49:15.000Z",
    "dependency_tree": {
      "accepts": {
        "name": "accepts",
        "version": "1.3.8",
        "version_range": "~1.3.8",
        "dependencies": {
          "mime-types": {
            "name": "mime-types",
            "version": "2.1.35",
            "version_range": "~2.1.34",
            "dependencies": {
              "mime-db": {
                "name": "mime-db",
                "version": "1.52.0",
                "version_range": "1.52.0"
              }
            }
          },
          "negotiator": {
            "name": "negotiator",
            "version": "0.6.3",
            "version_range": "0.6.3"
          }
        }
      },
      "body-parser": {
        "name": "body-parser",
        "version": "2.1.0",
        "version_range": "2.1.0",
        "dependencies": {
          "...": "..."
        }
      }
    },
    "total_dependencies": 43,
    "tree_depth": 3
  }
}
```

You can also analyze a specific version and control depth:
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "npm_dependencies_analyze",
    "arguments": {
      "package_name": "express",
      "version": "4.18.0",
      "max_depth": 2
    }
  }
}
```

</details>

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Format Code

```bash
make format
```

### Local Development with MCP Client

For local development and testing with an MCP client, you can configure it to run the server from source:

```json
{
  "mcpServers": {
    "devtools": {
      "command": "go",
      "args": ["run", "cmd/mcp-server-devtools/main.go"],
      "cwd": "/path/to/mcp-server-devtools"
    }
  }
}
```

## License

The [MIT License](LICENSE)
