---
title: Structure
description: Structure documentation for Base Framework.
---

# Directory Structure
Understanding Base Framework's organized project structure and file conventions.
## Project Overview
### Base Project Structure
```
`my-base-project/
├── app/                          # Your Application Layer
│   ├── models/                   # Database models (GORM)
│   └── init.go                   # App initialization
├── @base/core/                   # Base Framework Core (@base/core)
│   ├── app/                      # Built-in app structures
│   │   ├── authentication/       # Auth system
│   │   ├── authorization/        # Authorization
│   │   ├── media/                # Media handling
│   │   ├── oauth/                # OAuth providers
│   │   └── profile/              # User profiles
│   ├── base/                     # Base controller/service
│   ├── config/                   # Configuration
│   ├── database/                 # Database connection
│   ├── email/                    # Email providers
│   ├── emitter/                  # Event system
│   ├── errors/                   # Error handling
│   ├── helper/                   # Helper utilities
│   ├── http/                     # HTTP router
│   ├── logger/                   # Logging system
│   ├── module/                   # Module system
│   ├── router/                   # Router & middleware
│   │   └── middleware/           # Built-in middleware
│   ├── storage/                  # File storage (local/S3/R2)
│   ├── translation/              # Internationalization
│   ├── types/                    # Common types
│   ├── validator/                # Validation system
│   └── websocket/                # WebSocket support
├── storage/                      # Active Local Storage
│   ├── app/                      # Uploaded files
│   ├── logs/                     # Application logs
│   └── temp/                     # Temporary files
├── .env                          # Environment config
├── go.mod                        # Go module definition
└── main.go                       # Application entry point
`
```
## Directory Breakdown
### app/
Application structures and models
#### app/models/
Centralized location for all database models. This prevents circular dependencies
and makes model relationships clear.
`
user.go, post.go, comment.go, category.go...
`
#### app/&lt;module&gt;/
Each feature gets its own module directory with controller, service, and module files.
controller.go
HTTP handlers
service.go
Business logic
module.go
Module registration
#### app/init.go
Central module registration and application initialization.
### core/
Framework core systems powering your app
####
router/
HTTP routing system with middleware
• auth.go, logger.go, ratelimit.go
• context.go, router.go, tree.go
####
database/
Database connection management
• database.go - GORM integration
####
storage/
File storage with multiple providers
• local.go, s3.go, r2.go
• active_storage.go, types.go
####
logger/
Structured logging system
• init.go, logger.go
####
emitter/
Event system for module communication
• emitter.go
####
email/
Email service with multiple providers
• smtp.go, sendgrid.go, postmark.go
• default.go, email.go, manager.go
####
errors/
Error handling system
• errors.go
####
module/
Module registration and management
• app.go, core.go, initializer.go
• register.go
####
config/
Configuration management
• config.go
####
validator/
Input validation system
• validator.go
####
translation/
Internationalization system
• controller.go, model.go, service.go
• helper.go, module.go
####
websocket/
WebSocket support
• websocket.go
### @base/core/app/
Built-in application structures
####
authentication/
Complete authentication system
• controller.go, service.go, mod.go
• model.go, errors.go, template.go
####
authorization/
Role-based access control
• controller.go, service.go, middleware.go
• model.go, module.go
####
media/
Media file management
• controller.go, service.go
• mod.go, model.go
####
oauth/
OAuth provider integration
• controller.go, service.go, config.go
• mod.go, model.go
####
profile/
User profile management
• controller.go, service.go
• model.go, module.go
### storage/
Runtime storage directories
#### app/
Uploaded files and user content
#### logs/
Application and access logs
#### temp/
Temporary files and caches
### static/
WebSocket Examples folder
#### chat.html
WebSocket real-time chat with rooms example
#### draw.html
WebSocket real-time multi-user draw example
#### kanban.html
WebSocket real-time kanban example
#### editor.html
WebSocket real-time multi-user editor example
#### spreadsheet.html
WebSocket real-time multi-user spreadsheet example
#### monitor.html
WebSocket system monitor example
## Core-App Integration
### Module System Integration
####
Core Framework Structures Provide
- • HTTP Router & Context (@base/core)
- • Database Connection (@base/core)
- • Logging Infrastructure (@base/core)
- • Event Emitter (@base/core)
- • Storage Systems (@base/core)
- • Error Handling (@base/core)
####
Your App Structures Consume
- • Route Registration (user code)
- • Database Models (user code)
- • Log Messages (user code)
- • Event Listeners (user code)
- • File Operations (user code)
- • Custom Errors (user code)
#### Dependency Injection Flow
```
// In app/user/module.go
import (
"@base/core/database"
"@base/core/logger"
"@base/core/emitter"
"@base/core/storage"
"@base/core/router"
)
type UserModule struct {
DB      *database.DB          // Injected from @base/core/database
Logger  *logger.Logger        // Injected from @base/core/logger
Emitter *emitter.Emitter      // Injected from @base/core/emitter
Storage *storage.Storage      // Injected from @base/core/storage
}
func (m *UserModule) Register(r *router.Router) {
// Register routes using core router
r.POST("/api/users", m.controller.Create)
r.GET("/api/users", m.controller.List)
}
```
### Module Registration Process
1
#### Core Initialization
Framework starts up router, database, logger, emitter
2
#### App Init
app/init.go calls RegisterModule() for each app module
3
#### Dependency Injection
Core services injected into module constructors
4
#### Route Registration
Modules register their HTTP routes with core router
## Framework vs Application Code
### Built-in Framework Code
Part of Base Framework
#### core/ directory
Contains framework's built-in modules that you import and use.
You don't modify these files.
• core/router/ - HTTP routing engine
• core/database/ - DB connection management
• core/logger/ - Logging infrastructure
• core/emitter/ - Event system
• core/storage/ - File storage
• core/middleware/ - HTTP middleware
### Your Application Code
Your business logic
#### app/ directory
Contains your application's modules, models, and business logic.
This is where you write your code.
• app/models/ - Your database models
• app/user/ - User management module
• app/post/ - Blog post module
• app/auth/ - Authentication module
• app/init.go - Module registration
## File Naming Conventions
### Go Files
snake_case
File names
user_service.go, blog_post.go
PascalCase
Struct and function names
UserService, CreatePost()
camelCase
Private fields and methods
userID, validateEmail()
### Module Structure
```
app/user/
├── controller.go    # HTTP handlers and routing
├── service.go      # Business logic and validation
├── module.go       # Module registration and dependencies
└── README.md       # Module-specific documentation
```
## Generated Files
#### When running base g User name:string email:string
```
✓ Created app/models/user.go
✓ Created app/user/controller.go
✓ Created app/user/service.go
✓ Created app/user/module.go
✓ Updated app/init.go
```
## Configuration Files
### .env
Environment-specific configuration
• Database connection settings
• JWT secrets and API keys
• Storage provider configuration
• Email service settings
### go.mod
Go module dependencies
• Base Framework dependencies
• GORM and database drivers
• Third-party packages
## Best Practices
#### Keep modules focused
Each module should handle a single domain concept (User, Post, Order, etc.)
#### Use the models directory
Always place database models in app/models/ to avoid circular imports
#### Follow naming conventions
Consistent naming makes the codebase predictable and easier to navigate