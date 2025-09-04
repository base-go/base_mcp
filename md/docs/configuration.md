---
title: Configuration
description: Configuration documentation for Base Framework.
---

# Configuration
Complete guide to configuring Base Framework with environment variables, security best practices, and deployment strategies.
## Quick Start
### 1. Create Your Environment File
```bash
# Copy the template
cp .env.sample .env
# Edit with your settings
vim .env
```
### 2. Security Best Practices
### 3. Required vs Optional Settings
#### Required for Development
- `SERVER_PORT` - Server port (default: 8100)
- `ENV` - Environment mode (default: debug)
#### Optional but Recommended
- `JWT_SECRET` - Change from default for security
- `API_KEY` - Change from default for security
- `DB_DRIVER` - Database type (default: sqlite)
## Environment File Structure
### .env (Your local configuration)
- ✅ Contains your actual configuration values
- 🚫 Not tracked by git (listed in .gitignore)
- 🔒 Safe to include sensitive data like passwords and API keys
- 📋 Created by copying .env.sample
### .env.sample (Template)
- 📝 Template with example values
- ✅ Tracked by git for team sharing
- 📚 Contains documentation and examples
- 🚫 Never contains real secrets
## Configuration Sections
### Application Settings
```
APP_VERSION=1.0.0
ENV=debug                    # debug, development, production
```
### Server Configuration
```
SERVER_ADDRESS=localhost
SERVER_PORT=8100
APPHOST=http://localhost:8100
CORS_ALLOWED_ORIGINS=http://localhost:3000
```
### Feature Toggles
```
SWAGGER_ENABLED=true         # Set to false in production
WS_ENABLED=true             # Enable WebSocket functionality
```
Feature toggles allow you to enable/disable functionality based on your deployment needs.
### Security ⚠️ Change in production!
```
JWT_SECRET=change_me_in_production
API_KEY=change_me_in_production
SWAGGER_ENABLED=false
```
### Database
```
DB_DRIVER=sqlite            # sqlite, mysql, postgres
DB_PATH=storage/base.db     # For SQLite
# For other databases:
# DB_HOST=localhost
# DB_PORT=3306
# DB_USER=username
# DB_PASSWORD=password
# DB_NAME=database_name
```
### Email
```
EMAIL_PROVIDER=smtp         # smtp, sendgrid, postmark
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
# SMTP settings
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_username
SMTP_PASSWORD=your_password
```
### Storage
```
STORAGE_PROVIDER=local      # local, s3, r2
STORAGE_PATH=storage/uploads
STORAGE_MAX_SIZE=10485760   # 10MB in bytes
STORAGE_ALLOWED_EXT=.jpg,.jpeg,.png,.gif,.pdf
```
## Environment-Specific Examples
```
ENV=debug
SWAGGER_ENABLED=true
WS_ENABLED=true
DB_DRIVER=sqlite
DB_PATH=storage/dev.db
LOG_LEVEL=debug
JWT_SECRET=dev_secret_not_for_production
API_KEY=dev_api_key
```
```
ENV=production
SWAGGER_ENABLED=false
WS_ENABLED=true
DB_DRIVER=postgres
DB_URL=postgres://user:pass@db:5432/app
LOG_LEVEL=warn
JWT_SECRET=super_secure_production_secret_256_bits
API_KEY=secure_production_api_key
STORAGE_PROVIDER=s3
STORAGE_BUCKET=production-assets
```
```
ENV=testing
SWAGGER_ENABLED=false
WS_ENABLED=false
DB_DRIVER=sqlite
DB_PATH=storage/test.db
LOG_LEVEL=error
```
## Validation and Troubleshooting
### Automatic Validation
The framework automatically validates your configuration and provides helpful error messages:
```
config := config.NewConfig()
if errors := config.Validate(); len(errors) > 0 {
for _, err := range errors {
log.Printf("Configuration error: %v", err)
}
}
```
### Common Issues
#### Port already in use
```
# Change the port
SERVER_PORT=8101
```
#### Database connection failed
```
# For SQLite, ensure directory exists
mkdir -p storage
# For other databases, check connection settings
DB_HOST=localhost
DB_PORT=3306
```
#### Invalid boolean values
✅ Correct:
```
SWAGGER_ENABLED=true
SWAGGER_ENABLED=false
```
❌ Invalid (will use default):
```
SWAGGER_ENABLED=yes
```
### Testing Configuration
```
# Test with specific settings
SWAGGER_ENABLED=false base start
# Test production-like settings
ENV=production SWAGGER_ENABLED=false LOG_LEVEL=warn base start
```
## Security Guidelines
### Development Security
- ✅ Use default values for local development
- 🔄 Rotate secrets regularly
- 🚫 Don't share .env files between developers
### Production Security
- 🔐 Use strong, unique secrets (256-bit minimum)
- 🚫 Disable development features (`SWAGGER_ENABLED=false`)
- ☁️ Use secure storage providers
- 📊 Set appropriate log levels
- 🔑 Use environment variables or secret management systems
### Secret Management
❌ Bad - hardcoded in .env
```
JWT_SECRET=mysecret123
```
✅ Good - from environment/vault
```
JWT_SECRET=${VAULT_JWT_SECRET}
JWT_SECRET=${K8S_SECRET_JWT}
```
## Docker Integration
### Dockerfile
```
# Copy template but not actual .env
COPY .env.sample /app/.env.sample
# Set production defaults
ENV SWAGGER_ENABLED=false
ENV WS_ENABLED=true
ENV ENV=production
```
### Docker Compose
```
services:
base-app:
image: base-framework
environment:
- SWAGGER_ENABLED=false
- DB_URL=postgres://user:pass@db:5432/app
env_file:
- .env.production  # Separate production env file
```
## Git Best Practices
### What's Tracked
- ✅ `.env.sample` - Template file
- ✅ `.gitignore` - Includes .env exclusion
- ✅ Documentation files
### What's NOT Tracked
- ❌ `.env` - Your local configuration
- ❌ `.env.local` - Local overrides
- ❌ `.env.production` - Production secrets
- ❌ `storage/base.db` - Local database
### Team Collaboration
1. Update `.env.sample` when adding new configuration options
2. Document new settings in comments
3. Notify team of new required environment variables
4. Use defaults that work for most developers
## Migration from Old Configuration
If you have an existing .env file tracked by git:
```
# Remove from git but keep local file
git rm --cached .env
# Add to gitignore (if not already there)
echo ".env" >> .gitignore
# Commit the removal
git add .gitignore
git commit -m "Remove .env from git tracking for security"
# Update team
echo "Team: Please copy .env.sample to .env and configure locally"
```
## Feature Control Examples
### Development (all features enabled)
```
SWAGGER_ENABLED=true WS_ENABLED=true base start
```
### Production (secure configuration)
```
ENV=production SWAGGER_ENABLED=false base start
```
### API-only deployment (minimal features)
```
SWAGGER_ENABLED=false WS_ENABLED=false base start
```
## Configuration Validation
Base Framework includes comprehensive configuration validation with helpful error messages and security warnings.
```
config := config.NewConfig()
if errors := config.Validate(); len(errors) > 0 {
for _, err := range errors {
log.Printf("Configuration error: %v", err)
}
}
```
#### Type Safety
Automatic type conversion with fallback to safe defaults for invalid values.
#### Security Checks
Warns about insecure defaults in production environments.
#### Helper Methods
Convenient methods for common operations and environment checks.