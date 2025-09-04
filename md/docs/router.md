---
title: Router
description: Router documentation for Base Framework.
---

# Router
Zero-dependency HTTP router with radix tree routing, middleware chaining, and high-performance request handling
## Overview
🚀 High-Performance Router
Base Framework's router is a lightweight, zero-dependency HTTP router built with radix tree routing for optimal performance. It features middleware chaining, context pooling, and comprehensive request/response handling.
Radix Tree Routing
Efficient URL pattern matching with parameter extraction and wildcard support.
Middleware Chaining
Composable middleware with global and route-specific support for cross-cutting concerns.
Context Pooling
Memory-efficient context reuse with zero-allocation routing for high-throughput applications.
## Route Definition & Registration
Basic Route Registration
Basic HTTP Methods
```
r := router.New()
// Basic HTTP methods
r.GET("/users", handleGetUsers)
r.POST("/users", handleCreateUser)
r.PUT("/users/:id", handleUpdateUser)
r.DELETE("/users/:id", handleDeleteUser)
r.PATCH("/users/:id", handlePatchUser)
// Generic handler for any HTTP method
r.Handle("GET", "/custom", handleCustom)
```
Route Groups
Grouping Routes with Prefixes
```
// API v1 group
api := r.Group("/api/v1")
api.GET("/users", handleGetUsers)
api.POST("/users", handleCreateUser)
// Admin group with middleware
admin := r.Group("/admin", AuthMiddleware(), AdminMiddleware())
admin.GET("/users", handleAdminGetUsers)
admin.DELETE("/users/:id", handleAdminDeleteUser)
// Nested groups
v2 := api.Group("/v2")
v2.GET("/posts", handleGetPostsV2)
// Group with multiple middleware
protected := r.Group("/protected",
middleware.Logger(loggerConfig),
middleware.Auth(authConfig),
middleware.RateLimit(rateLimitConfig),
)
protected.GET("/profile", handleProfile)
```
## Context API & Request Handling
Context Interface
Basic Handler Function
```
func handleUser(c *router.Context) error {
// Access request information
method := c.Request.Method
path := c.Request.URL.Path
userAgent := c.Request.UserAgent()
clientIP := c.ClientIP()
// Get headers
authHeader := c.Header("Authorization")
contentType := c.ContentType()
// Store data in context for other middleware/handlers
c.Set("user_id", 123)
c.Set("request_start", time.Now())
// Retrieve stored data
userID, exists := c.Get("user_id")
if exists {
c.Logger.Info("Processing for user", "id", userID)
}
return c.JSON(200, map[string]any{
"message": "Success",
"ip": clientIP,
"method": method,
})
}
```
## Route Parameters & Query Strings
URL Parameters
Parameter Extraction
```
// Route with parameters
r.GET("/users/:id/posts/:postId", func(c *router.Context) error {
// Named parameters
userID := c.Param("id")
postID := c.Param("postId")
return c.JSON(200, map[string]string{
"user_id": userID,
"post_id": postID,
})
})
// Wildcard routes
r.GET("/files/*filepath", func(c *router.Context) error {
filepath := c.Param("filepath")
// filepath captures everything after /files/
return c.File("./static/" + filepath)
})
// Multiple parameters with validation
r.GET("/api/v1/users/:id/orders/:orderNumber", func(c *router.Context) error {
userID := c.Param("id")
orderNumber := c.Param("orderNumber")
// Convert to appropriate types
id, err := strconv.Atoi(userID)
if err != nil {
return c.JSON(400, map[string]string{"error": "Invalid user ID"})
}
return handleOrder(c, id, orderNumber)
})
```
Query Parameters
Query String Handling
```
func handleSearch(c *router.Context) error {
// Single query parameter
query := c.Query("q")
page := c.DefaultQuery("page", "1")
// Check if parameter exists
category, exists := c.GetQuery("category")
if !exists {
category = "all"
}
// Multiple values for same parameter
tags, hasTag := c.GetQueryArray("tags")
if hasTag {
// Process array of tags: ?tags=go&tags=web&tags=api
c.Logger.Info("Tags", "values", tags)
}
// Convert query parameters
limitStr := c.DefaultQuery("limit", "10")
limit, err := strconv.Atoi(limitStr)
if err != nil {
return c.JSON(400, map[string]string{
"error": "Invalid limit parameter",
})
}
return c.JSON(200, map[string]any{
"query": query,
"page": page,
"category": category,
"tags": tags,
"limit": limit,
})
}
```
## Request Binding & Validation
JSON & Form Binding
Auto Content-Type Detection
```
type CreateUserRequest struct {
Name     string `json:"name" form:"name" validate:"required,min=2"`
Email    string `json:"email" form:"email" validate:"required,email"`
Age      int    `json:"age" form:"age" validate:"min=18,max=100"`
IsActive bool   `json:"is_active" form:"is_active"`
}
func createUser(c *router.Context) error {
var req CreateUserRequest
// Auto-detects content type and binds accordingly
if err := c.Bind(&req); err != nil {
return c.JSON(400, map[string]string{
"error": "Invalid request format",
"details": err.Error(),
})
}
// Manual binding for specific content types
var jsonReq CreateUserRequest
if err := c.BindJSON(&jsonReq); err != nil {
return c.JSON(400, map[string]string{"error": "Invalid JSON"})
}
// Form binding
var formReq CreateUserRequest
if err := c.BindForm(&formReq); err != nil {
return c.JSON(400, map[string]string{"error": "Invalid form data"})
}
// Query parameter binding
var queryReq CreateUserRequest
if err := c.BindQuery(&queryReq); err != nil {
return c.JSON(400, map[string]string{"error": "Invalid query params"})
}
return c.JSON(201, req)
}
```
File Upload Handling
Multipart File Handling
```
func uploadFile(c *router.Context) error {
// Single file upload
file, err := c.FormFile("upload")
if err != nil {
return c.JSON(400, map[string]string{
"error": "No file provided",
})
}
// Access file metadata
filename := file.Filename
size := file.Size
contentType := file.Header.Get("Content-Type")
// Process multipart form
form, err := c.MultipartForm()
if err != nil {
return c.JSON(400, map[string]string{
"error": "Invalid multipart form",
})
}
// Handle multiple files
files := form.File["files"]
for _, fileHeader := range files {
// Process each file
file, err := fileHeader.Open()
if err != nil {
continue
}
defer file.Close()
// Save file logic here
c.Logger.Info("Processing file", "name", fileHeader.Filename)
}
// Get form values from multipart form
title := c.FormValue("title")
description := c.FormValue("description")
return c.JSON(200, map[string]any{
"message": "Files uploaded successfully",
"filename": filename,
"size": size,
"content_type": contentType,
"title": title,
})
}
```
## Response Methods & Content Types
JSON & Data Responses
Response Formats
```
func handleResponses(c *router.Context) error {
// JSON response
return c.JSON(200, map[string]any{
"message": "Success",
"data": []string{"item1", "item2"},
"count": 2,
})
// String response
return c.String(200, "Plain text response")
// HTML response
html := `&lt;h1&gt;Hello World&lt;/h1&gt;&lt;p&gt;This is HTML content&lt;/p&gt;`
return c.HTML(200, html)
// Raw data response
data := []byte("binary data here")
return c.Data(200, "application/octet-stream", data)
// File response
c.File("/path/to/file.pdf")
// No content response
return c.NoContent()
// Custom headers with response
c.SetHeader("X-Custom-Header", "custom-value")
c.SetHeader("Cache-Control", "max-age=3600")
return c.JSON(200, map[string]string{"status": "ok"})
}
// Error response helper
func handleError(c *router.Context) error {
err := errors.New("something went wrong")
return c.Error(500, err) // Automatically formats as JSON error
}
// Redirect responses
func handleRedirect(c *router.Context) error {
return c.Redirect(302, "https://example.com")
}
```
## Middleware Usage & Custom Middleware
Built-in Middleware
Common Middleware Setup
```
import (
"base/core/router"
"base/core/router/middleware"
"base/core/logger"
)
func setupRouter() *router.Router {
r := router.New()
log := logger.New()
// Global middleware (applied to all routes)
r.Use(middleware.Recovery(log))
r.Use(middleware.RequestID())
r.Use(middleware.Logger(middleware.DefaultLoggerConfig(log)))
r.Use(middleware.RateLimit(middleware.DefaultRateLimitConfig()))
// Authentication middleware
authConfig := middleware.DefaultAuthConfig()
authConfig.TokenValidator = validateJWTToken
authConfig.SkipPaths = []string{"/login", "/register", "/health"}
// Public routes
r.GET("/health", handleHealth)
r.POST("/login", handleLogin)
r.POST("/register", handleRegister)
// Protected routes group
api := r.Group("/api/v1", middleware.Auth(authConfig))
api.GET("/profile", handleProfile)
api.PUT("/profile", handleUpdateProfile)
// Admin routes with additional middleware
admin := r.Group("/admin",
middleware.Auth(authConfig),
middleware.RequireAuth("user"),
adminMiddleware,
)
admin.GET("/users", handleGetUsers)
admin.DELETE("/users/:id", handleDeleteUser)
return r
}
```
Custom Middleware
Creating Custom Middleware
```
// Custom CORS middleware
func CORSMiddleware() router.MiddlewareFunc {
return func(next router.HandlerFunc) router.HandlerFunc {
return func(c *router.Context) error {
// Set CORS headers
c.SetHeader("Access-Control-Allow-Origin", "*")
c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
c.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
// Handle preflight requests
if c.Request.Method == "OPTIONS" {
return c.NoContent()
}
return next(c)
}
}
}
// Timing middleware
func TimingMiddleware() router.MiddlewareFunc {
return func(next router.HandlerFunc) router.HandlerFunc {
return func(c *router.Context) error {
start := time.Now()
// Process request
err := next(c)
// Add timing header
duration := time.Since(start)
c.SetHeader("X-Response-Time", duration.String())
return err
}
}
}
// Role-based access control middleware
func RequireRole(role string) router.MiddlewareFunc {
return func(next router.HandlerFunc) router.HandlerFunc {
return func(c *router.Context) error {
user, exists := c.Get("user")
if !exists {
return c.JSON(401, map[string]string{
"error": "Authentication required",
})
}
userRole := getUserRole(user)
if userRole != role && userRole != "admin" {
return c.JSON(403, map[string]string{
"error": "Insufficient permissions",
})
}
return next(c)
}
}
}
// API versioning middleware
func APIVersionMiddleware() router.MiddlewareFunc {
return func(next router.HandlerFunc) router.HandlerFunc {
return func(c *router.Context) error {
version := c.Header("API-Version")
if version == "" {
version = "v1" // default
}
c.Set("api_version", version)
c.SetHeader("API-Version", version)
return next(c)
}
}
}
```
## Error Handling & HTTP Status Codes
Error Responses & Recovery
Error Handling Patterns
```
// Custom error types
type APIError struct {
Code    int    `json:"code"`
Message string `json:"message"`
Details any    `json:"details,omitempty"`
}
func (e APIError) Error() string {
return e.Message
}
// Error handler with detailed responses
func handleUserCreate(c *router.Context) error {
var req CreateUserRequest
if err := c.Bind(&req); err != nil {
return c.JSON(400, APIError{
Code:    40001,
Message: "Invalid request format",
Details: err.Error(),
})
}
// Validation
if req.Email == "" {
return c.JSON(400, APIError{
Code:    40002,
Message: "Email is required",
})
}
user, err := userService.Create(req)
if err != nil {
if errors.Is(err, ErrUserExists) {
return c.JSON(409, APIError{
Code:    40901,
Message: "User already exists",
})
}
// Log internal errors but don't expose details
c.Logger.Error("Failed to create user", "error", err)
return c.JSON(500, APIError{
Code:    50001,
Message: "Internal server error",
})
}
return c.JSON(201, user)
}
// Global error handler
func ErrorHandler(c *router.Context, err error) {
var apiErr APIError
switch e := err.(type) {
case APIError:
apiErr = e
case *validation.ValidationError:
apiErr = APIError{
Code:    40003,
Message: "Validation failed",
Details: e.Fields,
}
default:
apiErr = APIError{
Code:    50000,
Message: "Internal server error",
}
}
c.JSON(apiErr.Code/100, apiErr)
}
// Custom 404 handler
func notFoundHandler(c *router.Context) error {
return c.JSON(404, APIError{
Code:    40401,
Message: "Endpoint not found",
Details: map[string]string{
"path":   c.Request.URL.Path,
"method": c.Request.Method,
},
})
}
// Setup router with error handling
func setupErrorHandling() {
r := router.New()
// Set custom 404 handler
r.NotFound(notFoundHandler)
// Recovery middleware
r.Use(middleware.Recovery(logger))
}
```
## Performance Considerations & Best Practices
Performance Features
Zero-allocation radix tree routing
Context pooling for memory efficiency
Fast parameter extraction
Efficient middleware chaining
Thread-safe operations
Benchmarks
Route lookup:
~50ns
Context creation:
~0 allocs
Parameter extraction:
~30ns
Memory per request:
~0.5KB
Throughput:
100K+ RPS
✅ Performance Best Practices
- Use route groups to minimize middleware overhead
- Place more specific routes before wildcard routes
- Minimize middleware chain length for critical paths
- Use context pooling (automatic in Base router)
- Implement proper error handling to avoid panics
- Use built-in response methods for better performance
❌ Performance Anti-patterns
- Don't create new context instances manually
- Avoid heavy computation in middleware for all routes
- Don't ignore context cancellation signals
- Avoid deep middleware nesting without caching
- Don't use blocking operations in hot paths
- Avoid large response payloads without streaming
## Static File Serving & Advanced Features
Static File Configuration
Static File Setup
```
// Basic static file serving
r.Static("/static", "./public")
r.Static("/uploads", "./storage/uploads")
// Group-specific static files
api := r.Group("/api/v1")
api.Static("/assets", "./api-assets")
// Custom file handler with middleware
r.GET("/files/*filepath", middleware.Auth(authConfig), func(c *router.Context) error {
filepath := c.Param("filepath")
// Security check
if strings.Contains(filepath, "..") {
return c.JSON(400, map[string]string{"error": "Invalid path"})
}
fullPath := path.Join("./secure-files", filepath)
// Check file exists and user has access
if !hasFileAccess(c.MustGet("user"), fullPath) {
return c.JSON(403, map[string]string{"error": "Access denied"})
}
c.File(fullPath)
return nil
})
// Server setup with graceful shutdown
func main() {
r := router.New()
setupRoutes(r)
// Start server
if err := r.Run(":8080"); err != nil {
log.Fatal("Server failed to start:", err)
}
}
```