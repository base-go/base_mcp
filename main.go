package main

import (
	"context"
	"fmt"
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

	// Check if running in web mode (with PORT env var) or local stdio mode
	if port := os.Getenv("PORT"); port != "" {
		// Web mode - serve installer page
		log.Printf("Starting web installer server on port %s", port)
		
		mux := http.NewServeMux()
		
		// Serve installer page
		mux.HandleFunc("/", serveInstaller)
		mux.HandleFunc("/install", serveInstallScript)
		mux.HandleFunc("/releases/", serveReleases)
		
		log.Printf("Installer available at: http://localhost:%s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("HTTP Server error: %v", err)
		}
	} else {
		// Local stdio mode for editor integration
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

func serveInstaller(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Base Framework MCP Server</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .code { background: #f4f4f4; padding: 15px; border-radius: 5px; font-family: monospace; margin: 10px 0; }
        .step { margin: 30px 0; }
        h1 { color: #333; }
        h2 { color: #666; }
    </style>
</head>
<body>
    <h1>🚀 Base Framework MCP Server</h1>
    <p>Model Context Protocol server for Base Framework - integrates with Claude Code to provide Base Framework documentation and CLI context.</p>
    
    <div class="step">
        <h2>📦 Installation</h2>
        <p>Install the Base MCP server to your <code>~/.base</code> directory:</p>
        <div class="code">curl -fsSL https://mcp.base.al/install | bash</div>
    </div>
    
    <div class="step">
        <h2>⚙️ Configuration</h2>
        <p>Add to your Claude Code configuration file (<code>~/.claude.json</code> or <code>%APPDATA%\.claude.json</code>):</p>
        <div class="code">{
  "mcpServers": {
    "base": {
      "command": "~/.base/base-mcp",
      "args": []
    }
  }
}</div>
    </div>
    
    <div class="step">
        <h2>✅ Verification</h2>
        <p>Test the installation:</p>
        <div class="code">claude mcp list</div>
        <p>You should see:</p>
        <div class="code">base: ~/.base/base-mcp (stdio) - ✓ Connected
  Tools: base_info</div>
    </div>
    
    <div class="step">
        <h2>🔧 Usage</h2>
        <p>Once configured, you can use the Base Framework MCP server in Claude Code:</p>
        <ul>
            <li><strong>base_info</strong>: Get Base Framework information and documentation</li>
        </ul>
    </div>
    
    <div class="step">
        <h2>🌐 Links</h2>
        <ul>
            <li><a href="https://github.com/yourusername/base-framework">Base Framework GitHub</a></li>
            <li><a href="https://claude.ai/code">Claude Code</a></li>
            <li><a href="https://spec.modelcontextprotocol.io">MCP Specification</a></li>
        </ul>
    </div>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func serveInstallScript(w http.ResponseWriter, r *http.Request) {
	script := `#!/bin/bash
set -e

# Base MCP Server Installer
echo "🚀 Installing Base Framework MCP Server..."

# Create ~/.base directory
BASE_DIR="$HOME/.base"
mkdir -p "$BASE_DIR"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

# Download URL (you'll need to build and host binaries)
BINARY_URL="https://mcp.base.al/releases/base-mcp-$OS-$ARCH"

# Download binary
echo "📥 Downloading base-mcp for $OS-$ARCH..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$BINARY_URL" -o "$BASE_DIR/base-mcp"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$BINARY_URL" -O "$BASE_DIR/base-mcp"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
chmod +x "$BASE_DIR/base-mcp"

echo "✅ Base MCP Server installed to $BASE_DIR/base-mcp"
echo ""
echo "📝 Next steps:"
echo "1. Add to your Claude Code configuration (~/.claude.json):"
echo '   {'
echo '     "mcpServers": {'
echo '       "base": {'
echo '         "command": "~/.base/base-mcp",'
echo '         "args": []'
echo '       }'
echo '     }'
echo '   }'
echo ""
echo "2. Test with: claude mcp list"
echo ""
echo "🎉 Installation complete!"
`
	
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=install.sh")
	fmt.Fprint(w, script)
}

func serveReleases(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := r.URL.Path[len("/releases/"):]
	
	// Basic security - only allow expected binary names
	if !isValidBinaryName(filename) {
		http.NotFound(w, r)
		return
	}
	
	// Serve the binary file from releases directory
	http.ServeFile(w, r, fmt.Sprintf("releases/%s", filename))
}

func isValidBinaryName(filename string) bool {
	validNames := []string{
		"base-mcp-darwin-amd64",
		"base-mcp-darwin-arm64", 
		"base-mcp-linux-amd64",
		"base-mcp-linux-arm64",
	}
	
	for _, valid := range validNames {
		if filename == valid {
			return true
		}
	}
	return false
}