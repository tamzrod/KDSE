// KDSE MCP Server - Model Context Protocol server for Knowledge-Driven Software Engineering
// Supports both STDIO (local development) and Streamable HTTP (remote deployment) transports
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kdse/mcp/tools"
)

// =============================================================================
// Transport Configuration
// =============================================================================

const (
	TransportStdio = "stdio"
	TransportHTTP  = "http"

	DefaultHTTPPort = "8080"
	ProtocolVersion = "2024-11-05"
)

// Transport is selected via MCP_TRANSPORT environment variable or --transport flag
func getTransport() string {
	if transport := os.Getenv("MCP_TRANSPORT"); transport != "" {
		return strings.ToLower(transport)
	}
	// Check legacy MCP_STDIO env var for backward compatibility
	if os.Getenv("MCP_STDIO") == "true" {
		return TransportStdio
	}
	// Default to stdio for local development
	return TransportStdio
}

func getHTTPPort() string {
	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		return port
	}
	return DefaultHTTPPort
}

// =============================================================================
// MCP Protocol Types
// =============================================================================

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

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// =============================================================================
// Main Entry Point
// =============================================================================

func main() {
	transport := getTransport()

	switch transport {
	case TransportHTTP:
		log.Println("Starting KDSE MCP Server in HTTP mode...")
		runHTTPServer()
	default:
		log.Println("Starting KDSE MCP Server in STDIO mode...")
		runStdioServer()
	}
}

// =============================================================================
// KDSE Service Layer (shared by all transports)
// =============================================================================

type KDSEService struct {
	tools *tools.ToolHandler
}

func NewKDSEService() *KDSEService {
	return &KDSEService{
		tools: tools.NewToolHandler(),
	}
}

func (s *KDSEService) HandleRequest(req *MCPRequest) (interface{}, *MCPError) {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req), nil
	case "notifications/initialized":
		return nil, nil
	case "tools/list":
		return s.handleListTools(), nil
	case "tools/call":
		return s.handleToolCall(req)
	case "help":
		return s.tools.Help(), nil
	case "status":
		return s.tools.Status(), nil
	default:
		return nil, &MCPError{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)}
	}
}

func (s *KDSEService) handleInitialize(req *MCPRequest) map[string]interface{} {
	var params InitializeParams
	if req.Params != nil {
		json.Unmarshal(req.Params, &params)
	}

	return map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"serverInfo": map[string]interface{}{
			"name":    "kdse-mcp",
			"version": "0.2.0",
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
		},
		"instructions": "KDSE MCP Server v0.1 - Provides access to Knowledge-Driven Software Engineering repository information. Available tools: help, initialize, status",
	}
}

func (s *KDSEService) handleListTools() map[string]interface{} {
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

	return map[string]interface{}{
		"tools": toolDefs,
	}
}

func (s *KDSEService) handleToolCall(req *MCPRequest) (interface{}, *MCPError) {
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
		return nil, &MCPError{Code: -32602, Message: fmt.Sprintf("Unknown tool: %s", params.Name)}
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatJSON(result),
			},
		},
	}, nil
}

// =============================================================================
// STDIO Transport (Local Development)
// =============================================================================

func runStdioServer() {
	server := &StdioTransport{
		service: NewKDSEService(),
	}
	server.Run()
}

type StdioTransport struct {
	service *KDSEService
}

func (t *StdioTransport) Run() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			t.sendError(encoder, nil, -32700, "Parse error")
			continue
		}

		t.handleRequest(encoder, &req)
	}
}

func (t *StdioTransport) handleRequest(encoder *json.Encoder, req *MCPRequest) {
	result, err := t.service.HandleRequest(req)
	if err != nil {
		if req.ID != nil {
			t.sendError(encoder, req.ID, err.Code, err.Message)
		}
		return
	}

	// Handle notifications that don't need responses
	if result == nil && req.Method == "notifications/initialized" {
		return
	}

	if req.ID != nil {
		t.sendResult(encoder, req.ID, result)
	}
}

func (t *StdioTransport) sendResult(encoder *json.Encoder, id interface{}, result interface{}) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	encoder.Encode(resp)
}

func (t *StdioTransport) sendError(encoder *json.Encoder, id interface{}, code int, message string) {
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

// =============================================================================
// HTTP Transport (Remote Deployment)
// =============================================================================

func runHTTPServer() {
	server := &HTTPTransport{
		service: NewKDSEService(),
		port:    getHTTPPort(),
	}
	log.Printf("HTTP server listening on port %s", server.port)
	log.Fatal(http.ListenAndServe(":"+server.port, server))
}

type HTTPTransport struct {
	service *KDSEService
	port    string
}

func (t *HTTPTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for remote clients
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, MCP-Session-ID")
	w.Header().Set("Content-Type", "application/json")

	// Handle CORS preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Health endpoint
	if r.URL.Path == "/health" || r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		return
	}

	// MCP endpoint - requires POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle streaming endpoint for MCP Streamable HTTP
	if strings.HasPrefix(r.URL.Path, "/mcp/") || r.URL.Path == "/mcp" {
		t.handleMCP(w, r)
		return
	}

	// Default to MCP endpoint
	t.handleMCP(w, r)
}

func (t *HTTPTransport) handleMCP(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := t.service.HandleRequest(&req)
	if err != nil {
		resp := MCPResponse{
			JSONRPC: "2.0",
			Error:   err,
			ID:      req.ID,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Handle notifications that don't need responses
	if result == nil && req.Method == "notifications/initialized" {
		w.WriteHeader(http.StatusOK)
		return
	}

	resp := MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// =============================================================================
// Utility Functions
// =============================================================================

func formatJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"failed to format JSON: %v\"}", err)
	}
	return string(b)
}
