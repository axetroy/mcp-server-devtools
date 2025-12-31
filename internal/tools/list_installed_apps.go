package tools

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type app struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type listInstalledAppsOutput struct {
	Apps []app `json:"apps" jsonschema:"List of installed applications"`
}

func getInstalledAppsOnMacOS() ([]app, error) {
	applicationsPath := "/Applications"
	entries, err := os.ReadDir(applicationsPath)
	if err != nil {
		return nil, err
	}

	var apps []app
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".app") {
			apps = append(apps, app{Name: strings.TrimSuffix(entry.Name(), ".app")})
		}
	}

	return apps, nil
}

func ListInstalledApps(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *listInstalledAppsOutput, error) {
	var apps []app
	var err error

	if runtime.GOOS == "darwin" {
		apps, err = getInstalledAppsOnMacOS()
		if err != nil {
			return nil, nil, err
		}
	}

	return nil, &listInstalledAppsOutput{
		Apps: apps,
	}, nil
}
