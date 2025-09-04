#!/bin/bash
set -e

# Build script for Base MCP Server
echo "🔨 Building Base MCP Server for multiple platforms..."

# Create releases directory
mkdir -p releases

# Build configurations
declare -a platforms=(
    "darwin/amd64"
    "darwin/arm64" 
    "linux/amd64"
    "linux/arm64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -ra ADDR <<< "$platform"
    OS=${ADDR[0]}
    ARCH=${ADDR[1]}
    
    echo "Building for $OS-$ARCH..."
    
    # Set environment variables and build
    CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -a -installsuffix cgo -o "releases/base-mcp-$OS-$ARCH" .
    
    # Make executable (for unix systems)
    chmod +x "releases/base-mcp-$OS-$ARCH"
    
    echo "✅ Built releases/base-mcp-$OS-$ARCH"
done

echo ""
echo "🎉 All builds completed!"
echo "📦 Built binaries:"
ls -la releases/