package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create simple MCP server
	mcpServer := server.NewMCPServer("Base Framework", "1.0.0")

	// Add one simple tool
	tool := mcp.NewTool("base_info", mcp.WithDescription("Get Base Framework information"))
	mcpServer.AddTool(tool, handleBaseInfo)

	// Check mode
	if port := os.Getenv("PORT"); port != "" {
		// SSE mode with request logging
		log.Printf("Starting SSE server with logging on port %s", port)
		
		// Create SSE server
		sseServer := server.NewSSEServer(mcpServer)
		log.Printf("SSE server created")
		
		// Create custom handler for both SSE and direct HTTP requests
		mux := http.NewServeMux()
		
		// Add SSE handler for SSE transport
		mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("SSE Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			sseServer.ServeHTTP(w, r)
		})
		
		// Add direct HTTP MCP handler for HTTP transport (simplified)
		mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("HTTP MCP Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			
			if r.Method == "POST" {
				// For now, respond with a simple error to Claude Code
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"HTTP transport not fully implemented - use SSE transport"}}`))
			} else {
				// Fallback to SSE server for other methods
				sseServer.ServeHTTP(w, r)
			}
		})
		
		log.Printf("Starting HTTP server with combined SSE/HTTP handler...")
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("HTTP Server error: %v", err)
		}
	} else {
		// stdio mode for local
		log.Println("Starting stdio mode")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Stdio error: %v", err)
		}
	}
}

func handleBaseInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent("Base Framework is a modern Go web framework for rapid development."),
		},
	}, nil
}