package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Luminance coefficients for relative luminance calculation (ITU-R BT.709)
const (
	RedLuminance   = 0.2126
	GreenLuminance = 0.7152
	BlueLuminance  = 0.0722
)

// max returns the maximum of three float64 values
func max(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

// rgbToCMYK converts RGB values (0-1 range) to CMYK values (0-1 range)
func rgbToCMYK(r, g, b float64) (c, m, y, k float64) {
	k = 1.0 - max(r, g, b)
	if k < 1.0 {
		c = (1.0 - r - k) / (1.0 - k)
		m = (1.0 - g - k) / (1.0 - k)
		y = (1.0 - b - k) / (1.0 - k)
	}
	return c, m, y, k
}

// ColorInput represents the input for color conversion tool
type ColorInput struct {
	Color string `json:"color" jsonschema:"CSS color value (e.g., '#ff5733', 'rgb(255, 87, 51)', 'hsl(9, 100%, 60%)', 'red')"`
}

// ColorOutput represents the output of color conversion
type ColorOutput struct {
	Hex       string  `json:"hex" jsonschema:"Hexadecimal color representation"`
	RGB       string  `json:"rgb" jsonschema:"RGB color representation"`
	HSL       string  `json:"hsl" jsonschema:"HSL color representation"`
	HSV       string  `json:"hsv" jsonschema:"HSV color representation"`
	CMYK      string  `json:"cmyk" jsonschema:"CMYK color representation"`
	LAB       string  `json:"lab" jsonschema:"LAB color representation"`
	XYZ       string  `json:"xyz" jsonschema:"XYZ color representation"`
	LinearRGB string  `json:"linear_rgb" jsonschema:"Linear RGB color representation"`
	Luminance float64 `json:"luminance" jsonschema:"Relative luminance (0-1)"`
	IsLight   bool    `json:"is_light" jsonschema:"Whether the color is light (luminance > 0.5)"`
	IsDark    bool    `json:"is_dark" jsonschema:"Whether the color is dark (luminance <= 0.5)"`
	Original  string  `json:"original" jsonschema:"Original input color value"`
}

// IPAddressOutput represents the output of IP address tool
type IPAddressOutput struct {
	Addresses []string `json:"addresses" jsonschema:"List of IP addresses"`
	Primary   string   `json:"primary" jsonschema:"Primary IP address (first non-loopback IPv4)"`
}

type CurrentTimeOutput struct {
	Time string `json:"time" jsonschema:"Current server time in RFC1123 format"`
}

type ListOldDownloadsOutput struct {
	System string                `json:"system" jsonschema:"Operating system of the server"`
	Files  []ListOldDownloadFile `json:"files" jsonschema:"List of file paths to check for old downloads"`
}

// ColorConversionTool converts CSS color values to various color formats
func ColorConversionTool(ctx context.Context, req *mcp.CallToolRequest, input ColorInput) (*mcp.CallToolResult, *ColorOutput, error) {
	// Parse the color
	color, err := colorful.Hex(input.Color)
	if err != nil {
		// Try parsing as named color or other CSS format
		color, err = parseColor(input.Color)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse color '%s': %w", input.Color, err)
		}
	}

	// Get various color representations
	r, g, b := color.RGB255()
	h, s, l := color.Hsl()
	hv, sv, v := color.Hsv()

	// Calculate CMYK
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	c, m, y, k := rgbToCMYK(rf, gf, bf)

	lab_l, lab_a, lab_b := color.Lab()
	x, yv, z := color.Xyz()
	lr, lg, lb := color.LinearRgb()
	luminance := (RedLuminance*float64(r) + GreenLuminance*float64(g) + BlueLuminance*float64(b)) / 255.0

	output := &ColorOutput{
		Hex:       color.Hex(),
		RGB:       fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL:       fmt.Sprintf("hsl(%.1f, %.1f%%, %.1f%%)", h, s*100, l*100),
		HSV:       fmt.Sprintf("hsv(%.1f, %.1f%%, %.1f%%)", hv, sv*100, v*100),
		CMYK:      fmt.Sprintf("cmyk(%.1f%%, %.1f%%, %.1f%%, %.1f%%)", c*100, m*100, y*100, k*100),
		LAB:       fmt.Sprintf("lab(%.2f, %.2f, %.2f)", lab_l, lab_a, lab_b),
		XYZ:       fmt.Sprintf("xyz(%.3f, %.3f, %.3f)", x, yv, z),
		LinearRGB: fmt.Sprintf("linear-rgb(%.3f, %.3f, %.3f)", lr, lg, lb),
		Luminance: luminance,
		IsLight:   luminance > 0.5,
		IsDark:    luminance <= 0.5,
		Original:  input.Color,
	}

	return nil, output, nil
}

// parseColor attempts to parse various CSS color formats
func parseColor(colorStr string) (colorful.Color, error) {
	colorStr = strings.TrimSpace(colorStr)

	// Try hex format
	if strings.HasPrefix(colorStr, "#") {
		return colorful.Hex(colorStr)
	}

	// Try RGB format
	if strings.HasPrefix(colorStr, "rgb") {
		var r, g, b uint8
		_, err := fmt.Sscanf(colorStr, "rgb(%d,%d,%d)", &r, &g, &b)
		if err != nil {
			_, err = fmt.Sscanf(colorStr, "rgb(%d, %d, %d)", &r, &g, &b)
		}
		if err == nil {
			return colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}, nil
		}
	}

	// Try HSL format
	if strings.HasPrefix(colorStr, "hsl") {
		var h, s, l float64
		_, err := fmt.Sscanf(colorStr, "hsl(%f,%f%%,%f%%)", &h, &s, &l)
		if err != nil {
			_, err = fmt.Sscanf(colorStr, "hsl(%f, %f%%, %f%%)", &h, &s, &l)
		}
		if err == nil {
			return colorful.Hsl(h, s/100.0, l/100.0), nil
		}
	}

	// Try named colors
	namedColors := map[string]string{
		"red": "#ff0000", "green": "#008000", "blue": "#0000ff",
		"white": "#ffffff", "black": "#000000", "yellow": "#ffff00",
		"cyan": "#00ffff", "magenta": "#ff00ff", "gray": "#808080",
		"orange": "#ffa500", "purple": "#800080", "pink": "#ffc0cb",
		"brown": "#a52a2a", "lime": "#00ff00", "navy": "#000080",
	}

	if hex, ok := namedColors[strings.ToLower(colorStr)]; ok {
		return colorful.Hex(hex)
	}

	return colorful.Color{}, fmt.Errorf("unable to parse color")
}

// GetIPAddressTool returns the current computer's IP addresses
func GetIPAddressTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *IPAddressOutput, error) {
	addresses := []string{}
	primary := ""

	// Get all network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ipStr := ip.String()
			addresses = append(addresses, ipStr)

			// Set primary as the first non-loopback IPv4 address
			if primary == "" && ip.To4() != nil {
				primary = ipStr
			}
		}
	}

	if len(addresses) == 0 {
		return nil, nil, fmt.Errorf("no IP addresses found")
	}

	if primary == "" && len(addresses) > 0 {
		primary = addresses[0]
	}

	return nil, &IPAddressOutput{
		Addresses: addresses,
		Primary:   primary,
	}, nil
}

func GetCurrentTimeTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *CurrentTimeOutput, error) {
	currentTime := fmt.Sprintf("Current server time is: %s", time.Now().Format(time.RFC1123))
	return nil, &CurrentTimeOutput{Time: currentTime}, nil
}

type ListOldDownloadFile struct {
	Name           string    `json:"name" jsonschema:"Name of the old file"`
	LastModifyTime time.Time `json:"last_modify" jsonschema:"Last modify time of the file"`
	Size           int64     `json:"size" jsonschema:"Size of the file in bytes"`
}

// ListOldDownloadsTool lists files in the Download directory that haven't been accessed in a long time.
func ListOldDownloadsTool(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, *ListOldDownloadsOutput, error) {
	downloadDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	downloadDir = downloadDir + string(os.PathSeparator) + "Downloads"

	files, err := os.ReadDir(downloadDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read Downloads directory: %w", err)
	}

	var oldFiles []ListOldDownloadFile

	cutoff := time.Now().AddDate(0, -3, 0) // 3 months ago

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			oldFiles = append(oldFiles, ListOldDownloadFile{
				Name:           info.Name(),
				LastModifyTime: info.ModTime(),
				Size:           info.Size(),
			})
		}
	}

	return nil, &ListOldDownloadsOutput{
		System: runtime.GOOS,
		Files:  oldFiles,
	}, nil
}

func main() {
	// Create MCP server using the official SDK
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-devtools",
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
	}, ColorConversionTool)

	// Register IP address tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_ip_address",
		Description: "Get the current computer's IP addresses, including all network interfaces and the primary IP address",
	}, GetIPAddressTool)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_time",
		Description: "Get the current server time in RFC1123 format",
	}, GetCurrentTimeTool)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_old_downloads",
		Description: "List files in the Download directory that haven't been modified in a long time.",
	}, ListOldDownloadsTool)

	log.Println("MCP server started (version:", version, "commit:", commit, "date:", date, "builtBy:", builtBy+")")

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
