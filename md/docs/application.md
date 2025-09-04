---
title: Application
description: Application documentation for Base Framework.
---

# App Directory
Understanding Base Framework's application layer, module initialization, and business logic organization.
## Overview
The `app/` directory is the heart of your Base Framework application. It's where your business logic lives, modules are organized, and application-specific functionality is implemented. Unlike the `core/` directory which contains framework infrastructure, the `app/` directory is entirely yours to customize.
## Directory Structure
```
app/
├── init.go              # Module registration hub
├── models/              # Shared data models
│   ├── user.go
│   ├── post.go
│   └── category.go
├── post/
|   ├── controller.go  # HTTP handlers
|   ├── service.go     # Business logic
|   └── module.go      # Module definition
├── user/
|   ├── controller.go
|   ├── service.go
|   └── module.go
└── category/
├── controller.go
├── service.go
└── module.go
```
### init.go
Central module registry that tells the framework which modules to load and initialize. This is the entry point for all your business logic.
### models/
Shared data models (GORM structs) that can be used across multiple modules. Centralized to prevent **circular dependencies**.
## Module Initialization Flow
Understanding how modules get initialized is crucial for building maintainable Base applications. Here's what happens under the hood:
1
### Application Startup
When your Base application starts, `main.go` creates the core application infrastructure:
- • Database connection (GORM)
- • Router with middleware
- • Logger, Emitter, Storage systems
- • Email sender and configuration
- • Websocket server
- • Translation system
- • Validator system
- • Helper utilities
- • Error handling
- • Configuration
2
### Core Modules First
The framework initializes core modules (authentication, permissions, etc.) before your app modules. This ensures all foundation services are available.
- • Authentication module
- • Authorization module
- • Media module
- • OAuth module
- • Profile module
3
### App Module Discovery
The framework calls `app/init.go:GetAppModules()` to discover your business modules:
```go
func (am *AppModules) GetAppModules(deps module.Dependencies) map[string]module.Module {
modules := make(map[string]module.Module)
// Add your modules here:
modules["post"] = post.NewModule(deps)
modules["user"] = user.NewModule(deps)
return modules
}
```
4
### Module Lifecycle
For each module, the framework executes the complete lifecycle:
1. $1
2. $1
3. $1
4. $1
## Dependency Injection
Base Framework uses dependency injection to provide your modules with access to framework services. Every module receives a `Dependencies` struct containing:
### Core Dependencies
- `DB` - GORM database instance
- `Router` - HTTP router for defining endpoints
- `Logger` - Structured logging service
- `Config` - Application configuration
### Extended Services
- `Emitter` - Event system for module communication
- `Storage` - File upload and management
- `EmailSender` - Email service integration
- `Translation` - Internationalization service
- `Validator` - Validation service
You can check each service's documentation for more details.
### Usage Example
```go
// In your module's service
func (s *PostService) CreatePost(data CreatePostRequest) (*models.Post, error) {
// Use injected database
post := &models.Post{
Title:   data.Title,
Content: data.Content,
UserID:  data.UserID,
}
if err := s.deps.DB.Create(post).Error; err != nil {
s.deps.Logger.Error("Failed to create post", logger.String("error", err.Error()))
return nil, err
}
// Emit event for other modules to react
s.deps.Emitter.Emit("post.created", post)
return post, nil
}
```
## MCS Architecture Pattern
Base Framework implements the Model-Controller-Service (MCS) architecture pattern to ensure clean separation of concerns, maintainability, and testability. Each layer has distinct responsibilities and clear boundaries.
### Model
**Data Layer:** Defines data structures, validation rules, and data relationships.
- • GORM struct definitions
- • Request/Response types
- • Data validation tags
- • Database relationships
- • Serialization methods
Located in: `app/models/`
### Controller
**HTTP Layer:** Handles web requests, responses, and HTTP-specific concerns.
- • Route definitions
- • Request parsing
- • Response formatting
- • HTTP status codes
- • Input validation
Located in: `app/[module]/controller.go`
### Service
**Business Layer:** Contains business logic, data access, and external integrations.
- • Business rules & validation
- • Database operations
- • External API calls
- • Event emission
- • Inter-service communication
Located in: `app/[module]/service.go`
### Data Flow & Responsibilities
HTTP Request
Controller
Service
Model/Database
## Implementation Examples
Here's how the MCS pattern is implemented in practice with real code examples from generated modules.
### Controller Layer (HTTP)
Handles HTTP requests, validation, response formatting, and route definitions. Controllers are thin and delegate business logic to services.
```go
func (c *PostController) CreatePost(ctx *router.Context) error {
var req CreatePostRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
return ctx.JSON(400, map[string]any{"error": err.Error()})
}
post, err := c.service.CreatePost(req)
if err != nil {
return ctx.JSON(500, map[string]any{"error": err.Error()})
}
return ctx.JSON(201, post)
}
func (c *PostController) Routes(router *router.RouterGroup) {
// Import authorization middleware
// import "base/core/app/authorization"
// Main CRUD endpoints with authorization middleware
router.GET("/posts", c.List)                                               // Public read
router.POST("/posts", authorization.Can("create", "post"), c.Create)       // Requires create permission
router.GET("/posts/all", c.ListAll)                                        // Public read
router.GET("/posts/:id", c.Get)                                            // Public read
router.PUT("/posts/:id", authorization.CanAccess("update", "post", "id"), c.Update)    // Requires resource permission
router.DELETE("/posts/:id", authorization.CanAccess("delete", "post", "id"), c.Delete) // Requires resource permission
}
```
### Service Layer (Business Logic)
Contains the core business logic, data validation, and business rules. Services can call other services and emit events.
```go
func (s *PostService) CreatePost(data CreatePostRequest) (*models.Post, error) {
// Business validation only - permissions handled by middleware
if len(data.Title)
### Module Definition
Ties everything together - handles initialization, database migrations, and connects controllers to the router.
```go
type Module struct {
service    *Service
controller *Controller
deps       module.Dependencies
}
func (m *Module) Init() error {
// Module initialization logic
return nil
}
func (m *Module) Migrate() error {
// Auto-migrate models
return m.deps.DB.AutoMigrate(&models.Post{})
}
func (m *Module) Routes(router *router.RouterGroup) {
// Delegate route registration to the controller
m.Controller.Routes(router)
}
```
## Module Communication
Modules in Base Framework communicate through events and direct service injection, maintaining loose coupling while enabling powerful integrations.
### Event-Driven Communication
Asynchronous, decoupled communication using the built-in event emitter:
```go
// Emit events
emitter.Emit("user.created", user)
emitter.Emit("post.published", post)
// Listen for events
emitter.On("user.created", handleUserCreated)
emitter.On("post.published", handlePostPublished)
```
### Direct Service Injection
Synchronous communication by injecting services directly:
```go
type PostService struct {
userService *user.Service
deps        module.Dependencies
}
func (s *PostService) CreatePost(data PostData) error {
// Use injected user service
user, err := s.userService.GetUser(data.UserID)
if err != nil {
return err
}
// ... rest of logic
}
```
## CLI Integration
The Base CLI integrates seamlessly with the app directory structure, generating modules that follow the established patterns.
### Generating Modules
```bash
# Generate a complete module with CRUD operations
base g post title:string content:text published:bool
# Generate with relationships (auto-detected from _id suffix)
base g comment post_id:uint content:text user_id:uint
# Generate with file uploads
base g profile user_id:uint avatar:image bio:text
```
This automatically creates the module directory, controller, service, model, and updates `app/init.go` to register the new module.
### What Gets Generated
`app/models/post.go` - GORM model with relationships
`app/post/` - Module directory
`app/post/controller.go` - HTTP handlers
`app/post/service.go` - Business logic
`app/post/module.go` - Module definition
`app/init.go` - Updated with new module registration
## Best Practices
### ✅ Do
- • Keep controllers thin - delegate to services
- • Use events for cross-module communication
- • Place shared models in `app/models/`
- • Follow the generated module structure
- • Use dependency injection for testability
- • Emit events for important business actions
### ❌ Don't
- • Put business logic in controllers
- • Create circular dependencies between modules
- • Access database directly from controllers
- • Hardcode configuration values
- • Skip error handling and logging
- • Ignore the module lifecycle methods
## Example: Complete Module Implementation
Here's how a complete module looks when following Base Framework patterns:
### Controller (HTTP Layer)
Handles HTTP requests, validation, and delegates to service layer:
```go
// app/post/controller.go
package post
import (
"base/core/router"
"net/http"
"strconv"
)
type Controller struct {
service *Service
}
func NewController(service *Service) *Controller {
return &Controller{service: service}
}
func (c *Controller) Create(ctx *router.Context) error {
var req models.CreatePostRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
return ctx.JSON(http.StatusBadRequest, map[string]any{
"error": err.Error(),
})
}
post, err := c.service.Create(&req)
if err != nil {
return ctx.JSON(http.StatusInternalServerError, map[string]any{
"error": err.Error(),
})
}
return ctx.JSON(http.StatusCreated, post.ToResponse())
}
func (c *Controller) List(ctx *router.Context) error {
posts, err := c.service.GetAll(nil, nil, nil, nil)
if err != nil {
return ctx.JSON(http.StatusInternalServerError, map[string]any{
"error": err.Error(),
})
}
return ctx.JSON(http.StatusOK, posts)
}
func (c *Controller) Get(ctx *router.Context) error {
id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
if err != nil {
return ctx.JSON(http.StatusBadRequest, map[string]any{
"error": "Invalid ID format",
})
}
post, err := c.service.GetById(uint(id))
if err != nil {
return ctx.JSON(http.StatusNotFound, map[string]any{
"error": "Post not found",
})
}
return ctx.JSON(http.StatusOK, post.ToResponse())
}
func (c *Controller) ListAll(ctx *router.Context) error {
posts, err := c.service.GetAllForSelect()
if err != nil {
return ctx.JSON(http.StatusInternalServerError, map[string]any{
"error": err.Error(),
})
}
return ctx.JSON(http.StatusOK, posts)
}
func (c *Controller) Update(ctx *router.Context) error {
id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
if err != nil {
return ctx.JSON(http.StatusBadRequest, map[string]any{
"error": "Invalid ID format",
})
}
var req models.UpdatePostRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
return ctx.JSON(http.StatusBadRequest, map[string]any{
"error": err.Error(),
})
}
post, err := c.service.Update(uint(id), &req)
if err != nil {
return ctx.JSON(http.StatusInternalServerError, map[string]any{
"error": err.Error(),
})
}
return ctx.JSON(http.StatusOK, post.ToResponse())
}
func (c *Controller) Delete(ctx *router.Context) error {
id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
if err != nil {
return ctx.JSON(http.StatusBadRequest, map[string]any{
"error": "Invalid ID format",
})
}
if err := c.service.Delete(uint(id)); err != nil {
return ctx.JSON(http.StatusInternalServerError, map[string]any{
"error": err.Error(),
})
}
return ctx.NoContent()
}
// Routes defines all the HTTP routes for this controller
func (c *Controller) Routes(router *router.RouterGroup) {
// Main CRUD endpoints - specific routes MUST come before parameterized routes
router.GET("/posts", c.List)          // Paginated list
router.POST("/posts", c.Create)       // Create
router.GET("/posts/all", c.ListAll)   // Unpaginated list - MUST be before /:id
router.GET("/posts/:id", c.Get)       // Get by ID - MUST be after /all
router.PUT("/posts/:id", c.Update)    // Update
router.DELETE("/posts/:id", c.Delete) // Delete
}
```
### Service (Business Layer)
Contains business logic, data access, and event emission:
```go
// app/post/service.go
package post
import (
"base/app/models"
"base/core/module"
"base/core/logger"
"base/core/types"
"errors"
)
type Service struct {
deps module.Dependencies
}
func NewService(deps module.Dependencies) *Service {
return &Service{deps: deps}
}
func (s *Service) Create(data *models.CreatePostRequest) (*models.Post, error) {
// Business validation
if len(data.Title)
### Module (Registration & Wiring)
Ties everything together with dependency injection and route registration:
```go
// app/posts/module.go
package posts
import (
"base/app/models"
"base/core/module"
"base/core/router"
"gorm.io/gorm"
)
type Module struct {
module.DefaultModule
DB         *gorm.DB
Service    *Service
Controller *Controller
}
// Init creates and initializes the Post module with all dependencies
func Init(deps module.Dependencies) module.Module {
// Initialize service and controller
service := NewService(deps)
controller := NewController(service)
// Create module
mod := &Module{
DB:         deps.DB,
Service:    service,
Controller: controller,
}
return mod
}
// Routes registers the module routes
func (m *Module) Routes(router *router.RouterGroup) {
m.Controller.Routes(router)
}
func (m *Module) Init() error {
return nil
}
func (m *Module) Migrate() error {
return m.DB.AutoMigrate(&models.Post{})
}
func (m *Module) GetModels() []any {
return []any{
&models.Post{},
}
}
```
### Registration (App Init)
How the module gets registered in the application startup:
```go
// app/init.go
package app
import (
"base/app/posts"
"base/core/module"
)
// AppModules implements module.AppModuleProvider interface
type AppModules struct{}
// GetAppModules returns the list of app modules to initialize
// This is the only function that needs to be updated when adding new app modules
func (am *AppModules) GetAppModules(deps module.Dependencies) map[string]module.Module {
modules := make(map[string]module.Module)
// Posts module
modules["posts"] = posts.Init(deps)
return modules
}
// NewAppModules creates a new AppModules provider
func NewAppModules() *AppModules {
return &AppModules{}
}
```