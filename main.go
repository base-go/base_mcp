package main

import (
	"context"
	"log"
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
		// Try stdio over HTTP approach
		log.Printf("Starting stdio over HTTP on port %s", port)
		if err := server.ServeHTTP(mcpServer, ":"+port); err != nil {
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