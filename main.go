package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("Base Framework MCP Server v1.0.0")
		return
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"Base Framework",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add Base Framework tools
	addBaseFrameworkTools(mcpServer)

	// Check if running in HTTP/SSE mode (for Caprover deployment)
	if port := os.Getenv("PORT"); port != "" {
		// HTTP/SSE mode for web deployment
		log.Printf("Starting Base Framework MCP Server on port %s", port)
		
		// Determine the correct base URL for the SSE server
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = fmt.Sprintf("http://localhost:%s", port)
		}
		
		// Check if documentation routes should be enabled
		enableDocs := os.Getenv("ENABLE_DOCS")
		if enableDocs == "" {
			enableDocs = "true" // enabled by default
		}
		
		log.Printf("Server ready at %s", baseURL)
		log.Printf("- MCP SSE endpoint: %s/sse", baseURL)
		
		// Create HTTP server with custom mux
		mux := http.NewServeMux()
		
		// Add documentation routes if enabled
		if enableDocs == "true" {
			log.Printf("- Documentation: %s/docs/", baseURL)
			log.Printf("- Landing page: %s/", baseURL)
			setupHTTPRoutes(mux, baseURL)
		}
		
		// Create SSE server and integrate with custom mux
		sseServer := server.NewSSEServer(mcpServer,
			server.WithBaseURL(baseURL),
		)
		
		// Add SSE handler to our mux (using default paths)
		mux.Handle("/sse", sseServer.SSEHandler())
		mux.Handle("/message", sseServer.MessageHandler())
		
		// Start unified HTTP server
		httpServer := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}
		
		log.Printf("Starting unified HTTP server...")
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP Server error: %v\n", err)
		}
	} else {
		// stdio mode for local MCP clients
		log.Println("Starting Base Framework MCP Server in stdio mode")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Stdio Server error: %v\n", err)
		}
	}
}

func addBaseFrameworkTools(s *server.MCPServer) {
	// CLI Command Context tool
	cmdContextTool := mcp.NewTool(
		"base_cmd_context",
		mcp.WithDescription("Get context and information about Base Framework CLI commands"),
	)
	s.AddTool(cmdContextTool, handleCmdContextNew)

	// Framework Documentation tool
	frameworkDocsTool := mcp.NewTool(
		"base_framework_docs",
		mcp.WithDescription("Get Base Framework documentation, architecture, and development patterns"),
	)
	s.AddTool(frameworkDocsTool, handleFrameworkDocsNew)

	// Specific Doc File tool - allows accessing any specific documentation file
	docFileTool := mcp.NewTool(
		"base_doc_file",
		mcp.WithDescription("Get specific Base Framework documentation file by name"),
		mcp.WithString("filename", 
			mcp.Description("Name of the documentation file (e.g. 'websocket.md', 'auth.md', 'storage.md')"),
			mcp.Required(),
		),
	)
	s.AddTool(docFileTool, handleDocFileNew)
}

// readMarkdownFile reads a markdown file from the docs directory
func readMarkdownFile(filename string) (string, error) {
	docsDir := "./md/docs"
	filePath := filepath.Join(docsDir, filename)
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %v", filename, err)
	}
	
	return string(content), nil
}

// combineMarkdownFiles combines multiple markdown files
func combineMarkdownFiles(files []string) string {
	var combined strings.Builder
	
	for i, filename := range files {
		content, err := readMarkdownFile(filename)
		if err != nil {
			log.Printf("Warning: Could not read %s: %v", filename, err)
			continue
		}
		
		// Remove YAML frontmatter if present
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				content = parts[2]
			}
		}
		
		combined.WriteString(content)
		if i < len(files)-1 {
			combined.WriteString("\n\n---\n\n")
		}
	}
	
	return combined.String()
}

// New API handlers
func handleCmdContextNew(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleCmdContext()
}

func handleFrameworkDocsNew(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleFrameworkDocs(nil)
}

func handleDocFileNew(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filename := mcp.ParseString(request, "filename", "")
	if filename == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent("Error: filename parameter is required"),
			},
		}, nil
	}
	
	arguments := map[string]any{"filename": filename}
	return handleDocFile(arguments)
}

func handleCmdContext() (*mcp.CallToolResult, error) {
	// Read the CLI documentation directly from file
	content, err := readMarkdownFile("cli.md")
	if err == nil {
		// Remove YAML frontmatter if present
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				content = strings.TrimSpace(parts[2])
			}
		}
	} else {
		// Fallback to embedded content if file read fails
		content = `# Base Framework CLI Commands

## Core Commands Reference

### base new [project_name]
**Purpose**: Creates a new Base Framework project by downloading the latest framework template and setting up the project structure.

**What it does**:
- Downloads the latest Base Framework template from GitHub
- Creates a new directory with your project name
- Sets up the complete project structure
- Initializes Go modules and dependencies
- Configures development environment

**Example Usage**:
- base new my-api - Creates new project called "my-api"
- cd my-api && base start - Navigate and start development server

---

### base g/generate [module_name] [field:type...] [options]
**Purpose**: Generates complete modules with models, controllers, services, validators, and auto-updates your app initialization. The most powerful command in the Base CLI.

**What gets generated**:
- **Model**: app/models/{name}.go - GORM model with relationships
- **Controller**: app/{name}/controller.go - HTTP handlers and routes
- **Service**: app/{name}/service.go - Business logic layer
- **Module**: app/{name}/module.go - Module definition and wiring
- **Validator**: app/{name}/validator.go - Request validation
- **Registration**: Updates app/init.go automatically

**Field Types Reference**:

**Basic Data Types**:
- string - Text field (VARCHAR) - name:string
- text - Large text field (TEXT) - content:text
- email - Email field (VARCHAR) - email:email
- url - URL field (VARCHAR) - website:url
- slug - URL-friendly string (VARCHAR) - slug:slug

**Numeric Types**:
- int - Signed integers - age:int
- uint - IDs, positive numbers - user_id:uint
- float - Prices, measurements - price:float
- decimal - Financial calculations - amount:decimal
- sort - Sort order - order:sort
- bool - True/false flags - active:bool

**Date & Time Types**:
- datetime - Full date and time - created_at:datetime
- time - Timestamp - started_at:time
- date - Date only - birth_date:date
- timestamp - Unix timestamp - modified:timestamp

**File & Media Types**:
- image - Image field with 5MB limit, validation for image formats - avatar:image
- file - File field with 50MB limit, validation for documents - document:file

**Special & Data Types**:
- json - JSON data - metadata:json
- jsonb - Binary JSON (PostgreSQL) - settings:jsonb
- translation - Multi-language text - title:translation
- translatedField - i18n field - name:translatedField

**Relationship Types**:
- field:belongsTo:Model - Belongs to relationship - author:belongsTo:User
- field:hasOne:Model - Has one relationship - profile:hasOne:Profile
- field:hasMany:Model - Has many relationship - posts:hasMany:Post
- field:manyToMany:Model - Many to many relationship - tags:manyToMany:Tag

**Example Usage**:
- base g post title:string content:text published:bool
- base g product name:string price:float category_id:uint image:image
- base g article title:string content:text author:belongsTo:User category:belongsTo:Category tags:manyToMany:Tag

---

### base d/destroy [module_name1] [module_name2] ...
**Purpose**: Safely removes generated modules and cleans up all associated files, models, and registrations.

**What gets removed**:
- **Module Directory**: app/{module_name}/ and all contents
- **Model File**: app/models/{module_name}.go
- **Test Files**: test/app_test/{module_name}_test/
- **Registration**: Import and initialization from app/init.go

**Safety Features**:
- Confirmation prompts before destroying modules
- Batch support for multiple modules
- Orphan cleanup removes registry entries even if files are missing
- Status reporting shows success/failure for each operation

---

### base start [flags]
**Purpose**: Starts your Base Framework application server with optional documentation generation.

**Flags**:
- --docs/-d: Generate Swagger documentation before starting

**What it does**:
- Ensures dependencies are up to date with go mod tidy
- Validates project structure (looks for main.go)
- Generates Swagger documentation (if -d flag used)
- Starts the application server
- Sets environment variables for documentation

---

### base docs [flags]
**Purpose**: Generates Swagger 2.0 documentation using go-swagger by scanning controller annotations.

**Flags**:
- --output/-o: Output directory for generated files (default: "docs")
- --static/-s: Generate static swagger files (default: true)
- --no-static: Skip generating static files

**Generated Files**:
- swagger.json - Swagger 2.0 specification in JSON format
- swagger.yaml - Swagger 2.0 specification in YAML format
- docs.go - Go package with embedded documentation

---

### base update
**Purpose**: Updates the Base Framework core directory to the latest version while preserving your application code.

**What it does**:
- Downloads the latest Base Framework core matching your CLI version
- Creates automatic backup of existing core directory
- Replaces core directory with updated version
- Preserves your app directory and custom code
- Rollback support if update fails

---

### base upgrade [flags]
**Purpose**: Upgrades the Base CLI tool itself to the latest version.

**Flags**:
- --major: Allow upgrade to new major version (may contain breaking changes)

**Upgrade Behavior**:
- Minor/Patch: Upgrades automatically within same major version
- Major: Requires --major flag and user confirmation
- Validation: Downloads, verifies, and tests binary before installation

## Development Workflow Commands

1. **Project Setup**: base new my-project && cd my-project
2. **Generate Modules**: base g user email:string name:string
3. **Start Development**: base start -d
4. **Access & Test**: Visit http://localhost:8100/swagger/
5. **Deploy**: Build with go build -o main

This context provides comprehensive information about all Base Framework CLI commands and their usage patterns based on the official documentation.`
	}

	return &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			mcp.NewTextContent(content),
		},
	}, nil
}

func handleFrameworkDocs(arguments map[string]any) (*mcp.CallToolResult, error) {
	// Read multiple documentation files for comprehensive coverage
	docFiles := []string{
		"structure.md",
		"application.md", 
		"websocket.md",
		"router.md",
		"storage.md",
		"emitter.md",
		"auth.md",
		"middleware.md",
		"logger.md",
		"email.md",
		"configuration.md",
	}
	
	content := combineMarkdownFiles(docFiles)
	if content == "" {
		// Fallback to embedded content if files can't be read
		content = `# Base Framework Architecture & Development

## Framework Overview

Base Framework is a modern Go web framework built for rapid development with clean architecture principles. The framework follows an organized project structure with clear separation between your business logic (app/) and framework infrastructure (core/).

## Project Structure

### Base Project Structure
app/                          # Your Application Layer
├── models/                   # Database models (GORM) - centralized to prevent circular dependencies
├── <module>/                 # Each feature gets its own module directory
│   ├── controller.go         # HTTP handlers
│   ├── service.go           # Business logic
│   └── module.go            # Module registration
└── init.go                  # App initialization and module registration

core/                        # Base Framework Core
├── app/                     # Built-in app structures
│   ├── authentication/      # Auth system
│   ├── authorization/       # Authorization
│   ├── media/              # Media handling
│   ├── oauth/              # OAuth providers
│   └── profile/            # User profiles
├── base/                   # Base controller/service
├── config/                 # Configuration
├── database/               # Database connection
├── email/                  # Email providers
├── emitter/                # Event system
├── errors/                 # Error handling
├── helper/                 # Helper utilities
├── http/                   # HTTP router
├── logger/                 # Logging system
├── module/                 # Module system
├── router/                 # Router & middleware
├── storage/                # File storage (local/S3/R2)
├── translation/            # Internationalization
├── types/                  # Common types
├── validator/              # Validation system
└── websocket/              # WebSocket support

storage/                    # Active Local Storage
├── app/                    # Uploaded files
├── logs/                   # Application logs
└── temp/                   # Temporary files

## MCS Architecture Pattern

Base Framework implements the Model-Controller-Service (MCS) architecture pattern to ensure clean separation of concerns, maintainability, and testability.

### Model (Data Layer)
**Location**: app/models/
**Responsibilities**:
- GORM struct definitions
- Request/Response types
- Data validation tags
- Database relationships
- Serialization methods

### Controller (HTTP Layer)
**Location**: app/[module]/controller.go
**Responsibilities**:
- Route definitions
- Request parsing
- Response formatting
- HTTP status codes
- Input validation

### Service (Business Layer)
**Location**: app/[module]/service.go
**Responsibilities**:
- Business rules & validation
- Database operations
- External API calls
- Event emission
- Inter-service communication

## Module Initialization Flow

1. **Application Startup**: main.go creates core infrastructure (Database, Router, Logger, Emitter, Storage, Email, WebSocket, Translation, Validator, Helper, Error handling, Configuration)

2. **Core Modules First**: Framework initializes core modules (authentication, permissions, etc.) before app modules

3. **App Module Discovery**: Framework calls app/init.go:GetAppModules() to discover business modules

4. **Module Lifecycle**: For each module: Initialization → Migration → Routes → Dependencies

## Dependency Injection

Base Framework uses dependency injection to provide modules with access to framework services:

**Core Dependencies**:
- DB - GORM database instance
- Router - HTTP router for defining endpoints
- Logger - Structured logging service
- Config - Application configuration

**Extended Services**:
- Emitter - Event system for module communication
- Storage - File upload and management
- EmailSender - Email service integration
- Translation - Internationalization service
- Validator - Validation service

## Module Communication

### Event-Driven Communication (Asynchronous)
- emitter.Emit("user.created", user)
- emitter.Emit("post.published", post)
- emitter.On("user.created", handleUserCreated)

### Direct Service Injection (Synchronous)
Inject services directly for synchronous operations between modules.

## WebSocket System

### Real-Time Features
- **Endpoint**: /api/ws
- **Room-Based Messaging**: Clients join specific rooms for scoped communication
- **User Management**: Real-time user presence and notifications
- **Message Types**: message, system, users_update, draw, cursor_move
- **Connection Management**: Automatic cleanup and reconnection handling

### WebSocket Examples
Base includes complete example applications:
1. **Chat Application**: Multi-room chat with user presence
2. **Drawing Canvas**: Collaborative drawing with real-time sync
3. **Kanban Board**: Real-time task management
4. **Code Editor**: Collaborative code editing
5. **System Monitor**: Live system metrics dashboard
6. **Spreadsheet**: Real-time collaborative spreadsheet

### Connection Example
const socket = new WebSocket('ws://localhost:8100/api/ws?id=user123&nickname=John&room=general');

## Core Systems Integration

### Router System
- Custom radix tree implementation
- Zero external dependencies
- Built-in middleware support
- Static file serving

### Database Layer
- GORM integration with relationship support
- Auto-migration on startup
- Multi-database support (SQLite, PostgreSQL, MySQL)
- Connection pooling

### Storage System
- Multi-provider support (local filesystem, Amazon S3, Cloudflare R2)
- Automatic file validation based on field types
- Event integration for file operations
- Context-aware operations

### Event System
- Thread-safe asynchronous event handling
- Context support for proper propagation
- Module decoupling through events
- Built-in lifecycle events

### Email System
- Multi-provider support (SMTP, SendGrid, Postmark)
- HTML and text email templates
- Dependency injection integration
- Environment-based provider selection

### Error Handling
- Structured errors with typed error codes
- Context preservation and tracking
- HTTP integration with automatic status code mapping
- Structured error logging integration

## File Storage & Upload System

### Upload Field Types
- **image**: 5MB limit, image format validation (jpg, png, gif, webp)
- **file**: 50MB limit, document format validation (pdf, doc, txt, etc.)
- **attachment**: 10MB limit, mixed file type support

### Storage Events
- Upload Events: {module}.{field}.uploaded with file metadata
- Delete Events: {module}.{field}.deleted with cleanup confirmation
- Error Events: Storage operation error handling

## Development Patterns

### Module Registration Process
1. Core Initialization - Framework starts up core services
2. App Init - app/init.go calls RegisterModule() for each app module
3. Dependency Injection - Core services injected into module constructors
4. Route Registration - Modules register their HTTP routes with core router

### Best Practices
- Keep controllers thin - delegate to services
- Use events for cross-module communication
- Place shared models in app/models/
- Follow the generated module structure
- Use dependency injection for testability
- Emit events for important business actions

This comprehensive documentation covers the entire Base Framework architecture and development patterns based on the official documentation.`
	}

	return &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			mcp.NewTextContent(content),
		},
	}, nil
}

func handleDocFile(arguments map[string]any) (*mcp.CallToolResult, error) {
	filename, ok := arguments["filename"].(string)
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent("Error: filename parameter is required and must be a string"),
			},
		}, nil
	}

	// Ensure the filename ends with .md
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	content, err := readMarkdownFile(filename)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error: Could not read file '%s': %v", filename, err)),
			},
		}, nil
	}

	// Remove YAML frontmatter if present
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			content = strings.TrimSpace(parts[2])
		}
	}

	return &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			mcp.NewTextContent(content),
		},
	}, nil
}


func setupHTTPRoutes(mux *http.ServeMux, docsBaseURL string) {
	// Index page explaining how to use MCP
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		// Replace placeholder with actual MCP base URL (not docs URL)
		html := strings.ReplaceAll(indexHTML, "{{BASE_URL}}", getBaseURL())
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})
	
	// Serve documentation files
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		// Extract filename from path
		filename := strings.TrimPrefix(r.URL.Path, "/docs/")
		if filename == "" {
			// List available docs
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(docsIndexHTML))
			return
		}
		
		// Ensure .md extension
		if !strings.HasSuffix(filename, ".md") {
			filename += ".md"
		}
		
		// Read and serve markdown file
		content, err := readMarkdownFile(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Documentation file '%s' not found: %v", filename, err), http.StatusNotFound)
			return
		}
		
		// Remove YAML frontmatter if present
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				content = strings.TrimSpace(parts[2])
			}
		}
		
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Write([]byte(content))
	})
}

func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}
	return baseURL
}
