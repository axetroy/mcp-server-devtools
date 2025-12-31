package tools

import (
	"context"
	"os/exec"
	"runtime"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type openInBrowserInput struct {
	Url string `json:"url" jsonschema:"URL to open in the browser"`
}

type openInBrowserOutput struct {
}

func OpenInBrowser(ctx context.Context, req *mcp.CallToolRequest, input openInBrowserInput) (*mcp.CallToolResult, *openInBrowserOutput, error) {
	// Open the URL in the default browser
	switch runtime.GOOS {
	case "windows":
		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", input.Url).Start()
		if err != nil {
			return nil, nil, err
		}
	case "darwin":
		err := exec.Command("open", input.Url).Start()
		if err != nil {
			return nil, nil, err
		}
	case "linux":
		err := exec.Command("xdg-open", input.Url).Start()
		if err != nil {
			return nil, nil, err
		}
	default:
		// Unsupported OS
	}

	return nil, nil, nil
}
