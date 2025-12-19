package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Server represents an MCP server
type Server struct {
	info     ServerInfo
	tools    map[string]Tool
	handlers map[string]func(map[string]interface{}) (interface{}, error)
}

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
	return &Server{
		info: ServerInfo{
			Name:    name,
			Version: version,
		},
		tools:    make(map[string]Tool),
		handlers: make(map[string]func(map[string]interface{}) (interface{}, error)),
	}
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(tool Tool, handler func(map[string]interface{}) (interface{}, error)) {
	s.tools[tool.Name] = tool
	s.handlers[tool.Name] = handler
}

// Run starts the MCP server
func (s *Server) Run() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, -32700, "Parse error", nil)
			continue
		}

		s.handleRequest(&req)
	}
}

// handleRequest processes a JSON-RPC request
func (s *Server) handleRequest(req *JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		result := InitializeResult{
			ProtocolVersion: "2024-11-05",
			Capabilities: ServerCapabilities{
				Tools: map[string]bool{},
			},
			ServerInfo: s.info,
		}
		s.sendResponse(req.ID, result)

	case "tools/list":
		tools := make([]Tool, 0, len(s.tools))
		for _, tool := range s.tools {
			tools = append(tools, tool)
		}
		s.sendResponse(req.ID, map[string]interface{}{
			"tools": tools,
		})

	case "tools/call":
		toolName, ok := req.Params["name"].(string)
		if !ok {
			s.sendError(req.ID, -32602, "Invalid tool name", nil)
			return
		}

		handler, exists := s.handlers[toolName]
		if !exists {
			s.sendError(req.ID, -32601, fmt.Sprintf("Tool not found: %s", toolName), nil)
			return
		}

		arguments, _ := req.Params["arguments"].(map[string]interface{})
		result, err := handler(arguments)
		if err != nil {
			s.sendError(req.ID, -32603, err.Error(), nil)
			return
		}

		s.sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		})

	case "ping":
		s.sendResponse(req.ID, map[string]interface{}{})

	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

// sendResponse sends a JSON-RPC response
func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.send(&resp)
}

// sendError sends a JSON-RPC error response
func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	s.send(&resp)
}

// send writes a response to stdout
func (s *Server) send(resp *JSONRPCResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s\n", data)
}
