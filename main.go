package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed md
var docsFS embed.FS

func main() {
	// Create simple MCP server
	mcpServer := server.NewMCPServer("Base Framework", "1.0.0")

	// Add Base Framework tools
	infoTool := mcp.NewTool("base_info", mcp.WithDescription("Get Base Framework information"))
	mcpServer.AddTool(infoTool, handleBaseInfo)

	cliTool := mcp.NewTool("base_cli", mcp.WithDescription("Get Base Framework CLI commands and usage"))
	mcpServer.AddTool(cliTool, handleBaseCLI)

	docsTool := mcp.NewTool("base_docs", mcp.WithDescription("Get Base Framework documentation and features"))
	mcpServer.AddTool(docsTool, handleBaseDocs)

	// Check if running in web mode (with PORT env var) or local stdio mode
	if port := os.Getenv("PORT"); port != "" {
		// Web mode - serve installer page
		log.Printf("Starting web installer server on port %s", port)

		mux := http.NewServeMux()

		// Serve installer page
		mux.HandleFunc("/", serveInstaller)
		mux.HandleFunc("/install", serveInstallScript)

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
			mcp.NewTextContent("Base Framework is a modern Go web framework for rapid development with intelligent code generation, smart relationship detection, and production-ready features out of the box."),
		},
	}, nil
}

func handleBaseCLI(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Read CLI documentation from markdown file
	content, err := readMarkdownFile("md/docs/cli.md")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Error reading CLI documentation: " + err.Error()),
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(content),
		},
	}, nil
}

func handleBaseDocs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Combine multiple documentation files for comprehensive docs
	var allDocs strings.Builder

	// Main overview from index.md
	if content, err := readMarkdownFile("md/index.md"); err == nil {
		allDocs.WriteString(content)
		allDocs.WriteString("\n\n")
	}

	// Core modules documentation
	coreModules := []string{
		"router", "emitter", "storage", "middleware",
		"logger", "websocket", "auth", "email",
	}

	allDocs.WriteString("# Core Framework Modules\n\n")

	for _, module := range coreModules {
		if content, err := readMarkdownFile(fmt.Sprintf("md/docs/%s.md", module)); err == nil {
			allDocs.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(module)))
			allDocs.WriteString(content)
			allDocs.WriteString("\n\n")
		}
	}

	// Add configuration docs
	if content, err := readMarkdownFile("md/docs/configuration.md"); err == nil {
		allDocs.WriteString("# Configuration\n\n")
		allDocs.WriteString(content)
		allDocs.WriteString("\n\n")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(allDocs.String()),
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
        <p><strong>Option 1:</strong> Add as global MCP server (recommended):</p>
        <div class="code">claude mcp add --scope user base ~/.base/base-mcp</div>
        <p><strong>Option 2:</strong> Manual configuration in <code>~/.claude.json</code>:</p>
        <div class="code">{
  "mcpServers": {
    "base": {
      "type": "stdio",
      "command": "~/.base/base-mcp",
      "args": [],
      "env": {}
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
        <p>Once configured, you can use these Base Framework tools in Claude Code:</p>
        <ul>
            <li><strong>base_info</strong>: Get Base Framework overview and information</li>
            <li><strong>base_cli</strong>: Complete CLI commands reference with examples</li>
            <li><strong>base_docs</strong>: Comprehensive framework documentation and features</li>
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

# Get latest release from GitHub  
REPO="base-go/base"
echo "📡 Fetching latest release info..."
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"

if command -v curl >/dev/null 2>&1; then
    RELEASE_INFO=$(curl -fsSL "$LATEST_URL")
elif command -v wget >/dev/null 2>&1; then
    RELEASE_INFO=$(wget -qO- "$LATEST_URL")
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Extract download URL for the specific binary
BINARY_NAME="base-mcp-$OS-$ARCH"
DOWNLOAD_URL=$(echo "$RELEASE_INFO" | grep -o "https://github.com/$REPO/releases/download/[^\"]*/$BINARY_NAME" | head -1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Could not find binary for $OS-$ARCH in latest release"
    echo "Available binaries:"
    echo "$RELEASE_INFO" | grep -o "base-mcp-[^\"]*" | sort | uniq
    exit 1
fi

# Download binary
echo "📥 Downloading $BINARY_NAME..."
if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$DOWNLOAD_URL" -o "$BASE_DIR/base-mcp"
else
    wget -q "$DOWNLOAD_URL" -O "$BASE_DIR/base-mcp"
fi

# Make executable
chmod +x "$BASE_DIR/base-mcp"

# Verify installation
if "$BASE_DIR/base-mcp" --help >/dev/null 2>&1 || [ $? -eq 0 ]; then
    echo "✅ Base MCP Server installed successfully to $BASE_DIR/base-mcp"
else
    echo "⚠️  Installation completed but binary verification failed"
fi

echo ""
echo "📝 Next steps:"
echo "1. Add as global MCP server (recommended):"
echo "   claude mcp add --scope user base ~/.base/base-mcp"
echo ""
echo "2. Or add manually to ~/.claude.json:"
echo '   {'
echo '     "mcpServers": {'
echo '       "base": {'
echo '         "command": "~/.base/base-mcp",'
echo '         "args": []'
echo '       }'
echo '     }'
echo '   }'
echo ""
echo "3. Verify installation:"
echo "   claude mcp list"
echo ""
echo "🎉 Installation complete!"
echo ""
echo "Available tools:"
echo "- base_info: Get Base Framework information"
echo "- base_cli: Base CLI commands and usage"
echo "- base_docs: Complete framework documentation"
`

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=install.sh")
	fmt.Fprint(w, script)
}

func readMarkdownFile(filePath string) (string, error) {
	// Read from embedded filesystem
	content, err := docsFS.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded file %s: %w", filePath, err)
	}

	// Remove YAML frontmatter if present
	contentStr := string(content)
	if strings.HasPrefix(contentStr, "---") {
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) >= 3 {
			contentStr = strings.TrimSpace(parts[2])
		}
	}

	return contentStr, nil
}
