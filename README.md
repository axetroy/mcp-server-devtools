[![Build Status](https://github.com/axetroy/mcp-server-devtools/workflows/ci/badge.svg)](https://github.com/axetroy/mcp-server-devtools/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/mcp-server-devtools)](https://goreportcard.com/report/github.com/axetroy/mcp-server-devtools)
![Latest Version](https://img.shields.io/github/v/release/axetroy/mcp-server-devtools.svg)
![License](https://img.shields.io/github/license/axetroy/mcp-server-devtools.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/mcp-server-devtools.svg)

## MCP DevTools

> A Model Context Protocol (MCP) server that provides useful developer tools for local development.

MCP DevTools is an MCP server implementation built with the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) that exposes useful development tools through the Model Context Protocol. It allows AI assistants and other MCP clients to interact with your local development environment through a standardized interface.

### Features

This MCP server provides the following tools:

#### Color Conversion
- **color_convert** - Convert CSS color values to various color formats
  - Input: CSS color value (e.g., `#ff5733`, `rgb(255, 87, 51)`, `hsl(9, 100%, 60%)`, or named colors like `red`)
  - Output: Hex, RGB, HSL, HSV, CMYK, LAB, XYZ, Linear RGB representations
  - Additional info: Luminance, whether the color is light or dark

#### Network Information
- **get_ip_address** - Get the current computer's IP addresses
  - Returns all active network interface IP addresses
  - Identifies the primary IP address (first non-loopback IPv4)

### Usage

The MCP DevTools server communicates via the Model Context Protocol over stdin/stdout. It follows the MCP specification and is built with the official Go SDK.

To use with an MCP client:

```bash
mcp-server-devtools
```

The server will start and wait for MCP requests on stdin, sending responses to stdout.

#### Example: Color Conversion

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

#### Example: Get IP Address

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

### Install

1. Shell (Mac/Linux)

   ```bash
   curl -fsSL https://github.com/release-lab/install/raw/v1/install.sh | bash -s -- -r=axetroy/mcp-server-devtools -e=mcp-server-devtools
   ```

2. PowerShell (Windows):

   ```powershell
   $r="axetroy/mcp-server-devtools";$e="mcp-server-devtools";iwr https://github.com/release-lab/install/raw/v1/install.ps1 -useb | iex
   ```

3. [Github release page](https://github.com/axetroy/mcp-server-devtools/releases) (All platforms)

   download the executable file and put the executable file to `$PATH`

4. Build and install from source using [Golang](https://golang.org) (All platforms)

   ```bash
   go install github.com/axetroy/mcp-server-devtools/cmd/mcp-server-devtools@latest
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
