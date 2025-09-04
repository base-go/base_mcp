---
title: Middleware
description: Middleware documentation for Base Framework.
---

# 🛡️ Configurable Middleware System
Convention over Configuration approach to middleware management with environment-driven settings and module-level overrides.
## 📋 Overview
Base Framework's configurable middleware system follows "Convention over Configuration" principles, providing sensible defaults that work out-of-the-box while allowing fine-grained control when needed. The system supports environment-based configuration, path-based rules, and module-level overrides.
### 🎯 Key Benefits
- • **Zero Configuration:** Works immediately with sensible defaults
- • **Environment Driven:** Configure via environment variables
- • **Path-Based Rules:** Different middleware for different endpoints
- • **Module Overrides:** Modules can customize their middleware
- • **Webhook Support:** Special handling for third-party webhooks
### ⚙️ Configuration Layers
Module-specific overrides
Highest
Route-specific configuration
High
Global environment settings
Medium
Framework defaults
Lowest
## 🌍 Environment Configuration
Configure middleware behavior through environment variables in your `.env` file. The system provides sensible defaults for immediate use.
### 🌐 Global Middleware Settings
```yml
# Global middleware toggles
MIDDLEWARE_API_KEY_ENABLED=true
MIDDLEWARE_AUTH_ENABLED=false
MIDDLEWARE_RATE_LIMIT_ENABLED=true
MIDDLEWARE_LOGGING_ENABLED=true
MIDDLEWARE_RECOVERY_ENABLED=true
MIDDLEWARE_CORS_ENABLED=true
```
#### Security Middleware
- • `API_KEY_ENABLED`: Require X-Api-Key header
- • `AUTH_ENABLED`: Require JWT/Bearer tokens
- • `CORS_ENABLED`: Enable CORS headers
#### Performance & Monitoring
- • `RATE_LIMIT_ENABLED`: Enable request throttling
- • `LOGGING_ENABLED`: Enable request logging
- • `RECOVERY_ENABLED`: Enable panic recovery
### 🛤️ Path-Based Skip Rules
```yml
# Skip middleware for specific paths (comma-separated)
MIDDLEWARE_API_KEY_SKIP_PATHS=/health,/,/docs,/swagger
MIDDLEWARE_AUTH_SKIP_PATHS=/api/auth/login,/api/auth/register,/api/auth/forgot-password
MIDDLEWARE_RATE_LIMIT_SKIP_PATHS=/health,/
MIDDLEWARE_LOGGING_SKIP_PATHS=
```
💡 Wildcard Support
Use wildcards for path patterns: `/api/public/*` matches all public endpoints, `/webhooks/*` matches all webhook endpoints.
### ⚡ Rate Limiting Configuration
```yml
# Rate limiting settings
MIDDLEWARE_RATE_LIMIT_REQUESTS=60
MIDDLEWARE_RATE_LIMIT_WINDOW=1m
# Webhook-specific rate limiting
MIDDLEWARE_WEBHOOK_RATE_LIMIT_REQUESTS=1000
MIDDLEWARE_WEBHOOK_RATE_LIMIT_WINDOW=1h
```
#### Standard API Endpoints
- • **60 requests** per minute (default)
- • Applied to all API endpoints
- • Based on client IP address
#### Webhook Endpoints
- • **1000 requests** per hour (default)
- • Higher limits for webhook traffic
- • Prevents webhook spam/abuse
## 🎣 Webhook-Specific Configuration
Webhooks from third-party services (Stripe, GitHub, PayPal) have different security requirements than regular API endpoints. The system provides special webhook handling.
### 🔧 Webhook Settings
```yml
# Webhook paths and security settings
MIDDLEWARE_WEBHOOK_PATHS=/api/webhooks/*,/webhooks/*
MIDDLEWARE_WEBHOOK_API_KEY_ENABLED=false
MIDDLEWARE_WEBHOOK_AUTH_ENABLED=false
MIDDLEWARE_WEBHOOK_SIGNATURE_ENABLED=true
```
#### 🎯 Webhook Security Model
- • **No API Key:** Third parties don't have your internal API keys
- • **No Authentication:** No JWT tokens from external services
- • **Signature Verification:** HMAC-based webhook signatures
- • **Rate Limiting:** Higher limits to handle webhook bursts
### 🔐 Webhook Signature Secrets
```yml
# Webhook signature verification secrets
STRIPE_WEBHOOK_SECRET=whsec_your_stripe_webhook_secret
GITHUB_WEBHOOK_SECRET=your_github_webhook_secret
PAYPAL_WEBHOOK_SECRET=your_paypal_webhook_secret
```
🔒 Security Note
Always use environment variables for webhook secrets. Never hardcode secrets in your source code. Each webhook provider has a different signature format and verification process.
## 🧩 Module-Level Overrides
Modules can override global middleware settings for their specific routes. This provides fine-grained control while maintaining global defaults.
### 📋 ConfigurableModule Interface
```go
// Extend your module with middleware configuration
type ConfigurableModule interface {
Module
// Return middleware overrides for this module
MiddlewareConfig() *MiddlewareOverrides
}
```
Implement the `ConfigurableModule` interface to provide custom middleware settings for your module's routes.
### 🎣 Webhook Module Example
```go
func (m *WebhooksModule) MiddlewareConfig() *module.MiddlewareOverrides {
return &module.MiddlewareOverrides{
PathRules: map[string]module.MiddlewareSettings{
"/api/webhooks/stripe": {
APIKey: module.DisableAPIKey().APIKey,
Auth:   module.DisableAuth().Auth,
WebhookSignature: module.WebhookSignature(
"stripe", "stripe-signature", "STRIPE_WEBHOOK_SECRET",
).WebhookSignature,
RateLimit: module.CustomRateLimit(500, "1h").RateLimit,
},
"/api/webhooks/github": {
APIKey: module.DisableAPIKey().APIKey,
Auth:   module.DisableAuth().Auth,
WebhookSignature: module.WebhookSignature(
"github", "x-hub-signature-256", "GITHUB_WEBHOOK_SECRET",
).WebhookSignature,
},
},
}
}
```
#### Helper Functions
- • `DisableAPIKey()`: Skip API key requirement
- • `DisableAuth()`: Skip authentication
- • `CustomRateLimit()`: Custom rate limiting
- • `WebhookSignature()`: Signature verification
#### Override Behavior
- • `nil`: Use global setting
- • `true`: Enable middleware
- • `false`: Disable middleware
- • `config`: Custom configuration
### 🔐 Authentication Module Example
```go
func (m *AuthModule) MiddlewareConfig() *module.MiddlewareOverrides {
return &module.MiddlewareOverrides{
PathRules: map[string]module.MiddlewareSettings{
// Public auth endpoints - no API key or auth required
"/api/auth/login":    module.DisableAuthAndAPIKey(),
"/api/auth/register": module.DisableAuthAndAPIKey(),
"/api/auth/forgot-password": module.DisableAuthAndAPIKey(),
// Protected auth endpoints - require both API key and auth
"/api/auth/profile":  module.RequireAuth(),
"/api/auth/logout":   module.RequireAuth(),
},
}
}
```
Authentication modules typically disable security for login/register endpoints while requiring it for protected operations.
## ⚙️ Advanced Configuration
### 📝 JSON Configuration Overrides
For complex routing rules, use JSON configuration in environment variables:
```yml
# Complex middleware overrides (JSON format)
MIDDLEWARE_OVERRIDES={
"api/admin/*": {
"api_key": "required",
"auth": "required",
"rate_limit": "strict"
},
"api/public/*": {
"api_key": "disabled",
"auth": "disabled"
},
"api/webhooks/*": {
"api_key": "disabled",
"auth": "disabled",
"rate_limit": "relaxed"
}
}
```
💡 JSON vs Module Overrides
Use JSON overrides for simple path-based rules. Use module overrides for complex logic, custom rate limiting, and webhook signature verification.
### 🚀 Production Considerations
```yml
# Production middleware settings
MIDDLEWARE_API_KEY_ENABLED=true
MIDDLEWARE_AUTH_ENABLED=true
MIDDLEWARE_RATE_LIMIT_ENABLED=true
MIDDLEWARE_RATE_LIMIT_REQUESTS=100
MIDDLEWARE_RATE_LIMIT_WINDOW=1h
MIDDLEWARE_LOGGING_ENABLED=true
MIDDLEWARE_RECOVERY_ENABLED=true
```
#### Security Best Practices
- • Enable API key validation in production
- • Require authentication for sensitive endpoints
- • Use stricter rate limiting
- • Always enable recovery middleware
#### Performance Optimization
- • Disable logging for high-traffic endpoints
- • Use appropriate rate limit windows
- • Skip middleware for health checks
- • Monitor middleware performance
## 🎯 Common Use Cases
### 🌐 API Gateway Pattern
```yml
# All endpoints require API key
MIDDLEWARE_API_KEY_ENABLED=true
MIDDLEWARE_API_KEY_SKIP=/health,/docs
# Auth required except for public endpoints
MIDDLEWARE_AUTH_ENABLED=true
MIDDLEWARE_AUTH_SKIP=/api/public/*,/api/auth/*
```
Secure all endpoints by default, with explicit exceptions for public routes.
### 🔧 Microservice Pattern
```yml
# Minimal middleware for internal services
MIDDLEWARE_API_KEY_ENABLED=false
MIDDLEWARE_AUTH_ENABLED=false
MIDDLEWARE_RATE_LIMIT_ENABLED=false
MIDDLEWARE_LOGGING_ENABLED=true
```
Lightweight configuration for internal microservices with just logging enabled.
### 🌍 Public API Pattern
```yml
# Strict rate limiting for public APIs
MIDDLEWARE_RATE_LIMIT_ENABLED=true
MIDDLEWARE_RATE_LIMIT_REQUESTS=100
MIDDLEWARE_RATE_LIMIT_WINDOW=1h
MIDDLEWARE_API_KEY_ENABLED=true
```
Conservative rate limiting and API key requirements for public-facing APIs.
### 🛠️ Development Pattern
```yml
# Relaxed settings for development
MIDDLEWARE_API_KEY_ENABLED=false
MIDDLEWARE_AUTH_ENABLED=false
MIDDLEWARE_RATE_LIMIT_ENABLED=false
MIDDLEWARE_LOGGING_ENABLED=true
```
Minimal restrictions during development with comprehensive logging for debugging.
## 🔄 Migration Guide
Migrating from manual middleware configuration to the configurable system is straightforward and backward compatible.
### 📋 Migration Steps
1
#### Add Environment Configuration
Add middleware configuration to your `.env` file using the examples above.
2
#### Update main.go
Replace manual middleware setup with `ApplyConfigurableMiddleware()` call.
3
#### Implement Module Overrides
Add `MiddlewareConfig()` method to modules that need custom middleware.
4
#### Test Configuration
Verify middleware behavior matches your expectations across different endpoints.
✅ Backward Compatibility
The configurable middleware system is fully backward compatible. Existing middleware configurations will continue to work while you gradually migrate to the new system.