---
title: Logger
description: Logger documentation for Base Framework.
---

# Logger
Structured logging system with typed fields, context integration, and beautiful console output
## Overview
📊 Structured Logging
Base Framework's logger provides structured logging with typed fields, context integration, and multiple output formats. Built on Zap for high performance, it includes beautiful console output with emojis and colors plus JSON file logging for production.
Typed Fields
Use strongly-typed fields like String(), Int(), Duration() for structured, searchable logs.
Beautiful Console
Colorized console output with emojis and clear formatting for development.
JSON Files
Production-ready JSON logs written to files for analysis and monitoring.
## Configuration
Environment Variables
.env Configuration
```
# Log Level (debug, info, warn, error, fatal)
LOG_LEVEL=info
# Log Directory Path
LOG_PATH=logs
# Environment (development/production)
ENVIRONMENT=development
```
Manual Configuration
logger initialization
```
import "base/core/logger"
// Create logger with custom config
config := logger.Config{
Environment: "development", // or "production"
LogPath:     "logs",
Level:       "debug",
}
log, err := logger.NewLogger(config)
if err != nil {
panic(err)
}
// Use the logger
log.Info("Application started",
logger.String("version", "1.0.0"),
logger.String("environment", config.Environment))
```
## Log Levels
Base Framework supports five log levels, each with specific use cases and visual indicators.
🔍
DEBUG
Detailed information for diagnosing problems. Only shown in development.
`log.Debug("User query", logger.String("sql", query))`
ℹ️
INFO
General information about application flow and important events.
`log.Info("User created", logger.Uint("user_id", 123))`
⚠️
WARN
Something unexpected happened but the application can continue.
`log.Warn("Email failed", logger.String("recipient", email))`
❌
ERROR
Error conditions that should be investigated and fixed.
`log.Error("Database error", logger.String("error", err.Error()))`
💀
FATAL
Critical errors that cause the application to exit immediately.
`log.Fatal("Cannot connect to database", logger.String("dsn", dsn))`
## Structured Logging with Typed Fields
Available Field Types
Typed Field Examples
```
// String fields
logger.String("username", "john_doe")
logger.String("error", err.Error())
// Numeric fields
logger.Int("status_code", 200)
logger.Int64("user_id", 123456)
logger.Uint("count", 42)
logger.Float64("price", 29.99)
// Boolean fields
logger.Bool("is_admin", true)
logger.Bool("email_verified", false)
// Time and duration
logger.Duration("request_time", time.Since(start))
// Complex data
logger.Any("metadata", map[string]interface{}{
"ip": "192.168.1.1",
"user_agent": "Mozilla/5.0...",
})
```
## Context Integration & Request Tracing
Creating Child Loggers
Context-aware logging
```
// Create a child logger with context
requestLogger := logger.With(
logger.String("request_id", requestID),
logger.String("user_id", userID),
logger.String("endpoint", "/api/users"),
)
// All subsequent logs include context
requestLogger.Info("Processing request")
requestLogger.Debug("Validating input",
logger.Any("input", requestData))
requestLogger.Info("Request completed",
logger.Duration("duration", time.Since(start)))
```
💡 Context Best Practice
Create child loggers with request-specific context at the beginning of request handlers. This ensures all logs within that request scope are automatically tagged with relevant information like request ID, user ID, and endpoint.
## Usage Examples
Controller Logging
controller.go
```
type UserController struct {
service *UserService
logger  logger.Logger
}
func (c *UserController) CreateUser(ctx *router.Context) error {
// Create request-scoped logger
requestLogger := c.logger.With(
logger.String("endpoint", "POST /api/users"),
logger.String("request_id", ctx.Get("request_id").(string)),
)
requestLogger.Info("Creating new user")
var req CreateUserRequest
if err := ctx.ShouldBind(&req); err != nil {
requestLogger.Warn("Invalid request data",
logger.String("error", err.Error()),
logger.Any("request", req))
return ctx.JSON(400, types.ErrorResponse{Error: err.Error()})
}
user, err := c.service.Create(&req)
if err != nil {
requestLogger.Error("Failed to create user",
logger.String("error", err.Error()),
logger.String("email", req.Email))
return ctx.JSON(500, types.ErrorResponse{Error: "Internal error"})
}
requestLogger.Info("User created successfully",
logger.Uint("user_id", user.ID),
logger.String("email", user.Email))
return ctx.JSON(201, user)
}
```
Service Layer Logging
service.go
```
type EmailService struct {
db     *gorm.DB
logger logger.Logger
}
func (s *EmailService) SendWelcomeEmail(userID uint, email string) error {
start := time.Now()
// Create service-scoped logger
serviceLogger := s.logger.With(
logger.String("service", "EmailService"),
logger.String("operation", "SendWelcomeEmail"),
logger.Uint("user_id", userID),
)
serviceLogger.Info("Starting email send",
logger.String("recipient", email),
logger.String("template", "welcome"))
// Fetch user details
var user User
if err := s.db.First(&user, userID).Error; err != nil {
serviceLogger.Error("Failed to fetch user",
logger.String("error", err.Error()))
return err
}
serviceLogger.Debug("User fetched",
logger.String("name", user.Name),
logger.Bool("email_verified", user.EmailVerified))
// Send email (simulated)
if err := s.sendEmail(email, "welcome", user); err != nil {
serviceLogger.Error("Email delivery failed",
logger.String("error", err.Error()),
logger.Duration("attempt_duration", time.Since(start)))
return err
}
serviceLogger.Info("Welcome email sent successfully",
logger.Duration("total_duration", time.Since(start)))
return nil
}
```
## Middleware Integration
Request Logging Middleware
middleware setup
```
import "base/core/router/middleware"
// Configure logger middleware
loggerConfig := &middleware.LoggerConfig{
Logger:         log,
LogLevel:       "info",
SkipPaths:      []string{"/health", "/metrics"},
IncludeHeaders: true,  // Include request headers in logs
IncludeBody:    false, // Don't log request bodies by default
}
// Add to router
router.Use(middleware.Logger(loggerConfig))
// Add recovery middleware with logging
router.Use(middleware.Recovery(log))
// Add request ID middleware for tracing
router.Use(middleware.RequestID())
```
Custom Access Log Format
custom access logs
```
// Create custom access log with specific format
format := ":method :path :status :latency :ip - :user_agent"
router.Use(middleware.AccessLog(format, log))
// Example output:
// INFO   GET /api/users 200 45.2ms 127.0.0.1 - Mozilla/5.0...
// Available tokens:
// :method     - HTTP method (GET, POST, etc.)
// :path       - Request path
// :status     - Response status code
// :latency    - Request duration
// :ip         - Client IP address
// :user_agent - User agent string
```
## Output Examples
Console Output (Development)
Development Console
```
2024-01-15 14:30:25  ℹ️  INFO   main.go:118  🚀 Starting Base Framework  version=1.0.0 environment=development
2024-01-15 14:30:25  ℹ️  INFO   database.go:45  ✅ Database initialized
2024-01-15 14:30:26  ℹ️  INFO   middleware.go:102  Request  method=GET path=/api/users status=200 latency=24.5ms
2024-01-15 14:30:27  ⚠️  WARN   user.go:67  Email delivery failed  user_id=123 error=smtp timeout
2024-01-15 14:30:28  ❌ ERROR  database.go:89  Connection lost  error=connection refused retry_count=3
```
JSON File Output (Production)
logs/app.log
```
{"level":"info","ts":"2024-01-15T14:30:25.123Z","caller":"main.go:118","msg":"Starting Base Framework","version":"1.0.0","environment":"production"}
{"level":"info","ts":"2024-01-15T14:30:26.456Z","caller":"middleware.go:102","msg":"Request","method":"POST","path":"/api/users","status":201,"latency":"45.2ms","ip":"192.168.1.100","user_agent":"PostmanRuntime/7.28.4"}
{"level":"error","ts":"2024-01-15T14:30:27.789Z","caller":"user.go:67","msg":"Failed to create user","user_id":123,"email":"user@example.com","error":"validation failed","request_id":"req_abc123"}
```
## Performance Tips & Best Practices
✅ Performance Best Practices
- Use typed fields (String, Int, etc.) instead of Any() when possible
- Create child loggers with With() for request context instead of repeating fields
- Set appropriate log levels - use DEBUG sparingly in production
- Avoid logging large objects with Any() - serialize critical fields only
- Use Duration() fields to track performance metrics
- Configure log rotation to prevent disk space issues
💡 Logging Guidelines
- Log at entry/exit points of important functions with context
- Always log errors with relevant context (user_id, request_id, etc.)
- Use consistent field names across your application (user_id, not userId)
- Include correlation IDs to trace requests across services
- Log business events for analytics and monitoring
- Don't log sensitive information like passwords or tokens
⚠️ Common Pitfalls
- Don't create new loggers in hot paths - reuse instances
- Avoid string concatenation in log messages - use fields instead
- Don't log at DEBUG level in production unless necessary
- Be careful with Any() - it can impact performance with large objects
- Don't ignore errors from logging operations in critical paths
- Avoid excessive logging in loops - consider sampling or batching