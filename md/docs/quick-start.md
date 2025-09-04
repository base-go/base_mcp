---
title: Quick Start
description: Quick Start documentation for Base Framework.
---

# Quick Start
Create your first Base Framework project in minutes with our step-by-step guide.
## Prerequisites
### Go 1.24+
Base Framework requires Go 1.24 or later for modern language features and performance.
### Base CLI
Install the Base CLI tool for project scaffolding and module generation.
## Create Your First Project
1
### Create New Project
Create a new Base Framework project and navigate to the directory:
```bash
# Create a new project
base new my-blog-api
cd my-blog-api
# Project structure is automatically created with:
# - Core framework files
# - Database configuration
# - Initial app structure
```
2
### Generate Your First Module
Generate a complete blog post module with models, controllers, and services:
```bash
# Generate a blog post module with various field types
base g post title:string content:text slug:slug published:bool author_id:uint
# This creates:
# - app/models/post.go (GORM model)
# - app/post/controller.go (HTTP handlers)
# - app/post/service.go (business logic)
# - app/post/module.go (module definition)
# - Updates app/init.go automatically
```
3
### Start Development Server
Start the server with hot reloading and automatic documentation generation:
```bash
# Start with hot reloading and docs
base docs
base start -r
# Your API is now running at:
# http://localhost:8080 - Main API
# http://localhost:8080/swagger/ - Interactive API docs
```
## Generated API Endpoints
Your blog post module automatically generates these REST endpoints:
GET
/api/posts
List all posts with pagination
POST
/api/posts
Create a new post
GET
/api/posts/{id}
Get post by ID
PUT
/api/posts/{id}
Update existing post
DELETE
/api/posts/{id}
Delete post
## Test Your API
### Create a Post
```bash
curl -X POST http://localhost:8080/api/posts \
-H "Content-Type: application/json" \
-d '{
"title": "My First Blog Post",
"content": "This is the content of my first blog post using Base Framework!",
"slug": "my-first-post",
"published": true,
"author_id": 1
}'
```
### List All Posts
```bash
curl http://localhost:8080/api/posts
# With pagination
curl "http://localhost:8080/api/posts?page=1&limit=10"
```
## What You Get Out of the Box
### Database & ORM
- • GORM integration with auto-migrations
- • Connection pooling and optimization
- • Multiple database support (PostgreSQL, MySQL, SQLite)
- • Automatic relationship detection
### Security & Validation
- • Built-in input validation
- • JWT authentication support
- • CORS configuration
- • Rate limiting middleware
### Documentation
- • Automatic Swagger/OpenAPI generation
- • Interactive API documentation
- • Code annotation support
- • Export to JSON/YAML formats
### Developer Experience
- • Hot reloading with Air
- • Structured logging
- • Error handling middleware
- • Environment configuration
## Next Steps
### Blog Tutorial
Follow our comprehensive tutorial to build a complete blog application with authentication, file uploads, and more.
Start Tutorial
### CLI Reference
Explore all CLI commands and learn about advanced field types, relationships, and module generation.
View CLI Docs
### Architecture
Understand Base Framework's HMVC architecture, module system, and design patterns.
Learn Architecture