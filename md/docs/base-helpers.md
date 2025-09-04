---
title: Base Helpers
description: Base Helpers documentation for Base Framework.
---

# Base Helpers
Base Framework provides a comprehensive set of helper functions and utilities to streamline development and reduce boilerplate code.
## Overview
Base Framework includes several helper packages designed to handle common development tasks like JWT management, context operations, slug generation, dynamic model retrieval, standardized responses, error handling, and custom data types. These helpers promote code consistency and reduce repetitive implementation across your modules.
## Helper Categories
### Core Helpers
- • `helper/` - General utilities
- • `base/` - Base classes for modules
- • `types/` - Custom types and response formats
### Authentication
- • JWT token generation
- • Token validation
- • Context-based user retrieval
### HTTP Responses
- • Standardized success responses
- • Error handling utilities
- • Pagination helpers
### Dynamic Operations
- • Dynamic model retrieval
- • Slug generation and validation
- • Context type conversion
## Core Helper Functions
### JWT Authentication Helpers
Helper functions for JWT token generation and validation:
```go
// Generate a JWT token for a user
token, err := helper.GenerateJWT(userID)
if err != nil {
return err
}
// Generate JWT with extended claims
token, err := types.GenerateJWT(userID, map[string]any{
"role": "admin",
"permissions": []string{"read", "write"},
})
// Validate JWT token and extract user ID
userID, err := types.ValidateJWT(tokenString)
if err != nil {
return fmt.Errorf("invalid token: %w", err)
}
// Legacy validation (from helper package)
_, userID, err := helper.ValidateJWT(tokenString)
if err != nil {
return fmt.Errorf("token validation failed: %w", err)
}
```
### Context Utilities
Type-safe context value retrieval with automatic conversion:
```go
// Get string values from context
userRole := helper.GetContextString(ctx, "user_role")
authToken := helper.GetContextString(ctx, "auth_token")
// Get numeric values with automatic conversion
userID := helper.GetContextUint(ctx, "user_id")        // Returns uint
priority := helper.GetContextInt(ctx, "priority")      // Returns int
weight := helper.GetContextFloat(ctx, "weight")        // Returns float64
// Get boolean values
isAdmin := helper.GetContextBool(ctx, "is_admin")
isActive := helper.GetContextBool(ctx, "active")
// Context keys are automatically prefixed with "base_"
// So "user_id" becomes "base_user_id" internally
```
### Dynamic Model Operations
Dynamically retrieve objects based on field names and values:
```go
// First, register your models for dynamic retrieval
helper.RegisterModel("categories", func() interface{} {
return &models.Category{}
})
helper.RegisterModel("users", func() interface{} {
return &models.User{}
})
// Dynamically get objects by field relationships
// Automatically extracts table name from field name
category, err := helper.GetObject(db, "category_id", 5)
if err != nil {
return fmt.Errorf("failed to get category: %w", err)
}
user, err := helper.GetObject(db, "user_id", 123)
if err != nil {
return fmt.Errorf("failed to get user: %w", err)
}
// Type-safe version with generics
category, err := helper.GetObjectAs[models.Category](db, "category_id", 5)
if err != nil {
return fmt.Errorf("failed to get category: %w", err)
}
// The function automatically:
// 1. Extracts "category" from "category_id"
// 2. Pluralizes to "categories"
// 3. Looks up the registered model constructor
// 4. Queries the database for id = 5
```
### Slug Generation
Generate URL-friendly slugs with uniqueness guarantees:
```go
slugHelper := helper.NewSlugHelper()
// Generate slug from title
slug := slugHelper.Normalize("My Blog Post", "", "en")
// Result: "my-blog-post"
// Use custom slug if provided
slug := slugHelper.Normalize("My Blog Post", "custom-slug", "en")
// Result: "custom-slug"
// Generate unique slug with conflict resolution
existsFunc := func(slug string) (bool, error) {
var count int64
err := db.Model(&models.Post{}).Where("slug = ?", slug).Count(&count).Error
return count > 0, err
}
uniqueSlug, err := slugHelper.GenerateUniqueSlug("my-blog-post", existsFunc)
if err != nil {
return err
}
// Result: "my-blog-post" or "my-blog-post-2" if conflict exists
```
## Base Classes
Base Framework provides base classes that you can embed in your modules to get common functionality automatically.
### Base Module
Provides common module functionality with dependency injection:
```go
import "base/core/base"
type MyModule struct {
*base.Module
// Your module-specific fields
}
func NewMyModule(deps base.ModuleDependencies) *MyModule {
return &MyModule{
Module: base.NewModule("my-module", deps),
}
}
func (m *MyModule) SomeMethod() {
// Logging with module context
m.LogInfo("Operation started", logger.String("operation", "example"))
m.LogError("Operation failed", logger.String("error", "example"))
// Emit events with module context
m.EmitEvent("my-module.operation.completed", data)
// Auto-migrate models
err := m.AutoMigrate(&models.MyModel{})
// Register routes with logging
m.RegisterRoutes(router, func(r *router.RouterGroup) {
r.GET("/api/my-endpoint", handler)
})
// Validate configuration
err := m.ValidateConfig(func(cfg *config.Config) error {
if cfg.MyRequiredField == "" {
return errors.New("my_required_field is required")
}
return nil
})
}
```
### Base Service
Common service operations like logging, pagination, and database transactions:
```go
type MyService struct {
*base.Service
}
func NewMyService(db *gorm.DB, logger logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) *MyService {
return &MyService{
Service: base.NewService(db, logger, emitter, storage),
}
}
func (s *MyService) CreatePost(data *CreatePostRequest) (*models.Post, error) {
// Validate ID
if err := s.ValidateID(data.UserID); err != nil {
return nil, err
}
// Use transaction wrapper
var post *models.Post
err := s.WithTransaction(func(tx *gorm.DB) error {
post = &models.Post{Title: data.Title, Content: data.Content}
return tx.Create(post).Error
})
if err != nil {
s.LogError("create_post", err, logger.String("title", data.Title))
return nil, err
}
// Emit event
s.EmitEvent("post.created", post)
s.LogInfo("create_post", "Post created successfully", logger.Uint("post_id", post.ID))
return post, nil
}
func (s *MyService) GetPosts(page, limit int) (*types.PaginatedResponse, error) {
var posts []models.Post
var total int64
// Count total
s.DB.Model(&models.Post{}).Count(&total)
// Get paginated data
offset := (page - 1) * limit
s.DB.Offset(offset).Limit(limit).Find(&posts)
// Create paginated response
return s.CreatePaginatedResponse(posts, total, page, limit), nil
}
```
### Base Controller
Standardized HTTP response methods and parameter extraction:
```go
type MyController struct {
*base.Controller
service *MyService
}
func NewMyController(logger logger.Logger, storage *storage.ActiveStorage, service *MyService) *MyController {
return &MyController{
Controller: base.NewController(logger, storage),
service:    service,
}
}
func (c *MyController) CreatePost(ctx *router.Context) error {
var req CreatePostRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
c.RespondValidationError(ctx, err)
return nil
}
post, err := c.service.CreatePost(&req)
if err != nil {
c.RespondInternalError(ctx, "Failed to create post")
return nil
}
// Standardized success response
c.RespondCreated(ctx, post, "Post created successfully")
return nil
}
func (c *MyController) GetPost(ctx *router.Context) error {
// Extract ID parameter with validation
id, err := c.GetIDParam(ctx)
if err != nil {
c.RespondError(ctx, http.StatusBadRequest, "Invalid ID format")
return nil
}
post, err := c.service.GetPostByID(id)
if err != nil {
c.RespondNotFound(ctx, "Post")
return nil
}
c.RespondSuccess(ctx, post)
return nil
}
func (c *MyController) ListPosts(ctx *router.Context) error {
// Extract pagination parameters
page, limit := c.GetPaginationParams(ctx)
posts, err := c.service.GetPosts(page, limit)
if err != nil {
c.RespondInternalError(ctx, "Failed to retrieve posts")
return nil
}
c.RespondSuccess(ctx, posts)
return nil
}
```
## Custom Types
### DateTime Type
Enhanced datetime handling with flexible parsing and database compatibility:
```go
import "base/core/types"
type MyModel struct {
ID        uint            `json:"id"`
Name      string          `json:"name"`
CreatedAt types.DateTime  `json:"created_at"`
UpdatedAt types.DateTime  `json:"updated_at"`
}
// Usage examples
now := types.Now()                           // Current time
dt := types.DateTime{Time: time.Now()}       // From time.Time
// JSON parsing supports multiple formats:
// - RFC3339: "2006-01-02T15:04:05Z07:00"
// - MySQL: "2006-01-02 15:04:05"
// - Date only: "2006-01-02"
// Methods available on DateTime
if !dt.IsZero() {
formatted := dt.Format("2006-01-02")     // Custom formatting
future := dt.Add(24 * time.Hour)         // Add duration
duration := dt.Sub(other)                // Duration between times
isAfter := dt.After(other)               // Comparison
isBefore := dt.Before(other)             // Comparison
isEqual := dt.Equal(other)               // Equality
}
```
### Standardized Response Types
Consistent response structures for success, error, and paginated responses:
```go
// Success Response
successResponse := types.SuccessResponse{
Success: true,
Message: "Operation completed",
Data:    myData,
}
// Error Response
errorResponse := types.ErrorResponse{
Success: false,
Error:   "Something went wrong",
Details: errorDetails, // Optional additional info
}
// Paginated Response
paginatedResponse := types.PaginatedResponse{
Data: posts,
Pagination: types.Pagination{
Total:      100,
Page:       1,
PageSize:   10,
TotalPages: 10,
},
}
// JSON output examples:
{
"success": true,
"message": "Posts retrieved successfully",
"data": [...]
}
{
"success": false,
"error": "Validation failed",
"details": "Title is required"
}
{
"data": [...],
"pagination": {
"total": 100,
"page": 1,
"page_size": 10,
"total_pages": 10
}
}
```
## Error Handling Types
### Validation Error Types
Structured validation error handling with field-specific error messages:
```go
import "base/core/types"
// Single validation error
validationErr := types.ValidationError{
Field:   "email",
Message: "Email address is required",
}
// Multiple validation errors
validationResponse := types.ValidationErrorResponse{
Errors: []types.ValidationError{
{Field: "title", Message: "Title must be at least 3 characters"},
{Field: "content", Message: "Content cannot be empty"},
{Field: "category_id", Message: "Valid category is required"},
},
}
// Usage in controllers
func (c *Controller) ValidatePost(ctx *router.Context) error {
var errors []types.ValidationError
if req.Title == "" {
errors = append(errors, types.ValidationError{
Field: "title", Message: "Title is required",
})
}
if len(errors) > 0 {
return ctx.JSON(400, types.ValidationErrorResponse{Errors: errors})
}
return nil
}
```
### Predefined Common Errors
Ready-to-use error constants for common scenarios:
```go
import "base/core/types"
// Authentication and authorization errors
if err := validateToken(token); err != nil {
return types.ErrInvalidToken
}
if time.Now().After(claims.ExpiresAt.Time) {
return types.ErrTokenExpired
}
// User management errors
user, err := getUserByEmail(email)
if err != nil {
return types.ErrUserNotFound
}
if !bcrypt.CompareHashAndPassword(user.Password, []byte(password)) {
return types.ErrInvalidPassword
}
// Registration validation
if userExists(email) {
return types.ErrEmailExists
}
if !isValidEmail(email) {
return types.ErrInvalidEmail
}
// Custom error wrapping
func (s *AuthService) Login(email, password string) error {
user, err := s.GetUserByEmail(email)
if err != nil {
s.Logger.Error("User lookup failed", logger.String("email", email))
return fmt.Errorf("login failed: %w", types.ErrUserNotFound)
}
if !s.ValidatePassword(user, password) {
s.Logger.Warn("Invalid password attempt", logger.String("email", email))
return fmt.Errorf("login failed: %w", types.ErrInvalidPassword)
}
return nil
}
```
### Standard User Data Type
Standardized user data structure for consistent user representation:
```go
import "base/core/types"
// Standard user data structure
userData := types.UserData{
Id:        user.ID,
FirstName: user.FirstName,
LastName:  user.LastName,
Email:     user.Email,
Username:  user.Username,
}
// Usage in responses
func (c *AuthController) GetProfile(ctx *router.Context) error {
userID := c.GetContextUint(ctx, "user_id")
user, err := c.service.GetUser(userID)
if err != nil {
return c.RespondNotFound(ctx, "User")
}
// Convert to standard format
userData := types.UserData{
Id:        user.ID,
FirstName: user.FirstName,
LastName:  user.LastName,
Email:     user.Email,
Username:  user.Username,
}
return c.RespondSuccess(ctx, userData)
}
// JSON output
{
"success": true,
"data": {
"id": 123,
"firstname": "John",
"lastname": "Doe",
"email": "john@example.com",
"username": "johndoe"
}
}
```
## Best Practices
### ✅ Do
- • Use base classes to inherit common functionality
- • Register models for dynamic retrieval early in initialization
- • Use standardized response types for consistency
- • Leverage transaction helpers for data integrity
- • Use context helpers for type-safe value extraction
- • Generate unique slugs to avoid conflicts
- • Use the DateTime type for flexible time handling
- • Use predefined error constants for common scenarios
- • Structure validation errors with field-specific messages
### ❌ Don't
- • Reimplement functionality already provided by helpers
- • Forget to register models before using dynamic retrieval
- • Skip error handling in helper functions
- • Use raw time.Time when DateTime provides better features
- • Hardcode response structures instead of using types
- • Ignore transaction boundaries for related operations
- • Access context values without type checking
- • Create custom error types for common error scenarios
- • Return generic errors without field-specific validation messages
## Integration Example
Here's how to use multiple helpers together in a real module:
```go
// Module with integrated helpers
type BlogModule struct {
*base.Module
service    *BlogService
controller *BlogController
}
func NewBlogModule(deps base.ModuleDependencies) *BlogModule {
// Register models for dynamic retrieval
helper.RegisterModel("categories", func() interface{} { return &models.Category{} })
helper.RegisterModel("users", func() interface{} { return &models.User{} })
// Create services with base functionality
service := &BlogService{
Service: base.NewService(deps.GetDB(), deps.GetLogger(), deps.GetEmitter(), deps.GetStorage()),
slugHelper: helper.NewSlugHelper(),
}
// Create controller with base functionality
controller := &BlogController{
Controller: base.NewController(deps.GetLogger(), deps.GetStorage()),
service:    service,
}
return &BlogModule{
Module:     base.NewModule("blog", deps),
service:    service,
controller: controller,
}
}
type BlogService struct {
*base.Service
slugHelper *helper.SlugHelper
}
func (s *BlogService) CreatePost(data *CreatePostRequest) (*models.Post, error) {
// Validate relationships using dynamic helpers
category, err := helper.GetObjectAs[models.Category](s.DB, "category_id", data.CategoryID)
if err != nil {
return nil, fmt.Errorf("invalid category: %w", err)
}
author, err := helper.GetObjectAs[models.User](s.DB, "user_id", data.AuthorID)
if err != nil {
return nil, fmt.Errorf("invalid author: %w", err)
}
// Generate unique slug
baseSlug := s.slugHelper.Normalize(data.Title, data.Slug, "en")
uniqueSlug, err := s.slugHelper.GenerateUniqueSlug(baseSlug, func(slug string) (bool, error) {
var count int64
err := s.DB.Model(&models.Post{}).Where("slug = ?", slug).Count(&count).Error
return count > 0, err
})
if err != nil {
return nil, err
}
// Create post with transaction
var post *models.Post
err = s.WithTransaction(func(tx *gorm.DB) error {
post = &models.Post{
Title:      data.Title,
Content:    data.Content,
Slug:       uniqueSlug,
CategoryID: category.ID,
AuthorID:   author.ID,
CreatedAt:  types.Now(),
}
return tx.Create(post).Error
})
if err != nil {
s.LogError("create_post", err,
logger.String("title", data.Title),
logger.String("slug", uniqueSlug))
return nil, err
}
// Emit event and log success
s.EmitEvent("blog.post.created", post)
s.LogInfo("create_post", "Post created successfully",
logger.Uint("post_id", post.ID),
logger.String("slug", uniqueSlug))
return post, nil
}
```