package main

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Base Framework MCP Server</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; line-height: 1.6; }
        .header { text-align: center; margin-bottom: 40px; }
        .logo { font-size: 2.5em; font-weight: bold; color: #2563eb; margin-bottom: 10px; }
        .tagline { color: #6b7280; font-size: 1.1em; }
        .section { margin-bottom: 30px; padding: 20px; background: #f8fafc; border-radius: 8px; }
        .section h2 { color: #1f2937; margin-top: 0; }
        .endpoint { background: #e5f3ff; padding: 15px; border-radius: 6px; margin: 10px 0; }
        .endpoint code { background: #1f2937; color: #f3f4f6; padding: 4px 8px; border-radius: 4px; }
        .tool-list { list-style: none; padding: 0; }
        .tool-list li { background: white; margin: 8px 0; padding: 12px; border-radius: 6px; border-left: 4px solid #10b981; }
        .config-block { background: #1f2937; color: #f3f4f6; padding: 15px; border-radius: 6px; margin: 10px 0; overflow-x: auto; }
        a { color: #2563eb; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">🚀 Base Framework MCP Server</div>
        <div class="tagline">Model Context Protocol server for Base Framework CLI tools and documentation</div>
    </div>

    <div class="section">
        <h2>🎯 What is this?</h2>
        <p>This is a <strong>Model Context Protocol (MCP) Server</strong> that provides AI assistants (like Claude) with access to:</p>
        <ul class="tool-list">
            <li><strong>base_cmd_context</strong> - Complete CLI command reference with examples</li>
            <li><strong>base_framework_docs</strong> - Framework architecture and development patterns</li>
            <li><strong>base_doc_file</strong> - Access to specific documentation files</li>
        </ul>
    </div>

    <div class="section">
        <h2>📡 Available Endpoints</h2>
        <div class="endpoint">
            <strong>MCP Connection:</strong> <code>/sse</code> - Server-Sent Events endpoint for MCP clients
        </div>
        <div class="endpoint">
            <strong>Documentation:</strong> <code>/docs/</code> - Browse and access markdown documentation files
        </div>
        <div class="endpoint">
            <strong>This Page:</strong> <code>/</code> - Usage instructions and setup guide
        </div>
    </div>

    <div class="section">
        <h2>🔧 How to Use with Claude Desktop</h2>
        <p>Add this configuration to your Claude Desktop config:</p>
        <div class="config-block">{
  "mcpServers": {
    "base-framework": {
      "url": "{{BASE_URL}}/sse"
    }
  }
}</div>
        <p><strong>Config file locations:</strong></p>
        <ul>
            <li><strong>macOS:</strong> <code>~/Library/Application Support/Claude/claude_desktop_config.json</code></li>
            <li><strong>Windows:</strong> <code>%APPDATA%\Claude\claude_desktop_config.json</code></li>
        </ul>
    </div>

    <div class="section">
        <h2>🏗️ How to Use with Windsurf</h2>
        <p>Add this configuration to your Windsurf MCP config:</p>
        <div class="config-block">{
  "mcpServers": {
    "base-framework": {
      "url": "{{BASE_URL}}/sse"
    }
  }
}</div>
        <p><strong>Config file locations:</strong></p>
        <ul>
            <li><strong>macOS/Linux:</strong> <code>~/.codeium/windsurf/mcp_config.json</code></li>
            <li><strong>Windows:</strong> <code>%APPDATA%\.codeium\windsurf\mcp_config.json</code></li>
        </ul>
    </div>

    <div class="section">
        <h2>💡 Example Usage</h2>
        <p>Once connected, you can ask your AI assistant:</p>
        <ul>
            <li>"Show me how to generate a Base Framework module with relationships"</li>
            <li>"What are all the field types available in Base Framework?"</li>
            <li>"Explain the WebSocket system in Base Framework"</li>
            <li>"How do I set up authentication in Base Framework?"</li>
        </ul>
    </div>

    <div class="section">
        <h2>📚 Browse Documentation</h2>
        <p><a href="/docs/">Browse available documentation files →</a></p>
    </div>

    <div class="section">
        <h2>🌟 About Base Framework</h2>
        <p>Base Framework is a modern Go web framework for rapid development with clean architecture, built-in WebSocket support, file storage, and comprehensive tooling.</p>
        <p><a href="https://github.com/base-go/base" target="_blank">Learn more about Base Framework →</a></p>
    </div>
</body>
</html>`

const docsIndexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Base Framework Documentation</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; line-height: 1.6; }
        .header { text-align: center; margin-bottom: 40px; }
        .logo { font-size: 2em; font-weight: bold; color: #2563eb; margin-bottom: 10px; }
        .doc-list { list-style: none; padding: 0; }
        .doc-list li { background: #f8fafc; margin: 8px 0; padding: 15px; border-radius: 6px; border-left: 4px solid #10b981; }
        .doc-list a { color: #2563eb; text-decoration: none; font-weight: 500; }
        .doc-list a:hover { text-decoration: underline; }
        .description { color: #6b7280; margin-top: 5px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2563eb; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <a href="/" class="back-link">← Back to MCP Server</a>
    
    <div class="header">
        <div class="logo">📚 Base Framework Documentation</div>
    </div>

    <ul class="doc-list">
        <li>
            <a href="/docs/cli">cli.md</a>
            <div class="description">Complete CLI command reference with examples and field types</div>
        </li>
        <li>
            <a href="/docs/structure">structure.md</a>
            <div class="description">Project structure and architecture overview</div>
        </li>
        <li>
            <a href="/docs/application">application.md</a>
            <div class="description">Application layer patterns and module development</div>
        </li>
        <li>
            <a href="/docs/websocket">websocket.md</a>
            <div class="description">Real-time WebSocket system and examples</div>
        </li>
        <li>
            <a href="/docs/router">router.md</a>
            <div class="description">HTTP router system and middleware</div>
        </li>
        <li>
            <a href="/docs/storage">storage.md</a>
            <div class="description">File storage system (local, S3, R2)</div>
        </li>
        <li>
            <a href="/docs/emitter">emitter.md</a>
            <div class="description">Event-driven communication system</div>
        </li>
        <li>
            <a href="/docs/auth">auth.md</a>
            <div class="description">Authentication and authorization patterns</div>
        </li>
        <li>
            <a href="/docs/middleware">middleware.md</a>
            <div class="description">HTTP middleware system and custom middleware</div>
        </li>
        <li>
            <a href="/docs/logger">logger.md</a>
            <div class="description">Structured logging system</div>
        </li>
        <li>
            <a href="/docs/email">email.md</a>
            <div class="description">Email system with multiple providers</div>
        </li>
        <li>
            <a href="/docs/configuration">configuration.md</a>
            <div class="description">Environment configuration and settings</div>
        </li>
    </ul>

    <p style="text-align: center; margin-top: 40px; color: #6b7280;">
        Access any documentation file directly: <code>/docs/[filename]</code>
    </p>
</body>
</html>`