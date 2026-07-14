// KDSE MCP Server - Model Context Protocol server for Knowledge-Driven Software Engineering
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kdse/mcp-server/tools"
)

// MCP Protocol types
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion,omitempty"`
	Capabilities    map[string]interface{} `json:"capabilities,omitempty"`
	ClientInfo      *ClientInfo            `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// Protocol constants
const (
	ProtocolVersion = "2024-11-05"
)

func main() {
	server := NewMCPServer()
	server.Run()
}

type MCPServer struct {
	tools *tools.ToolHandler
}

func NewMCPServer() *MCPServer {
	return &MCPServer{
		tools: tools.NewToolHandler(),
	}
}

func (s *MCPServer) Run() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			s.sendError(encoder, nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(encoder, &req)
	}
}

func (s *MCPServer) handleRequest(encoder *json.Encoder, req *MCPRequest) {
	// Handle MCP protocol methods
	switch req.Method {
	case "initialize":
		s.handleInitialize(encoder, req)
	case "notifications/initialized":
		// Client has finished initialization - no response needed
		return
	case "tools/list":
		s.handleListTools(encoder, req)
	case "tools/call":
		s.handleToolCall(encoder, req)
	case "help":
		s.handleHelp(encoder, req)
	case "status":
		s.handleStatus(encoder, req)
	default:
		// Try tool name directly
		if req.ID != nil {
			s.sendError(encoder, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		}
	}
}

func (s *MCPServer) handleInitialize(encoder *json.Encoder, req *MCPRequest) {
	var params InitializeParams
	if req.Params != nil {
		json.Unmarshal(req.Params, &params)
	}

	result := map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"serverInfo": map[string]interface{}{
			"name":    "kdse-mcp-server",
			"version": "0.1.0",
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
		},
		"instructions": "KDSE MCP Server v0.1 - Provides access to Knowledge-Driven Software Engineering repository information. Available tools: help, initialize, status",
	}

	s.sendResult(encoder, req.ID, result)
}

func (s *MCPServer) handleListTools(encoder *json.Encoder, req *MCPRequest) {
	toolDefs := []map[string]interface{}{
		{
			"name":        "help",
			"description": "Returns available tools and their usage information",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "initialize",
			"description": "Returns repository initialization information including module name, version, and supported features",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "status",
			"description": "Returns current repository status information including git state, file counts, and KDSE compliance",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	s.sendResult(encoder, req.ID, map[string]interface{}{
		"tools": toolDefs,
	})
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func (s *MCPServer) handleToolCall(encoder *json.Encoder, req *MCPRequest) {
	var params ToolCallParams
	if req.Params != nil {
		json.Unmarshal(req.Params, &params)
	}

	var result interface{}
	switch params.Name {
	case "help":
		result = s.tools.Help()
	case "initialize":
		result = s.tools.Initialize()
	case "status":
		result = s.tools.Status()
	default:
		s.sendError(encoder, req.ID, -32602, fmt.Sprintf("Unknown tool: %s", params.Name))
		return
	}

	s.sendResult(encoder, req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatJSON(result),
			},
		},
	})
}

func (s *MCPServer) handleHelp(encoder *json.Encoder, req *MCPRequest) {
	result := s.tools.Help()
	s.sendResult(encoder, req.ID, result)
}

func (s *MCPServer) handleStatus(encoder *json.Encoder, req *MCPRequest) {
	result := s.tools.Status()
	s.sendResult(encoder, req.ID, result)
}

func (s *MCPServer) sendResult(encoder *json.Encoder, id interface{}, result interface{}) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	encoder.Encode(resp)
}

func (s *MCPServer) sendError(encoder *json.Encoder, id interface{}, code int, message string) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	encoder.Encode(resp)
}

func formatJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"failed to format JSON: %v\"}", err)
	}
	return string(b)
}
