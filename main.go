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
		
		// Add HTTP request logging middleware with SSE connection tracking
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("HTTP Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			log.Printf("Headers: %v", r.Header)
			
			if r.URL.Path == "/sse" && r.Header.Get("Accept") == "text/event-stream" {
				log.Printf("SSE connection attempt - keeping alive")
			}
			
			sseServer.ServeHTTP(w, r)
			log.Printf("Response completed for %s %s", r.Method, r.URL.Path)
		})
		
		log.Printf("Starting HTTP server with SSE handler...")
		if err := http.ListenAndServe(":"+port, handler); err != nil {
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