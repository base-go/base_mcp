---
title: Tutorial
description: Tutorial documentation for Base Framework.
---

# Blog Tutorial
Build a complete blog application from scratch and learn Base Framework fundamentals.
## Step 1: Create Your Blog Project
Let's start by creating a new Base Framework project for our blog API:
Terminal
```
base new blog
cd blog
```
#### You'll get initial structure:
```
blog/
├── app/
│ ├── init.go # Module registration
│ └── models/ # Your database models
└── core/ # Core modules
```
## Step 2: Understanding Core Modules
Base Framework comes with built-in core modules that provide essential functionality. Let's explore what's
available
in our blog project:
##### 👤 Authentication Module
- • User model with auth features
- • Login, register, forgot password, reset password, logout endpoints
- • Located at core/app/authentication/
##### 👤 Authorization Module
- • Role model with auth features
- • Permission model with auth features
- • Located at core/app/authorization/
##### 👤 Profile Module
- • User model with auth features
- • Profile management endpoints
- • Avatar upload support
- • Located at core/app/profile/
##### 📁 Media Module
- • File and image management
- • Upload validation
- • Storage integration
- • Located at core/app/media/
## Step 3: Generate Post Module with
Relationships
Now let's generate our blog post module with relationships to the built-in core modules using
belongsTo:profile.User and belongsTo:media.Media:
Terminal
```
`
base g Post \
title:string \
content:text \
slug:string \
published:bool \
author:belongsTo:profile.User \
featured_image:belongsTo:media.Media
`
```
#### Generated Output:
```
Generated app/models/post.go
Generated app/posts/service.go
Generated app/posts/controller.go
Generated app/posts/module.go
Generated app/posts/validator.go
✅ Added module to app/init.go
✅ Dependencies updated
Successfully generated Post module
```
#### What Was Created:
##### 📄 Model (app/models/post.go)
- • Post struct with GORM tags
- • CreatePostRequest & UpdatePostRequest
- • PostResponse & PostListResponse
- • ToResponse() conversion methods
##### 🎮 Controller (app/posts/controller.go)
- • GET /posts (paginated list)
- • POST /posts (create)
- • GET /posts/:id (get by ID)
- • PUT /posts/:id (update)
- • DELETE /posts/:id (delete)
##### ⚙️ Service (app/posts/service.go)
- • Business logic layer
- • Database operations
- • Data validation
- • Error handling
##### 🔧 Module Registration
- • Auto-registered in app/init.go
- • Dependency injection setup
- • Route registration
- • Database migrations
## Step 4: Start Your Blog API Server
Let's start the server with Swagger documentation enabled to see our API in action:
Terminal
```
base start -d
```
#### Server Output:
```
`
📚 Generating swagger documentation...
📋 Found 40 endpoints across 7 controller files
🚀 Starting Base Framework
✅ Database initialized
✅ Infrastructure initialized
✅ Router initialized
✅ Module initialization complete
🎉 Base Framework Ready!
📍 Server URLs:
- Local: http://localhost:8100
- Network: http://192.168.1.35:8100
📚 Documentation:
- Swagger: http://localhost:8100/swagger/index.html
`
```
## Step 5: Test Your Blog API
Your blog API is now ready! Here are the automatically generated endpoints:
POST
/api/posts
Create a new blog post
GET
/api/posts
Get paginated list of posts
GET
/api/posts/all
Get all posts (unpaginated)
GET
/api/posts/:id
Get specific post by ID
PUT
/api/posts/:id
Update existing post
DELETE
/api/posts/:id
Delete post
#### Try It Out - Create Your First User and
Post
First, create a user
```
curl -X POST \
http://localhost:8100/api/auth/register \
-H "Content-Type: application/json" \
-d '{
"first_name": "John",
"last_name": "Doe",
"email": "john@example.com",
"password": "password"
}'
```
Then, create a post with the user relationship
```
curl -X POST \
http://localhost:8100/api/posts \
-H "Content-Type: application/json" \
-d '{
"title": "My First Blog Post",
"content": "My first blog post with Base Framework!",
"slug": "my-first-blog-post",
"published": true,
"author_id": 1
}'
```
## Step 6: Explore What You Built
Interactive API Documentation
Visit http://localhost:8100/swagger/ to explore your API interactively.
- • Try all endpoints directly in the browser
- • See request/response schemas
- • Generate client code
Database Ready
Your database is automatically configured with:
- • Auto-migrations for Post model
- • GORM relationships
- • Optimized indexes
#### What's Next?
##### Complete Blog Example
See the complete implementation below with all
relationships
##### Add Authentication
Enable built-in authentication with JWT tokens
##### Add File Uploads
Use field types like image:image or
attachment:file
for media
## Complete Blog Implementation
Let's build a comprehensive blog system showcasing Base Framework's full capabilities including file
uploads,
relationships, various field types, and event system:
Full-Featured Blog Generation
##### 1. Generate Category Module
Terminal
```
`
base g Category \
name:string \
description:text \
slug:string \
color:string
`
```
Creates Category for organizing posts (one-to-many)
##### 2. Generate Tag Module
Terminal
```
`
base g Tag \
name:string \
slug:string \
color:string \
description:text
`
```
Creates Tag system for labeling posts (many-to-many)
##### 3. Generate Comprehensive Post Module
Terminal
```
`
base g Post \
title:string \
content:text \
slug:string \
excerpt:text \
published:bool \
featured_image:image \
gallery:attachment \
author:belongsTo:profile.User \
category:belongsTo:Category \
tags:toMany:Tag \
publish_date:datetime \
view_count:int
`
```
**Demonstrates multiple relationship types:** belongsTo core modules, belongsTo custom
modules,
plus all field types
##### 4. Generate Advanced Comment System
Terminal
```
`
base g Comment \
content:text \
post:belongsTo:Post \
author:belongsTo:profile.User \
parent_comment:belongsTo:Comment \
reply_to_email:email \
approved:bool \
ip_address:string \
user_agent:text
`
```
**Advanced features:** nested comments, email validation, moderation, tracking fields
##### 5. Start with Full Documentation
Terminal
```
base s -d
```
Includes automatic file upload endpoints, validation, and Swagger docs
Field Types Demonstrated
##### Basic Types
string
Short text fields
text
Long text content
bool
True/false values
int
Numeric values
datetime
Date and time
##### Advanced Types
email
Email validation
image
Image uploads
attachment
File uploads
belongsTo
Relationships
Generated API Endpoints (40+ endpoints)
##### Categories API
CRUD
/api/categories
##### Tags API
CRUD
/api/tags
##### Posts API + File Uploads
CRUD
/api/posts/*
POST
/api/posts/:id/featured_image
DELETE
/api/posts/:id/featured_image
##### Comments + Core
CRUD
/api/comments/*
GET
/api/profile/*
POST
/api/auth/*
**✨ Automatic Features:** 50+ endpoints with pagination, sorting, filtering, validation,
file
upload endpoints, relationship preloading, Swagger docs, and event emission!
File Upload Examples
##### Upload Featured Image to Post
```
`
curl -X POST http://localhost:8100/api/posts/1/featured_image \
-H "Content-Type: multipart/form-data" \
-F "file=@hero-image.jpg"
`
```
Automatic validation, resizing, storage, and database association
##### Create Post with Rich Content
```
`
curl -X POST http://localhost:8100/api/posts \
-H "Content-Type: application/json" \
-d '{
"title": "Advanced Blog Post",
"content": "Full HTML content...",
"excerpt": "Post summary",
"slug": "advanced-blog-post",
"published": true,
"tags": "tech,base-framework,tutorial",
"publish_date": "2024-12-01T10:00:00Z",
"view_count": 0,
"author_id": 1,
"category_id": 1
}'
`
```
Base Framework Features Demonstrated
##### 🔧 Code Generation
- • Models with GORM tags
- • Controllers with validation
- • Services with business logic
- • Module registration
##### 📁 File Management
- • Image upload & validation
- • File attachment handling
- • Storage abstraction
- • Automatic cleanup
##### 🔗 Relationship Types
- • belongsTo:profile.User (core)
- • belongsTo:Category (custom)
- • belongsTo:Tag (labeling)
- • belongsTo:Comment (self-ref)
##### ✅ Validation
- • Email format validation
- • Required field checks
- • File type validation
- • Custom validation rules
##### 🎯 Event System
- • Post creation events
- • File upload events
- • Comment moderation
- • Email notifications
##### 📚 Documentation
- • Swagger/OpenAPI specs
- • Interactive API explorer
- • Request/response schemas
- • Authentication docs