# Base Framework MCP Server

A unified Model Context Protocol (MCP) server that provides access to Base Framework CLI tools and comprehensive documentation. The server can run locally for development or be deployed to the cloud for remote access.

## 🎯 Overview

This MCP server provides complete access to Base Framework capabilities:

- **🛠️ 3 MCP Tools**: CLI commands, framework docs, and specific doc files
- **📚 Real Documentation**: Reads directly from Base Framework markdown files
- **🔄 Always Fresh**: Updates automatically when docs change
- **🌐 Dual Mode**: Local stdio and remote HTTP/SSE deployment
- **☁️ Cloud Ready**: Easy deployment with Caprover, Docker, and more

## 🛠️ Available Tools

### 1. `base_cmd_context`
Complete CLI command reference with all field types, relationship patterns, and real-world examples.

### 2. `base_framework_docs` 
Comprehensive framework documentation including architecture, patterns, WebSocket, storage, and more.

### 3. `base_doc_file`
Access specific documentation files by name (e.g., `websocket.md`, `auth.md`, `storage.md`).

## 🚀 Installation & Deployment

### 🏠 Local Development

#### 1. Build from Source
```bash
# Clone and build
git clone <repository>
cd base_mcp
make build

# Test the server
./base-mcp --version
```

#### 2. Claude Desktop Integration

**macOS:**
```bash
# Edit Claude Desktop config
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**Windows:**
```cmd
notepad %APPDATA%\Claude\claude_desktop_config.json
```

**Configuration:**
```json
{
  "mcpServers": {
    "base-framework": {
      "command": "/absolute/path/to/base-mcp/base-mcp",
      "args": [],
      "env": {}
    }
  }
}
```

#### 3. Windsurf Integration

**Configure Windsurf:**
```bash
# macOS/Linux
code ~/.codeium/windsurf/mcp_config.json

# Windows
code %APPDATA%\.codeium\windsurf\mcp_config.json
```

**Configuration:**
```json
{
  "mcpServers": {
    "base-framework": {
      "command": "/absolute/path/to/base-mcp/base-mcp"
    }
  }
}
```

### ☁️ Cloud Deployment

#### 🚢 Caprover Deployment (Recommended)

1. **Setup Caprover CLI:**
```bash
npm install -g caprover
```

2. **Deploy with Script:**
```bash
# Configure your Caprover URL
export CAPROVER_URL="your-caprover-instance.com"
export APP_URL="https://mcp.base.al"

# Run deployment
./deploy.sh
```

3. **Manual Deployment:**
```bash
# Login to Caprover
caprover login

# Create and deploy app
caprover apps:create base-mcp
caprover deploy --appName base-mcp

# Enable HTTPS
caprover apps:ssl:enable base-mcp --letsencrypt
```

4. **Claude Desktop Config (Remote):**
```json
{
  "mcpServers": {
    "base-framework": {
      "url": "https://mcp.base.al/sse"
    }
  }
}
```

#### 🐳 Docker Deployment

**Local Testing:**
```bash
# Build and run with Docker Compose
docker-compose up -d

# Test the deployment
curl http://localhost:8080/sse
```

**Production Deployment:**
```bash
# Build image
docker build -t base-framework-mcp .

# Run container
docker run -d \
  --name base-mcp \
  -p 8080:8080 \
  -e PORT=8080 \
  -e BASE_URL=https://your-domain.com \
  base-framework-mcp
```

#### ☁️ Other Cloud Platforms

**Heroku:**
```bash
# Create Heroku app
heroku create base-framework-mcp

# Set environment variables
heroku config:set PORT=80
heroku config:set BASE_URL=https://base-framework-mcp.herokuapp.com

# Deploy
git push heroku main
```

**Railway:**
```bash
# Install Railway CLI
npm install -g @railway/cli

# Deploy
railway login
railway deploy
```

**Fly.io:**
```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Initialize and deploy
fly apps create base-framework-mcp
fly deploy
```

**DigitalOcean App Platform:**
1. Connect your GitHub repository
2. Set environment variables: `PORT=8080`
3. Deploy automatically on git push

**AWS/GCP/Azure:**
- Use Docker image with container services
- Set `PORT` and `BASE_URL` environment variables
- Enable HTTPS/SSL certificates

### 🌐 DNS & Domain Setup

For `mcp.base.al` or custom domain:

1. **Point DNS to your deployment:**
   - Caprover: Point to your Caprover instance
   - Heroku: `CNAME` to `yourapp.herokuapp.com`
   - Other: Point to deployment IP/URL

2. **SSL Configuration:**
   - Most platforms provide automatic SSL
   - For custom setups, use Let's Encrypt

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP port for server | stdio mode |
| `BASE_URL` | Public URL for the server | `http://localhost:PORT` |
| `ENABLE_DOCS` | Enable documentation routes | `true` |

### Deployment Modes

- **No `PORT` env**: Runs in stdio mode (local MCP)
- **With `PORT` env**: Runs unified HTTP server with:
  - `/sse` - MCP Server-Sent Events endpoint
  - `/` - Usage guide and setup instructions (if `ENABLE_DOCS=true`)
  - `/docs/` - Documentation browser and files (if `ENABLE_DOCS=true`)
- **Documentation Routes**: 
  - Enabled by default when `PORT` is set
  - Served on same port as MCP endpoint
  - Disable with `ENABLE_DOCS=false` for production deployments

## 🧪 Testing

### Local Testing
```bash
# Test version
./base-mcp --version

# Test MCP communication
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{}}}' | ./base-mcp

# Test tools
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./base-mcp

# Test specific tool
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"base_doc_file","arguments":{"filename":"websocket.md"}}}' | ./base-mcp
```

### Remote Testing
```bash
# Test unified server endpoints
curl https://mcp.base.al/sse     # MCP SSE endpoint
curl https://mcp.base.al/        # Usage guide
curl https://mcp.base.al/docs/   # Documentation browser

# Test locally with all features
PORT=8080 ./base-mcp &
curl http://localhost:8080/sse   # MCP endpoint
curl http://localhost:8080/      # Usage guide  
curl http://localhost:8080/docs/ # Documentation

# Test with docs disabled (production mode)
PORT=8080 ENABLE_DOCS=false ./base-mcp &
curl http://localhost:8080/sse   # Only MCP endpoint available
curl http://localhost:8080/      # Returns 404
```

## 📖 Usage Examples

### In Claude Desktop

Once configured, you can ask Claude:

```
Show me how to generate a Base Framework module with relationships

What are all the field types available in Base Framework?

Explain the WebSocket system in Base Framework

How do I set up authentication in Base Framework?
```

### API Usage

**CLI Commands:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "base_cmd_context",
    "arguments": {}
  }
}
```

**Framework Architecture:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "base_framework_docs", 
    "arguments": {}
  }
}
```

**Specific Documentation:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "base_doc_file",
    "arguments": {
      "filename": "websocket.md"
    }
  }
}
```

## 🛠️ Development

### Build Commands
```bash
make build      # Build binary
make dev        # Development build with race detection
make test       # Run tests
make fmt        # Format code
make lint       # Lint code (requires golangci-lint)
make clean      # Clean build artifacts
```

### File Structure
```
base_mcp/
├── main.go           # MCP server implementation
├── md/docs/          # Base Framework documentation
├── Dockerfile        # Container configuration
├── captain-definition # Caprover deployment config
├── docker-compose.yml # Local Docker setup
├── deploy.sh         # Caprover deployment script
└── README.md         # This file
```

## 🐛 Troubleshooting

### Common Issues

**Local MCP not working:**
- Ensure absolute paths in configuration
- Check file permissions: `chmod +x base-mcp`
- Restart Claude Desktop/Windsurf after config changes

**Remote deployment not accessible:**
- Check environment variables (`PORT`, `BASE_URL`)
- Verify DNS configuration
- Check SSL certificate status
- Review deployment logs

**Documentation not loading:**
- Ensure `md/docs/` directory is included in deployment
- Check file permissions in container
- Verify markdown files exist

### Debug Steps

1. **Test locally first:** `./base-mcp --version`
2. **Check configuration syntax:** Validate JSON config files
3. **Review logs:** Check deployment platform logs
4. **Test endpoints:** Use curl to test HTTP deployment
5. **Verify tools:** Use echo commands to test MCP communication

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

## 📄 License

MIT License - see Base Framework license for details.

---

**🌟 Ready to use Base Framework MCP Server everywhere - from local development to global cloud deployment!**