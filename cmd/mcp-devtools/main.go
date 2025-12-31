package main

import (
	"context"
	"log"

	"github.com/axetroy/mcp-server-devtools/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Create MCP server using the official SDK
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-server-devtools",
			Version: version,
		},
		&mcp.ServerOptions{
			Instructions: "A collection of useful developer tools including color conversion and network information.",
		},
	)

	// Register color conversion tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "color_convert",
		Description: "Convert CSS color values to various color formats (Hex, RGB, HSL, HSV, CMYK, LAB, XYZ, Linear RGB). Supports hex (#ff5733), rgb(255, 87, 51), hsl(9, 100%, 60%), and named colors (red, blue, etc.)",
	}, tools.ColorConversion)

	// Register IP address tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_ip_address",
		Description: "Get the current computer's IP addresses, including all network interfaces and the primary IP address",
	}, tools.GetIPAddress)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_time",
		Description: "Get the current server time in RFC1123 format",
	}, tools.GetCurrentTime)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_old_downloads",
		Description: "List files in the Download directory that haven't been modified in a long time.",
	}, tools.ListOldDownloads)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_installed_apps",
		Description: "List installed applications on the system (currently supports macOS only).",
	}, tools.ListInstalledApps)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "open_in_browser",
		Description: "Open a specified URL in the default web browser of the system.",
	}, tools.OpenInBrowser)

	log.Println("MCP server started (version:", version, "commit:", commit, "date:", date, "builtBy:", builtBy+")")

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
