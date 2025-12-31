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

const (
	// npmRegistryTimeout is the timeout for HTTP requests to npm registry
	npmRegistryTimeout = 30 * time.Second
)

// npmPackageInput represents the input for npm dependencies analysis tool
type npmPackageInput struct {
	PackageName string `json:"package_name" jsonschema:"The npm package name to analyze (e.g., 'express', 'react', '@types/node')"`
	Version     string `json:"version,omitempty" jsonschema:"Optional: specific version to analyze (e.g., '4.18.0'). If not provided, analyzes the latest version."`
	MaxDepth    int    `json:"max_depth,omitempty" jsonschema:"Optional: maximum depth to traverse the dependency tree (default: 5, max: 10)."`
}

// DependencyNode represents a node in the dependency tree
type DependencyNode struct {
	Name         string                     `json:"name" jsonschema:"Package name"`
	Version      string                     `json:"version" jsonschema:"Resolved version"`
	VersionRange string                     `json:"version_range,omitempty" jsonschema:"Version range specified by parent"`
	Dependencies map[string]*DependencyNode `json:"dependencies,omitempty" jsonschema:"Nested dependencies"`
	Circular     bool                       `json:"circular,omitempty" jsonschema:"Whether this is a circular dependency reference"`
	DepthLimited bool                       `json:"depth_limited,omitempty" jsonschema:"Whether traversal stopped due to depth limit"`
	Error        string                     `json:"error,omitempty" jsonschema:"Error message if package info couldn't be fetched"`
}

// npmPackageOutput represents the analyzed npm package information
type npmPackageOutput struct {
	Name              string                     `json:"name" jsonschema:"Package name"`
	Version           string                     `json:"version" jsonschema:"Package version analyzed"`
	Description       string                     `json:"description" jsonschema:"Package description"`
	License           string                     `json:"license,omitempty" jsonschema:"Package license"`
	Homepage          string                     `json:"homepage,omitempty" jsonschema:"Package homepage URL"`
	Repository        string                     `json:"repository,omitempty" jsonschema:"Package repository URL"`
	Author            string                     `json:"author,omitempty" jsonschema:"Package author"`
	Keywords          []string                   `json:"keywords,omitempty" jsonschema:"Package keywords"`
	LatestVersion     string                     `json:"latest_version" jsonschema:"Latest available version of the package"`
	PublishTime       string                     `json:"publish_time,omitempty" jsonschema:"Time when this version was published"`
	DependencyTree    map[string]*DependencyNode `json:"dependency_tree" jsonschema:"Complete dependency tree with nested dependencies"`
	TotalDependencies int                        `json:"total_dependencies" jsonschema:"Total number of unique dependencies (including transitive)"`
	TreeDepth         int                        `json:"tree_depth" jsonschema:"Maximum depth of the dependency tree"`
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

// NpmDependenciesAnalyze fetches and analyzes npm package information and builds a complete dependency tree
func NpmDependenciesAnalyze(ctx context.Context, req *mcp.CallToolRequest, input npmPackageInput) (*mcp.CallToolResult, *npmPackageOutput, error) {
	if input.PackageName == "" {
		return nil, nil, fmt.Errorf("package_name is required")
	}

	// Set default max depth if not specified or cap at maximum
	maxDepth := input.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 5 // Default depth
	} else if maxDepth > 10 {
		maxDepth = 10 // Cap at 10 to prevent excessive API calls
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: npmRegistryTimeout,
	}

	// Fetch package metadata
	registryData, versionToAnalyze, err := fetchPackageInfo(client, input.PackageName, input.Version)
	if err != nil {
		return nil, nil, err
	}

	// Get the specific version details
	versionDetails, ok := registryData.Versions[versionToAnalyze]
	if !ok {
		return nil, nil, fmt.Errorf("version '%s' not found for package '%s'", versionToAnalyze, input.PackageName)
	}

	// Extract metadata
	license := extractLicense(registryData.License)
	repoURL := extractRepository(registryData.Repository)
	author := extractAuthor(registryData.Author)

	publishTime := ""
	if t, ok := registryData.Time[versionToAnalyze]; ok {
		publishTime = t
	}

	latestVersion := ""
	if registryData.DistTags != nil {
		latestVersion = registryData.DistTags["latest"]
	}

	description := versionDetails.Description
	if description == "" {
		description = registryData.Description
	}

	// Build dependency tree
	visited := make(map[string]bool)
	dependencyTree := make(map[string]*DependencyNode)
	var maxTreeDepth int

	for depName, depVersion := range versionDetails.Dependencies {
		node, depth := buildDependencyTree(client, depName, depVersion, visited, 1, maxDepth)
		if node != nil {
			dependencyTree[depName] = node
			if depth > maxTreeDepth {
				maxTreeDepth = depth
			}
		}
	}

	// Count total unique dependencies
	totalDeps := countUniqueDependencies(dependencyTree, make(map[string]bool))

	output := &npmPackageOutput{
		Name:              registryData.Name,
		Version:           versionToAnalyze,
		Description:       description,
		License:           license,
		Homepage:          registryData.Homepage,
		Repository:        repoURL,
		Author:            author,
		Keywords:          registryData.Keywords,
		LatestVersion:     latestVersion,
		PublishTime:       publishTime,
		DependencyTree:    dependencyTree,
		TotalDependencies: totalDeps,
		TreeDepth:         maxTreeDepth,
	}

	if output.Keywords == nil {
		output.Keywords = []string{}
	}
	if output.DependencyTree == nil {
		output.DependencyTree = make(map[string]*DependencyNode)
	}

	return nil, output, nil
}

// fetchPackageInfo fetches package information from npm registry
func fetchPackageInfo(client *http.Client, packageName, version string) (*npmRegistryResponse, string, error) {
	registryURL := fmt.Sprintf("https://registry.npmjs.org/%s", url.PathEscape(packageName))

	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch package information: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf("package '%s' not found in npm registry", packageName)
	}

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("npm registry returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	var registryData npmRegistryResponse
	if err := json.Unmarshal(body, &registryData); err != nil {
		return nil, "", fmt.Errorf("failed to parse npm registry response: %w", err)
	}

	// Determine which version to use
	versionToUse := version
	if versionToUse == "" {
		if registryData.DistTags != nil {
			if latestVersion, ok := registryData.DistTags["latest"]; ok {
				versionToUse = latestVersion
			}
		}
		if versionToUse == "" {
			return nil, "", fmt.Errorf("no latest version found for package '%s'", packageName)
		}
	}

	return &registryData, versionToUse, nil
}

// buildDependencyTree recursively builds the dependency tree for a package
func buildDependencyTree(client *http.Client, packageName, versionRange string, visited map[string]bool, currentDepth, maxDepth int) (*DependencyNode, int) {
	// Create a unique key for this package
	packageKey := packageName

	// Check if we've already visited this package (circular dependency)
	if visited[packageKey] {
		return &DependencyNode{
			Name:         packageName,
			VersionRange: versionRange,
			Circular:     true,
		}, currentDepth
	}

	// Check if we've reached max depth
	if currentDepth >= maxDepth {
		return &DependencyNode{
			Name:         packageName,
			VersionRange: versionRange,
			DepthLimited: true,
		}, currentDepth
	}

	// Mark as visited
	visited[packageKey] = true
	defer func(key string) { delete(visited, key) }(packageKey)

	// Fetch package info
	registryData, resolvedVersion, err := fetchPackageInfo(client, packageName, "")
	if err != nil {
		// If we can't fetch the package, return a node with error info
		return &DependencyNode{
			Name:         packageName,
			VersionRange: versionRange,
			Error:        err.Error(),
		}, currentDepth
	}

	versionDetails, ok := registryData.Versions[resolvedVersion]
	if !ok {
		return &DependencyNode{
			Name:         packageName,
			VersionRange: versionRange,
			Version:      resolvedVersion,
			Error:        "version not found in registry",
		}, currentDepth
	}

	// Create node for this dependency
	node := &DependencyNode{
		Name:         packageName,
		Version:      resolvedVersion,
		VersionRange: versionRange,
		Dependencies: make(map[string]*DependencyNode),
	}

	maxChildDepth := currentDepth

	// Recursively build tree for each dependency
	for depName, depVersion := range versionDetails.Dependencies {
		childNode, childDepth := buildDependencyTree(client, depName, depVersion, visited, currentDepth+1, maxDepth)
		if childNode != nil {
			node.Dependencies[depName] = childNode
			if childDepth > maxChildDepth {
				maxChildDepth = childDepth
			}
		}
	}

	return node, maxChildDepth
}

// countUniqueDependencies counts the total number of unique dependencies in the tree
func countUniqueDependencies(tree map[string]*DependencyNode, counted map[string]bool) int {
	total := 0
	for name, node := range tree {
		if !counted[name] && !node.Circular {
			counted[name] = true
			total++
			if node.Dependencies != nil {
				total += countUniqueDependencies(node.Dependencies, counted)
			}
		}
	}
	return total
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
