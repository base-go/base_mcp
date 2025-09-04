---
title: Validator
description: Validator documentation for Base Framework.
---

# Validator
Powerful request validation system with custom error messages and struct tag validation
## Overview
🔍 Comprehensive Validation
Base Framework's validator wraps go-playground/validator v10 with enhanced error messages, automatic JSON field name mapping, and seamless integration with request handling. It provides over 100+ built-in validation rules and supports custom validators.
Struct Validation
Validate entire structs using struct tags with automatic field name resolution.
Custom Messages
Automatically generated user-friendly error messages with field-specific context.
Custom Rules
Create and register custom validation rules for domain-specific validation logic.
## Quick Start
Basic Usage
Basic Validation Example
```
package main
import (
"base/core/validator"
)
type UserCreateRequest struct {
Name     string `json:"name" validate:"required,min=2,max=50"`
Email    string `json:"email" validate:"required,email"`
Age      int    `json:"age" validate:"gte=18,lte=120"`
Password string `json:"password" validate:"required,min=8"`
}
func validateUser() {
req := &UserCreateRequest{
Name:     "",           // Invalid - required
Email:    "invalid",   // Invalid - not email format
Age:      15,          // Invalid - too young
Password: "123",       // Invalid - too short
}
// Validate using Base validator
if errors := validator.Validate(req); errors != nil {
// Handle validation errors
for _, err := range errors {
fmt.Printf("%s: %s\n", err.Field, err.Message)
}
}
}
```
## Built-in Validation Rules
Base Framework includes over 100+ validation rules from go-playground/validator. Here are the most commonly used ones:
Required & Optional
#### Basic Rules
`required` - Field cannot be empty
`omitempty` - Skip validation if empty
`required_if` - Required based on other field
`required_unless` - Required unless other field matches
#### Examples
``validate:"required"``
``validate:"omitempty,email"``
``validate:"required_if=Type user"``
``validate:"required_unless=Premium true"``
String Validation
#### Length & Format
`min=n` - Minimum length
`max=n` - Maximum length
`len=n` - Exact length
`alpha` - Letters only
`alphanum` - Letters and numbers
`numeric` - Numbers only
#### Pattern Matching
`email` - Valid email format
`url` - Valid URL format
`uuid` - Valid UUID format
`contains=text` - Contains substring
`startswith=text` - Starts with text
`endswith=text` - Ends with text
Numeric Validation
#### Comparison
`gt=n` - Greater than
`gte=n` - Greater than or equal
`lt=n` - Less than
`lte=n` - Less than or equal
`eq=n` - Equal to
`ne=n` - Not equal to
#### Examples
``validate:"gte=18,lte=120"``
``validate:"gt=0"``
``validate:"min=1,max=100"``
``validate:"oneof=1 2 3 4 5"``
Collections & Arrays
#### Array/Slice Rules
`dive` - Validate each element
`unique` - All elements unique
`min=n` - Minimum array length
`max=n` - Maximum array length
#### Example Usage
``validate:"required,dive,email"``
``validate:"min=1,max=10,dive,gt=0"``
``validate:"unique"``
## Request Validation in Controllers
Model Request Structure
models/user.go
```
package models
import (
"time"
"gorm.io/gorm"
)
// User model
type User struct {
ID        uint           `json:"id" gorm:"primarykey"`
CreatedAt time.Time      `json:"created_at"`
UpdatedAt time.Time      `json:"updated_at"`
DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
Name      string         `json:"name" gorm:"not null"`
Email     string         `json:"email" gorm:"uniqueIndex"`
Age       int            `json:"age"`
Role      string         `json:"role" gorm:"default:user"`
}
// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
Name     string `json:"name" validate:"required,min=2,max=50"`
Email    string `json:"email" validate:"required,email"`
Age      int    `json:"age" validate:"gte=18,lte=120"`
Role     string `json:"role" validate:"omitempty,oneof=user admin moderator"`
Password string `json:"password" validate:"required,min=8,containsany=!@#$%^&*"`
}
// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
Name string `json:"name" validate:"omitempty,min=2,max=50"`
Age  *int   `json:"age" validate:"omitempty,gte=18,lte=120"`
Role string `json:"role" validate:"omitempty,oneof=user admin moderator"`
}
```
Controller Validation
user/controller.go
```
package user
import (
"net/http"
"base/app/models"
"base/core/router"
"base/core/types"
"base/core/validator"
)
func (c *UserController) Create(ctx *router.Context) error {
var req models.CreateUserRequest
// Bind JSON request
if err := ctx.ShouldBindJSON(&req); err != nil {
return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
Error:   "Invalid request format",
Success: false,
Details: err.Error(),
})
}
// Validate request using Base validator
if validationErrors := validator.Validate(&req); validationErrors != nil {
return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
Error:   "Validation failed",
Success: false,
Details: validationErrors,
})
}
// Create user through service
user, err := c.Service.Create(&req)
if err != nil {
return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
Error:   "Failed to create user",
Success: false,
Details: err.Error(),
})
}
return ctx.JSON(http.StatusCreated, types.SuccessResponse{
Message: "User created successfully",
Success: true,
Data:    user.ToResponse(),
})
}
```
## Custom Validation Rules
Creating Custom Validators
Custom Validator Example
```
package validators
import (
"regexp"
"strings"
"base/core/validator"
"github.com/go-playground/validator/v10"
)
// Init registers all custom validators
func Init() {
v := validator.New()
// Register custom validators
v.RegisterValidation("username", validateUsername)
v.RegisterValidation("strong_password", validateStrongPassword)
v.RegisterValidation("phone", validatePhone)
v.RegisterValidation("slug", validateSlug)
}
// validateUsername checks if username is valid (alphanumeric, underscore, hyphen)
func validateUsername(fl validator.FieldLevel) bool {
username := fl.Field().String()
if len(username) 3 || len(username) > 20 {
return false
}
// Only allow alphanumeric, underscore, and hyphen
matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
return matched
}
// validateStrongPassword checks password complexity
func validateStrongPassword(fl validator.FieldLevel) bool {
password := fl.Field().String()
// At least 8 characters
if len(password) 8 {
return false
}
// Must contain uppercase, lowercase, number, and special character
hasUpper, _ := regexp.MatchString("[A-Z]", password)
hasLower, _ := regexp.MatchString("[a-z]", password)
hasNumber, _ := regexp.MatchString("[0-9]", password)
hasSpecial, _ := regexp.MatchString("[!@#$%^&*(),.?\":{}|<>]", password)
return hasUpper && hasLower && hasNumber && hasSpecial
}
```
Using Custom Validators
Applying Custom Validation Tags
```
type UserRegistrationRequest struct {
Username string `json:"username" validate:"required,username"`
Email    string `json:"email" validate:"required,email"`
Phone    string `json:"phone" validate:"omitempty,phone"`
Password string `json:"password" validate:"required,strong_password"`
}
type BlogPostRequest struct {
Title   string `json:"title" validate:"required,min=3,max=200"`
Slug    string `json:"slug" validate:"required,slug"`
Content string `json:"content" validate:"required,min=10"`
Tags    []string `json:"tags" validate:"omitempty,dive,min=2,max=20"`
}
// Initialize validators in your main.go or init function
func init() {
validators.Init()
}
```
## Error Handling & Messages
🎯 User-Friendly Messages
Base Framework automatically generates user-friendly error messages from validation tags. Field names use JSON tag names for consistency with API responses.
Validation Error Structure
Error Response Example
```
{
"error": "Validation failed",
"success": false,
"details": [
{
"field": "name",
"tag": "required",
"value": "",
"message": "name is required"
},
{
"field": "email",
"tag": "email",
"value": "invalid-email",
"message": "email must be a valid email address"
},
{
"field": "age",
"tag": "gte",
"value": "15",
"message": "age must be greater than or equal to 18"
},
{
"field": "password",
"tag": "min",
"value": "123",
"message": "password must be at least 8 characters long"
}
]
}
```
Custom Error Messages
Enhanced Error Message Handling
```
package validators
import (
"fmt"
"base/core/validator"
)
// CustomMessageValidator extends the base validator with custom messages
type CustomMessageValidator struct {
*validator.Validator
messages map[string]string
}
func NewCustomValidator() *CustomMessageValidator {
return &CustomMessageValidator{
Validator: validator.New(),
messages: map[string]string{
"required":      "This field is required",
"email":         "Please enter a valid email address",
"min":           "This field must be at least %s characters",
"max":           "This field cannot exceed %s characters",
"gte":           "Value must be at least %s",
"lte":           "Value cannot exceed %s",
"username":      "Username must be 3-20 characters with letters, numbers, underscore or hyphen only",
"strong_password": "Password must be at least 8 characters with uppercase, lowercase, number and special character",
},
}
}
```
## Validation with File Uploads
Base Framework integrates validation seamlessly with file uploads, providing automatic file type and size validation.
File Upload Validation
File Upload with Validation
```
type ProductCreateRequest struct {
Name        string                `json:"name" validate:"required,min=2,max=100"`
Description string                `json:"description" validate:"required,min=10,max=1000"`
Price       float64               `json:"price" validate:"required,gt=0"`
Category    string                `json:"category" validate:"required,oneof=electronics clothing books"`
Images      []*multipart.FileHeader `json:"images" validate:"omitempty,dive,max=5242880"` // 5MB limit
Manual      *multipart.FileHeader   `json:"manual" validate:"omitempty"`
}
func (c *ProductController) CreateWithFiles(ctx *router.Context) error {
var req ProductCreateRequest
// Parse multipart form
if err := ctx.Request.ParseMultipartForm(32 20); err != nil { // 32MB
return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
Error: "Failed to parse form data",
})
}
// Bind form data
req.Name = ctx.FormValue("name")
req.Description = ctx.FormValue("description")
// Validate request
if validationErrors := validator.Validate(&req); validationErrors != nil {
return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
Error:   "Validation failed",
Details: validationErrors,
})
}
// Process upload through service
product, err := c.Service.CreateWithFiles(&req)
if err != nil {
return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
Error: "Failed to create product: " + err.Error(),
})
}
return ctx.JSON(http.StatusCreated, product.ToResponse())
}
```
## Advanced Validation Patterns
Cross-Field Validation
Cross-Field Validation Examples
```
type UserUpdateRequest struct {
CurrentPassword string `json:"current_password" validate:"required_with=NewPassword"`
NewPassword     string `json:"new_password" validate:"omitempty,min=8,nefield=CurrentPassword"`
ConfirmPassword string `json:"confirm_password" validate:"required_with=NewPassword,eqfield=NewPassword"`
}
type EventScheduleRequest struct {
StartDate time.Time `json:"start_date" validate:"required"`
EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
Duration  int       `json:"duration" validate:"required,gte=1"`
}
type PricingRequest struct {
RegularPrice float64 `json:"regular_price" validate:"required,gt=0"`
SalePrice    float64 `json:"sale_price" validate:"omitempty,gt=0,ltfield=RegularPrice"`
Currency     string  `json:"currency" validate:"required,oneof=USD EUR GBP"`
}
type ConditionalValidationRequest struct {
Type         string  `json:"type" validate:"required,oneof=personal business"`
CompanyName  string  `json:"company_name" validate:"required_if=Type business"`
TaxID        string  `json:"tax_id" validate:"required_if=Type business"`
PersonalName string  `json:"personal_name" validate:"required_if=Type personal"`
}
```
Nested Struct Validation
Nested Validation Example
```
type Address struct {
Street     string `json:"street" validate:"required,min=5,max=100"`
City       string `json:"city" validate:"required,min=2,max=50"`
State      string `json:"state" validate:"required,min=2,max=50"`
PostalCode string `json:"postal_code" validate:"required,len=5|len=9"`
Country    string `json:"country" validate:"required,len=2"`
}
type ContactInfo struct {
Email string `json:"email" validate:"required,email"`
Phone string `json:"phone" validate:"omitempty,phone"`
}
type UserRegistrationRequest struct {
Name        string      `json:"name" validate:"required,min=2,max=50"`
Contact     ContactInfo `json:"contact" validate:"required"`
Address     Address     `json:"address" validate:"required"`
ShippingAddress *Address `json:"shipping_address" validate:"omitempty"`
Preferences []string    `json:"preferences" validate:"omitempty,dive,oneof=email sms newsletter"`
Tags        []Tag       `json:"tags" validate:"omitempty,dive"`
}
type Tag struct {
Name  string `json:"name" validate:"required,min=2,max=20"`
Color string `json:"color" validate:"omitempty,hexcolor"`
}
```
## Best Practices
✅ Do's
- Use JSON tag names for consistent field naming in validation messages
- Combine multiple validation rules for comprehensive checks
- Use omitempty for optional fields to skip validation when empty
- Create custom validators for domain-specific business rules
- Validate at the controller level before processing requests
- Use dive tag for validating array/slice elements
- Implement cross-field validation for related fields
- Return structured validation errors with field-specific messages
❌ Don'ts
- Don't skip validation on user input
- Don't use validation tags on unexported fields
- Don't forget to handle validation errors properly
- Don't rely solely on client-side validation
- Don't validate in multiple layers unnecessarily
- Don't ignore the JSON tag for field naming
- Don't create overly complex validation rules in single tags
- Don't expose internal validation library errors to users
💡 Tips
- Use the Base CLI generator to automatically create validated request structs
- Test your validation rules with edge cases and boundary values
- Consider performance when using complex validation on large datasets
- Document custom validation rules for team members
- Use validation aliases for commonly repeated rule combinations
- Keep validation logic in the request layer, not in models or services