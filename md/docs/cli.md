---
title: Cli
description: Cli documentation for Base Framework.
---

# CLI Reference
Complete reference for Base Framework CLI commands, generators, and field types.
## Overview
The Base CLI provides a comprehensive set of commands for creating projects, generating modules, managing development servers, and maintaining your applications. Built with Cobra for powerful command-line interactions.
## Quick Reference
### Project Commands
- • `base new` - Create new project
- • `base start` - Start development server
- • `base docs` - Generate documentation
### Module Commands
- • `base g` - Generate module
- • `base d` - Destroy module
- • `base update` - Update framework
## base new
Creates a new Base Framework project by downloading the latest framework template and setting up the project structure.
### Syntax
```bash
base new [project_name]
```
### What it does
- • Downloads the latest Base Framework template from GitHub
- • Creates a new directory with your project name
- • Sets up the complete project structure
- • Initializes Go modules and dependencies
- • Configures development environment
### Example Usage
```bash
# Create a new project called "my-api"
base new my-api
# Navigate to the project
cd my-api
# Start development server
base start
```
## base g (generate)
Generates complete modules with models, controllers, services, validators, and auto-updates your app initialization. The most powerful command in the Base CLI.
### Syntax
```bash
base g [module_name] [field:type...] [options]
base generate [module_name] [field:type...] [options]
```
### What gets generated
- • **Model:** `app/models/{name}.go` - GORM model with relationships
- • **Controller:** `app/{name}/controller.go` - HTTP handlers and routes
- • **Service:** `app/{name}/service.go` - Business logic layer
- • **Module:** `app/{name}/module.go` - Module definition and wiring
- • **Validator:** `app/{name}/validator.go` - Request validation
- • **Registration:** Updates `app/init.go` automatically
### Example Usage
```bash
# Simple blog post
base g post title:string content:text published:bool
# E-commerce product with relationships
base g product name:string price:float category_id:uint image:image
# User profile with file upload
base g profile user_id:uint avatar:image bio:text social_links:json
# Article with complex relationships
base g article title:string content:text author:belongsTo:User category:belongsTo:Category tags:manyToMany:Tag
```
### Field Types Reference
#### Basic Data Types
CLI Type
Go Type
Database
Example
string
string
VARCHAR(255)
name:string
text
string
TEXT
content:text
email
string
VARCHAR(255)
email:email
url
string
VARCHAR(255)
website:url
slug
string
VARCHAR(255)
slug:slug
#### Numeric Types
CLI Type
Go Type
Use Case
Example
int
int
Signed integers
age:int
uint
uint
IDs, positive numbers
user_id:uint
float
float64
Prices, measurements
price:float
decimal
float64
Financial calculations
amount:decimal
sort
int
Sort order
order:sort
bool
bool
True/false flags
active:bool
#### Date & Time Types
CLI Type
Go Type
Description
Example
datetime
time.Time
Full date and time
created_at:datetime
time
time.Time
Timestamp
started_at:time
date
time.Time
Date only
birth_date:date
timestamp
time.Time
Unix timestamp
modified:timestamp
#### File & Media Types
CLI Type
Go Type
Validation
Example
image
*storage.Attachment
5MB, image formats
avatar:image
file
*storage.Attachment
50MB, documents
document:file
#### Special & Data Types
CLI Type
Go Type
Use Case
Example
json
datatypes.JSON
JSON data
metadata:json
jsonb
datatypes.JSON
Binary JSON (PostgreSQL)
settings:jsonb
translation
translation.Field
Multi-language text
title:translation
translatedField
translation.Field
i18n field
name:translatedField
#### Relationship Types
Base Framework supports sophisticated relationship definitions using colon syntax:
Relationship
Syntax
Go Type
Example
Belongs To
field:belongsTo:Model
Model
author:belongsTo:User
Has One
field:hasOne:Model
*Model
profile:hasOne:Profile
Has Many
field:hasMany:Model
[]Model
posts:hasMany:Post
Many to Many
field:manyToMany:Model
[]*Model
tags:manyToMany:Tag
### Real-World Examples
#### Blog System
```bash
# Create blog models with relationships
base g category name:string slug:string description:text
base g tag name:string color:string
base g post title:string slug:slug content:text excerpt:text published:bool author:belongsTo:User category:belongsTo:Category tags:manyToMany:Tag featured_image:image
base g comment content:text author:belongsTo:User post:belongsTo:Post
```
#### E-commerce System
```bash
# E-commerce product catalog
base g brand name:string logo:image website:url
base g category name:string parent:belongsTo:Category
base g product name:string sku:string price:decimal description:text images:json brand:belongsTo:Brand category:belongsTo:Category
base g review rating:int comment:text user:belongsTo:User product:belongsTo:Product
```
#### Project Management
```bash
# Project management system
base g project name:string description:text start_date:date end_date:date status:string owner:belongsTo:User
base g task title:string description:text priority:int due_date:datetime completed:bool project:belongsTo:Project assignee:belongsTo:User
base g time_entry hours:decimal description:text date:date task:belongsTo:Task user:belongsTo:User
```
## base d (destroy)
Safely removes generated modules and cleans up all associated files, models, and registrations.
### Syntax
```bash
base d [module_name1] [module_name2] ...
base destroy [module_name1] [module_name2] ...
```
### What gets removed
- • **Module Directory:** `app/{module_name}/` and all contents
- • **Model File:** `app/models/{module_name}.go`
- • **Test Files:** `test/app_test/{module_name}_test/`
- • **Registration:** Import and initialization from `app/init.go`
### Safety Features
- • **Confirmation:** Prompts before destroying modules
- • **Batch Support:** Can destroy multiple modules at once
- • **Orphan Cleanup:** Removes registry entries even if files are missing
- • **Status Reporting:** Shows success/failure for each operation
### Example Usage
```bash
# Destroy a single module
base d post
# Destroy multiple modules
base d post comment category
# Full command name
base destroy user profile
```
### Example Output
```bash
$ base d post comment
Modules to destroy: post, comment
Are you sure you want to destroy 2 module(s)? [Y/n] y
[1/2] Destroying module 'post'...
✅ Removed module directory: app/post
✅ Removed model file: app/models/post.go
✅ Removed 'post' from app/init.go
✅ Successfully destroyed module 'post'
[2/2] Destroying module 'comment'...
✅ Removed module directory: app/comment
✅ Removed model file: app/models/comment.go
✅ Removed 'comment' from app/init.go
✅ Successfully destroyed module 'comment'
✅ Successfully destroyed all 2 module(s).
```
## base start
Starts your Base Framework application server with optional hot reloading and documentation generation.
### Syntax
```bash
base start [flags]
```
### Flags
Flag
Short
Description
--docs
-d
Generate Swagger documentation before starting
### What it does
- • Ensures dependencies are up to date with `go mod tidy`
- • Validates project structure (looks for main.go)
- • Generates Swagger documentation (if -d flag used)
- • Starts the application server
- • Sets environment variables for documentation
### Example Usage
```bash
# Start server normally
base start
# Start with documentation generation
base start -d
base start --docs
```
## base docs
Generates Swagger 2.0 documentation using go-swagger by scanning controller annotations and creates static files for your API.
### Syntax
```bash
base docs [flags]
```
### Flags
Flag
Short
Description
--output
-o
Output directory for generated files (default: "docs")
--static
-s
Generate static swagger files (default: true)
--no-static
Skip generating static files
### Generated Files
- • **swagger.json** - Swagger 2.0 specification in JSON format
- • **swagger.yaml** - Swagger 2.0 specification in YAML format
- • **docs.go** - Go package with embedded documentation
### Supported Annotations
```go
// @Summary Create a new post
// @Description Create a new blog post with title and content
// @Tags posts
// @Param data body CreatePostRequest true "Post data"
// @Success 201 "Post created successfully"
// @Failure 400 "Bad request"
// @Router /posts [post]
// @Security BearerAuth
func (c *PostController) CreatePost(ctx *router.Context) error {
// Your implementation here
}
```
### Example Usage
```bash
# Generate docs in default directory
base docs
# Generate docs in custom directory
base docs -o documentation
# Skip static file generation
base docs --no-static
# Generate with custom output
base docs -o api-docs -s
```
## base update
Updates the Base Framework core directory to the latest version while preserving your application code and customizations.
### Syntax
```bash
base update
```
### What it does
- • Downloads the latest Base Framework core matching your CLI version
- • Creates automatic backup of existing core directory
- • Replaces core directory with updated version
- • Preserves your app directory and custom code
- • Rollback support if update fails
### Safety Features
- • **Automatic Backup:** Creates `core.bak` before updating
- • **Validation:** Verifies downloaded files before replacement
- • **Rollback:** Automatically restores backup if update fails
- • **Preservation:** Only updates core, leaves app code untouched
### Example Output
```bash
$ base update
Updating Base Core...
Downloading core from: https://github.com/base-go/base-core/archive/refs/tags/v2.1.7.zip
Core directory updated successfully.
Base Core updated successfully.
```
## base upgrade
Upgrades the Base CLI tool itself to the latest version with support for major version upgrades and breaking change warnings.
### Syntax
```bash
base upgrade [flags]
```
### Flags
Flag
Description
--major
Allow upgrade to new major version (may contain breaking changes)
### Upgrade Behavior
- • **Minor/Patch:** Upgrades automatically within same major version
- • **Major:** Requires `--major` flag and user confirmation
- • **Validation:** Downloads, verifies, and tests binary before installation
- • **Cross-platform:** Handles Windows, macOS, and Linux installations
- • **Permissions:** Uses sudo for system directory installation on Unix
### Major Version Warnings
```bash
# When only minor/patch updates are available
$ base upgrade
Checking for updates...
Downloading version 2.0.9...
Successfully upgraded to version 2.0.9!
# When major version is available without --major flag
$ base upgrade
Checking for updates...
You're already using the latest version (2.0.9)
# When using --major flag for major version upgrade
$ base upgrade --major
Checking for updates...
🚨 MAJOR VERSION UPGRADE DETECTED: 2.0.9 → 3.0.0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  This is a MAJOR version upgrade that may contain breaking changes!
📚 Full changelog: https://github.com/base-go/cmd/releases/tag/v3.0.0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Do you want to proceed with the upgrade? [y/N]: y
Downloading version 3.0.0...
Successfully upgraded to version 3.0.0!
```
### Example Usage
```bash
# Upgrade within current major version
base upgrade
# Allow major version upgrade with breaking changes
base upgrade --major
# Check what version would be installed
base version
```