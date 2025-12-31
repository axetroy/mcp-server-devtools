package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// npmPackageInput represents the input for npm dependencies analysis tool
type npmPackageInput struct {
	PackageName string `json:"package_name" jsonschema:"The npm package name to analyze (e.g., 'express', 'react', '@types/node')"`
	Version     string `json:"version,omitempty" jsonschema:"Optional: specific version to analyze (e.g., '4.18.0'). If not provided, analyzes the latest version."`
}

// npmPackageOutput represents the analyzed npm package information
type npmPackageOutput struct {
	Name             string            `json:"name" jsonschema:"Package name"`
	Version          string            `json:"version" jsonschema:"Package version analyzed"`
	Description      string            `json:"description" jsonschema:"Package description"`
	License          string            `json:"license,omitempty" jsonschema:"Package license"`
	Homepage         string            `json:"homepage,omitempty" jsonschema:"Package homepage URL"`
	Repository       string            `json:"repository,omitempty" jsonschema:"Package repository URL"`
	Dependencies     map[string]string `json:"dependencies,omitempty" jsonschema:"Production dependencies with their version ranges"`
	DevDependencies  map[string]string `json:"dev_dependencies,omitempty" jsonschema:"Development dependencies with their version ranges"`
	PeerDependencies map[string]string `json:"peer_dependencies,omitempty" jsonschema:"Peer dependencies with their version ranges"`
	DependencyCount  int               `json:"dependency_count" jsonschema:"Total number of production dependencies"`
	Author           string            `json:"author,omitempty" jsonschema:"Package author"`
	Keywords         []string          `json:"keywords,omitempty" jsonschema:"Package keywords"`
	LatestVersion    string            `json:"latest_version" jsonschema:"Latest available version of the package"`
	PublishTime      string            `json:"publish_time,omitempty" jsonschema:"Time when this version was published"`
}

// npmRegistryResponse represents the npm registry API response structure
type npmRegistryResponse struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	DistTags    map[string]string            `json:"dist-tags"`
	License     interface{}                  `json:"license"`
	Homepage    string                       `json:"homepage"`
	Repository  interface{}                  `json:"repository"`
	Author      interface{}                  `json:"author"`
	Keywords    []string                     `json:"keywords"`
	Versions    map[string]npmVersionDetails `json:"versions"`
	Time        map[string]string            `json:"time"`
}

type npmVersionDetails struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Description      string            `json:"description"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

// NpmDependenciesAnalyze fetches and analyzes npm package information and its dependencies
func NpmDependenciesAnalyze(ctx context.Context, req *mcp.CallToolRequest, input npmPackageInput) (*mcp.CallToolResult, *npmPackageOutput, error) {
	if input.PackageName == "" {
		return nil, nil, fmt.Errorf("package_name is required")
	}

	// Construct the npm registry URL (URL encode the package name for scoped packages)
	registryURL := fmt.Sprintf("https://registry.npmjs.org/%s", url.PathEscape(input.PackageName))

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch package information: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil, fmt.Errorf("package '%s' not found in npm registry", input.PackageName)
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("npm registry returned status code %d", resp.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var registryData npmRegistryResponse
	if err := json.Unmarshal(body, &registryData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse npm registry response: %w", err)
	}

	// Determine which version to analyze
	versionToAnalyze := input.Version
	if versionToAnalyze == "" {
		// Use the latest version
		if registryData.DistTags != nil {
			if latestVersion, ok := registryData.DistTags["latest"]; ok {
				versionToAnalyze = latestVersion
			}
		}
		if versionToAnalyze == "" {
			return nil, nil, fmt.Errorf("no latest version found for package '%s'", input.PackageName)
		}
	}

	// Get the specific version details
	versionDetails, ok := registryData.Versions[versionToAnalyze]
	if !ok {
		return nil, nil, fmt.Errorf("version '%s' not found for package '%s'", versionToAnalyze, input.PackageName)
	}

	// Extract license information
	license := extractLicense(registryData.License)

	// Extract repository URL
	repoURL := extractRepository(registryData.Repository)

	// Extract author information
	author := extractAuthor(registryData.Author)

	// Get publish time
	publishTime := ""
	if t, ok := registryData.Time[versionToAnalyze]; ok {
		publishTime = t
	}

	// Build the output
	latestVersion := ""
	if registryData.DistTags != nil {
		latestVersion = registryData.DistTags["latest"]
	}

	output := &npmPackageOutput{
		Name:             registryData.Name,
		Version:          versionToAnalyze,
		Description:      versionDetails.Description,
		License:          license,
		Homepage:         registryData.Homepage,
		Repository:       repoURL,
		Dependencies:     versionDetails.Dependencies,
		DevDependencies:  versionDetails.DevDependencies,
		PeerDependencies: versionDetails.PeerDependencies,
		DependencyCount:  len(versionDetails.Dependencies),
		Author:           author,
		Keywords:         registryData.Keywords,
		LatestVersion:    latestVersion,
		PublishTime:      publishTime,
	}

	// Initialize empty maps if nil
	if output.Dependencies == nil {
		output.Dependencies = make(map[string]string)
	}
	if output.DevDependencies == nil {
		output.DevDependencies = make(map[string]string)
	}
	if output.PeerDependencies == nil {
		output.PeerDependencies = make(map[string]string)
	}
	if output.Keywords == nil {
		output.Keywords = []string{}
	}

	return nil, output, nil
}

// extractLicense handles different license field formats
func extractLicense(license interface{}) string {
	if license == nil {
		return ""
	}

	switch v := license.(type) {
	case string:
		return v
	case map[string]interface{}:
		if licenseType, ok := v["type"].(string); ok {
			return licenseType
		}
	}
	return ""
}

// extractRepository handles different repository field formats
func extractRepository(repo interface{}) string {
	if repo == nil {
		return ""
	}

	switch v := repo.(type) {
	case string:
		return v
	case map[string]interface{}:
		if repoURL, ok := v["url"].(string); ok {
			// Clean up git+ prefix and .git suffix
			repoURL = strings.TrimPrefix(repoURL, "git+")
			repoURL = strings.TrimSuffix(repoURL, ".git")
			return repoURL
		}
	}
	return ""
}

// extractAuthor handles different author field formats
func extractAuthor(author interface{}) string {
	if author == nil {
		return ""
	}

	switch v := author.(type) {
	case string:
		return v
	case map[string]interface{}:
		name, hasName := v["name"].(string)
		email, hasEmail := v["email"].(string)

		if hasName && hasEmail {
			return fmt.Sprintf("%s <%s>", name, email)
		} else if hasName {
			return name
		} else if hasEmail {
			return email
		}
	}
	return ""
}
