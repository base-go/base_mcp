---
layout: home

hero:
  name: "Base Framework"
  text: "The Go Framework for Modern APIs"
  tagline: Build production-ready APIs in minutes, not hours
  image:
    src: /logo.png
    alt: Base Framework
  actions:
    - theme: brand
      text: Get Started
      link: /docs/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/base-go/base-core

features:
  - icon: ⚡
    title: Lightning Fast Setup
    details: Get started in seconds with our one-line installer. No complex configuration needed.
  - icon: 🏗️
    title: Code Generation
    details: Create complete CRUD structures with a single command. Models, services, controllers, and tests included.
  - icon: 🔗
    title: Smart Relationships
    details: Automatic relationship detection and GORM integration. Focus on business logic, not boilerplate.
  - icon: 🔐
    title: Built-in Authentication
    details: JWT authentication, role-based authorization, and session management out of the box.
  - icon: 📊
    title: Advanced Logging
    details: Structured logging with multiple outputs, log levels, and performance monitoring.
  - icon: 🚀
    title: Production Ready
    details: Docker support, health checks, graceful shutdowns, and deployment configurations included.
---

## Quick Installation
```
# Install Base CLI
curl -fsSL https://get.base.al | bash

# Create your first project
base new my-awesome-api
cd my-awesome-api

# Start developing immediately
base start # Server running on :8100
```

## Everything You NeedOut of the Box

Base combines the power of Go with intelligent code generation, modern development practices, and production-ready features that scale with your team.

### Automatic Code Generation

Generate complete structures with models, services, controllers, and tests from simple field definitions. No boilerplate, just results.

### Smart Relationship Detection

Fields ending with _id automatically create GORM relationships. Support for belongs_to, has_many, has_one, and many-to-many relationships.

### Modular Architecture

HMVC pattern with dependency injection. Each structure is self-contained with its own controllers, services, and models.

### 15+ Built-in Core Structures

Authentication, authorization, storage, email, websockets, scheduler, translation, and more. Everything you need for production APIs.

### Multi-Cloud Storage

Built-in support for local, AWS S3, and Cloudflare R2 storage with automatic file type validation and event emission.

### High Performance

Zero-dependency HTTP router with radix tree routing, built-in middleware, and optimized for speed and efficiency.

Ready to experience the fastest way to build Go APIs?

[Start Building Now](#get-started)
## Why We Built Base

Go is**incredibly powerful**— fast, concurrent, and built for scale. But unlike other languages, Go lacks full-featured frameworks that make development*effortless*.

While Laravel gives PHP developers magic, Rails gives Ruby developers elegance, and Django gives Python developers batteries-included power, Go developers have been stuck with basic HTTP routers and endless boilerplate.

Base changes that. We bring the power of Go together with the developer experience you deserve.

⚡
### Go's Power

Native performance, concurrency, and type safety

🚀
### Framework Magic

Code generation, auto-relationships, zero config

🛠️
### Powerful CLI

One command to rule them all

### What takes hours in other Go frameworks, takes seconds in Base

```
# Other Go frameworks: Write models, controllers, routes, validation...
# 📁 models/post.go (50+ lines)
# 📁 controllers/post.go (100+ lines)
# 📁 routes/post.go (20+ lines)
# 📁 validation/post.go (30+ lines)
# 📁 migrations/create_posts.sql
# Total: 200+ lines, 5+ files, 30+ minutes

# Base Framework: One command, everything ready
base g Post title:string content:text author_id:uint
# ✅ Complete CRUD API ready in 3 seconds
```

### The CLI that thinks ahead

```
# Smart relationship detection
base g Article title:string author_id:uint category_id:uint

# Base automatically knows:
# author_id → belongs to User model
# category_id → belongs to Category model
# Creates proper GORM relationships
# Sets up foreign key constraints
# Generates preloading queries

// Generated Go code:
type Article struct {
    AuthorId   uint `json:"author_id"`
    Author     User `json:"author" gorm:"foreignKey:AuthorId"`
    CategoryId uint `json:"category_id"`  
    Category   Category `json:"category" gorm:"foreignKey:CategoryId"`
}
```

### Production-ready from day one

```
# Generate a complete e-commerce API
base g Product name:string price:decimal image:attachment
base g Order user_id:uint total:decimal status:string
base g OrderItem order_id:uint product_id:uint quantity:int

# What you get instantly:
# 🔐 JWT Authentication
# 📁 File upload handling
# 🔍 Input validation
# 📊 Database relationships
# 📚 API documentation
# ⚡ High-performance routing
# 🎯 Ready for production
```

< 1minTo first API1Dependency∞PossibilitiescliPowered
## Ready to Build Something Amazing?

Join the growing community of developers revolutionizing Go development with Base Framework

🚀 Use BaseBuild your next Go API in minutes, not hours

🤝 ContributeHelp shape the future of Go development

```
# Install Base CLI
curl -fsSL https://get.base.al | sh

# Create your first project
base new my-awesome-api
cd my-awesome-api

# Generate your first module
base g User name:string email:string

# Start development server
base start -r
```

🚀
### Start Using Base

Build production-ready APIs in minutes

[Star us on GitHub](https://github.com/base-go/base-core)[📚 Read Documentation](/docs)🤝
### Contribute to Base

Help shape the future of Go development

[🐛 Report Issues & Ideas](https://github.com/base-go/base-core/issues)[💻 Submit Pull Requests](https://github.com/base-go/base-core/pulls)[💬 Join Community Discord](https://discord.gg/8wYFX3Wh)