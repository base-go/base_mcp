---
layout: home

hero:
  name: Base Framework
  text: The Go Framework for Modern APIs
  tagline: Simple, powerful, and performant. Meet the modern Go framework you've been waiting for.
  actions:
    - theme: brand
      text: Get Started
      link: /docs/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/base-go/base-core
---

<FeatureGrid 
  title="Everything You Need Out of the Box"
  description="Base combines the power of Go with intelligent code generation, modern development practices, and production-ready features that scale with your team."
  :features="[
    {
      icon: '⚡',
      title: 'Automatic Code Generation',
      details: 'Generate complete structures with models, services, controllers, and tests from simple field definitions. No boilerplate, just results.'
    },
    {
      icon: '🔗',
      title: 'Smart Relationship Detection',
      details: 'Fields ending with _id automatically create GORM relationships. Support for belongs_to, has_many, has_one, and many-to-many relationships.'
    },
    {
      icon: '🏗️',
      title: 'Modular Architecture',
      details: 'HMVC pattern with dependency injection. Each structure is self-contained with its own controllers, services, and models.'
    },
    {
      icon: '📦',
      title: '15+ Built-in Core Structures',
      details: 'Authentication, authorization, storage, email, websockets, scheduler, translation, and more. Everything you need for production APIs.'
    },
    {
      icon: '☁️',
      title: 'Multi-Cloud Storage',
      details: 'Built-in support for local, AWS S3, and Cloudflare R2 storage with automatic file type validation and event emission.'
    },
    {
      icon: '🚀',
      title: 'High Performance',
      details: 'Zero-dependency HTTP router with radix tree routing, built-in middleware, and optimized for speed and efficiency.'
    }
  ]"
/>

## Framework as Code

Base Framework generates complete API structures from simple field definitions, making it perfect for rapid prototyping and production applications.

```bash
# Create your first project
base new my-awesome-api
cd my-awesome-api

# Generate complete CRUD structure
base g User name:string email:string age:int
```

Each command generates models, controllers, services, validation, tests, and API routes. This allows you to focus on business logic while Base handles the infrastructure.

### What takes hours in other Go frameworks, takes seconds in Base

```bash
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

### Smart relationship detection

```go
// Smart relationship detection
base g Article title:string author_id:uint category_id:uint

// Base automatically knows:
// author_id → belongs to User model
// category_id → belongs to Category model
// Creates proper GORM relationships
// Sets up foreign key constraints
// Generates preloading queries

type Article struct {
    AuthorId   uint `json:"author_id"`
    Author     User `json:"author" gorm:"foreignKey:AuthorId"`
    CategoryId uint `json:"category_id"`  
    Category   Category `json:"category" gorm:"foreignKey:CategoryId"`
}
```

## Why We Built Base

### Go is **incredibly powerful** — fast, concurrent, and built for scale. But unlike other languages, Go lacks full-featured frameworks that make development *effortless*.

While Laravel gives PHP developers magic, Rails gives Ruby developers elegance, and Django gives Python developers batteries-included power, Go developers have been stuck with basic HTTP routers and endless boilerplate.

Base changes that. We bring the power of Go together with the developer experience you deserve.

<div class="vp-feature-grid">
  <div class="vp-feature">
    <div class="vp-feature-icon">⚡</div>
    <h3>Go's Power</h3>
    <p>Native performance, concurrency, and type safety</p>
  </div>
  
  <div class="vp-feature">
    <div class="vp-feature-icon">🚀</div>
    <h3>Framework Magic</h3>
    <p>Code generation, auto-relationships, zero config</p>
  </div>
  
  <div class="vp-feature">
    <div class="vp-feature-icon">🛠️</div>
    <h3>Powerful CLI</h3>
    <p>One command to rule them all</p>
  </div>
</div>
 

## Join the Community

Ready to revolutionize your Go development experience?

<div class="vp-feature-grid">
  <div class="vp-feature">
    <div class="vp-feature-icon">⭐</div>
    <h3>Star us on GitHub</h3>
    <p>Show your support and stay updated with the latest releases.</p>
    <a href="https://github.com/base-go/base-core" class="vp-feature-link">GitHub →</a>
  </div>
  
  <div class="vp-feature">
    <div class="vp-feature-icon">💬</div>
    <h3>Join Discord</h3>
    <p>Connect with other developers and get help from the community.</p>
    <a href="https://discord.gg/8wYFX3Wh" class="vp-feature-link">Discord →</a>
  </div>
  
  <div class="vp-feature">
    <div class="vp-feature-icon">🐛</div>
    <h3>Report Issues</h3>
    <p>Help us improve by reporting bugs and suggesting new features.</p>
    <a href="https://github.com/base-go/base-core/issues" class="vp-feature-link">Issues →</a>
  </div>
</div>