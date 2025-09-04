---
title: Auth
description: Auth documentation for Base Framework.
---

# Authentication & Authorization
Understand Base Framework's comprehensive authentication system with built-in OAuth, JWT tokens, and role-based access control.
## Core Authentication Modules
Base Framework comes with three built-in core modules that work together to provide complete authentication and authorization:
Authentication
Located at core/app/authentication/
- • JWT token generation
- • User login/logout
- • Password hashing
- • Token validation
- • Session management
OAuth
Located at core/app/oauth/
- • Google OAuth integration
- • GitHub OAuth support
- • Facebook login
- • Custom OAuth providers
- • Social account linking
Authorization
Located at core/app/authorization/
- • Role-based access control
- • Permission management
- • Route protection
- • API key authentication
- • Middleware integration
## Authentication Flow
JWT Authentication Process
1
##### User Registration/Login
POST to /api/auth/login or /api/auth/register
2
##### Credential Validation
Password verification using bcrypt hashing
3
##### JWT Token Generation
Signed JWT with user ID, extended data, and expiration
4
##### Token Storage
Client stores token in localStorage or httpOnly cookie
API Authentication Example
##### Login Request
cURL Example
```
curl -X POST http://localhost:8100/api/auth/login \\
-H "Content-Type: application/json" \\
-d '{
"email": "user@example.com",
"password": "securepassword"
}'
```
##### Response with Role Information
```json
{
"id": 1,
"first_name": "John",
"last_name": "Doe",
"username": "johndoe",
"phone": "+1234567890",
"email": "john@example.com",
"role_id": 1,
"role_name": "Owner",
"avatar_url": "",
"last_login": "2025-09-02T15:44:05+02:00",
"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
"exp": 1756908308,
"extend": {
"role": {
"id": 1,
"name": "Owner"
},
"user_id": 1
}
}
```
##### Authenticated API Call
Using JWT Token
```
curl -X GET http://localhost:8100/api/posts \\
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
## OAuth Integration
Supported OAuth Providers
##### Google OAuth
GET /api/oauth/google
##### GitHub OAuth
GET /api/oauth/github
##### Facebook OAuth
GET /api/oauth/facebook
##### Custom Provider
Easily add new providers
OAuth Configuration
Configure OAuth providers in your .env file:
.env Configuration
```
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8100/api/oauth/google/callback
# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=http://localhost:8100/api/oauth/github/callback
```
## Authorization & Permissions
First User Owner System
Base Framework automatically assigns the **Owner** role to the first user who registers, ensuring secure system initialization without manual intervention:
##### 👑 Automatic Role Assignment
**First User**
- • Automatically receives **Owner** role
- • Full system access and control
- • Can manage all users and permissions
- • No manual role assignment needed
**Subsequent Users**
- • Automatically receive **Member** role
- • Standard limited access
- • Can be promoted by Owner/Administrator
- • Secure default permissions
Role Hierarchy
Base Framework uses a hierarchical role system with four default roles, each with specific permissions and access levels:
##### Owner (Role ID: 1)
Full system access • Manage all users • System configuration
First User
##### Administrator (Role ID: 2)
Administrative access • User management • Most operations
##### Member (Role ID: 3)
Standard access • Limited permissions • Default for new users
Default
##### Viewer (Role ID: 4)
Read-only access • View permissions only
Role-Based Access Control (RBAC)
Base Framework includes a flexible RBAC system that allows you to define roles and permissions for your application:
##### 👤 Users
Assigned one or more roles
##### 🎭 Roles
Collections of permissions
##### 🔐 Permissions
Specific access rights
Route Protection
Protect your API endpoints using middleware decorators:
Go Code Example
```
// Require authentication (basic auth middleware)
router.GET("/api/protected", controller.ProtectedEndpoint, authMiddleware)
// Require specific role
router.POST("/api/admin", controller.AdminEndpoint, authorization.HasRole("Administrator"))
// Require specific permission
router.DELETE("/api/posts/:id", controller.DeletePost, authorization.CanAccess("delete", "Post", "id"))
// Multiple middleware combined
router.PUT("/api/posts/:id", controller.UpdatePost,
authMiddleware,
authorization.CanAccess("update", "Post", "id"))
```
API Key Authentication
Base Framework supports multiple authentication methods. **Bearer tokens (JWT) are the default for all generated modules**, while API keys provide additional flexibility for service-to-service communication:
##### 🔑 Authentication Methods in Base Framework
**Bearer Token (JWT) - Default**
- • Header: `Authorization: Bearer token`
- • **Default auth for generated modules**
- • User sessions, temporary access
- • Expires automatically
- • Contains user context
**API Key - Environment Variable**
- • Header: `X-API-Key: key`
- • **Set via environment variable**
- • Required for auth management operations
- • Service-to-service communication
- • Application-level access control
##### Configure API Key
.env Configuration
```
# JWT secret key
JWT_SECRET=my_super_secret_key_123
# Send this key in the header of the request for endpoints that require it
API_KEY=api
```
##### Use API Key in Requests
X-API-Key Header
```
curl -X GET http://localhost:8100/api/posts \\
-H "X-API-Key: api"
```
## Extended Authentication
Base Framework supports extending JWT tokens with custom data through the app.Extend() function. This allows you to include additional context like company IDs, roles, or any other data in both the authentication response and JWT token payload.
app.Extend() Function
Located at app/init.go, this function is called during authentication to extend the user context with additional data. **You must customize this function in your app to add extended data.**
app/init.go
```
func Extend(user_id uint) any {
// Get database instance
if database.DB == nil {
return map[string]any{
"user_id": user_id,
}
}
// Example: Add company_id to the context
return map[string]any{
"user_id":    user_id,
"company_id": 1,
}
}
```
**⚠️ Important:** You must modify the app.Extend() function in base/app/init.go to add your custom extended data. The default implementation only returns the user_id.
Enhanced Login Response
When app.Extend() returns data, it's included in both the login response and embedded in the JWT token payload.
##### Login Request
POST /api/auth/login
```
curl -X 'POST' \\
'http://localhost:8100/api/auth/login' \\
-H 'accept: application/json' \\
-H 'X-Api-Key: api' \\
-H 'Content-Type: application/json' \\
-d '{
"email": "john@example.com",
"password": "password123"
}'
```
##### Enhanced Response
Response with Extended Data
```
{
"id": 1,
"first_name": "John",
"last_name": "Doe",
"username": "johndoe",
"phone": "+1234567890",
"email": "john@example.com",
"avatar_url": "",
"last_login": "2025-09-02T15:42:51+02:00",
"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
"exp": 1756907010,
"extend": {
"company_id": 1,
"user_id": 1
}
}
```
JWT Token Structure
The JWT token payload includes the extended data, making it available for authorization checks and middleware without additional database queries.
##### JWT Payload Structure
Decoded JWT Payload
```
{
"exp": 1756907010,
"extend": {
"company_id": 1,
"user_id": 1
},
"user_id": 1
}
```
##### Common Use Cases
###### Multi-tenant Applications
Include company_id or tenant_id for data isolation
###### Role-based Access
Embed user roles and permissions in the token
###### Context Data
Include any additional context your app needs
##### 🚀 Benefits
- • **Reduced Database Queries:** Context data is embedded in the JWT token
- • **Consistent Data:** Same extended data available in both response and token
- • **Flexible Implementation:** Customize the Extend function for your specific needs
- • **Middleware Ready:** Extended data is available in authorization middleware
## Security Features
Password Security
-
bcrypt hashing with salt rounds
-
Password strength validation
-
Password reset flows
-
Account lockout protection
Token Management
-
JWT with configurable expiration
-
Refresh token rotation
-
Token blacklisting
-
Session management
Rate Limiting
-
IP-based rate limiting
-
User-based rate limiting
-
Endpoint-specific limits
-
Sliding window algorithm
Audit & Monitoring
-
Login attempt logging
-
Permission change tracking
-
Failed authentication alerts
-
Security event notifications
## Context Helpers &amp; Retrieving the Authenticated User
Using the Request Context
After successful authentication, the middleware saves the user object in the request `context.Context` using a type-safe key. Retrieve it anywhere deeper in your application with the helper below:
Go Example
```
// Inside a handler or service
if user, ok := middleware.UserFromContext[*models.User](c.Request.Context()); ok {
// user is *models.User
}
```
The same approach works for API-key data (`api_key_data`) or any other custom values you attach to the context.
## Attaching Custom Values to Context
Adding Your Own Context Keys
Need to pass request-scoped data (e.g. `project_id`) deeper into your services? Define a unique pointer key and use `context.WithValue` just like the built-in auth middleware:
Go Example
```
type ctxKey struct{ name string }
var projectIDKey = &amp;ctxKey{"project_id"}
// Middleware that stores project_id from route params
func ProjectIDMiddleware(next router.HandlerFunc) router.HandlerFunc {
return func(c *router.Context) error {
id := c.Param("project_id")
ctx := context.WithValue(c.Request.Context(), projectIDKey, id)
c.Request = c.Request.WithContext(ctx)
return next(c)
}
}
// Helper to retrieve it later
func ProjectIDFromContext(ctx context.Context) (string, bool) {
v, ok := ctx.Value(projectIDKey).(string)
return v, ok
}
```
Remember: context values should be used *sparingly* for data that needs to travel across API boundaries. For regular function parameters prefer explicit arguments.
## Authorization System
Authorization Models
The authorization system consists of four main data models that work together:
#####
Roles
System and custom roles with hierarchical permissions
• **Owner** - Full system access
• **Administrator** - Admin privileges
• **Member** - Standard user access
• **External** - Limited external access
#####
Permissions
Specific action rights on resource types
• **Actions:** create, read, update, delete, list
• **Resources:** posts, users, projects, etc.
#####
Role Permissions
Links roles to their allowed permissions
Many-to-many relationship mapping roles to specific permissions
#####
Resource Permissions
Fine-grained access to specific resources
• **Scopes:** own, team, all
• Resource-specific overrides
Authorization Middleware in Controllers
Base Framework provides intuitive authorization middleware functions with clean, readable syntax:
#####
✨ New Intuitive Syntax
Use the new Can() functions for clear, readable authorization checks:
Clean Authorization Syntax
```
import "authorization"
func (c *PostController) Routes(router *router.RouterGroup) {
// ✨ Clean, readable authorization
router.POST("/posts", c.Create, authorization.Can("create", "Post"))
router.GET("/posts", c.List, authorization.Can("list", "Post"))
router.PUT("/posts/:id", c.Update, authorization.CanAccess("update", "Post", "id"))
router.DELETE("/posts/:id", c.Delete, authorization.CanAccess("delete", "Post", "id"))
// Role-based protection
adminGroup := router.Group("/admin", authorization.HasRole("Administrator"))
ownerGroup := router.Group("/owner", authorization.HasRole("Owner"))
}
```
##### 1. Can() - General Resource Protection
The intuitive way to check if a user has permission to perform an action on a resource type.
###### ✨ New Syntax (Recommended)
Can() Function
```
// ✨ Clean, readable syntax
router.POST("/posts", c.Create, authorization.Can("create", "Post"))
router.GET("/posts", c.List, authorization.Can("list", "Post"))
router.PUT("/posts", c.Update, authorization.Can("update", "Post"))
// Multiple permissions
router.POST("/posts", c.Create,
authorization.Can("create", "Post"),
authorization.HasRole("Member"))
```
###### 🔧 Legacy Syntax
AuthMiddleware
```
// Legacy syntax (still supported)
router.POST("/posts", c.Create,
authorization.AuthMiddleware("post", "create"))
router.GET("/posts", c.List,
authorization.AuthMiddleware("post", "list"))
```
##### 2. CanAccess() - Specific Resource Protection
Checks if a user has permission to perform an action on a specific resource instance.
###### ✨ New Syntax (Recommended)
CanAccess() Function
```
// ✨ Check access to specific resources
router.PUT("/posts/:id", c.Update,
authorization.CanAccess("update", "Post", "id"))
router.DELETE("/posts/:id", c.Delete,
authorization.CanAccess("delete", "Post", "id"))
// Different parameter names
router.PUT("/users/:userId/posts/:postId", c.UpdateUserPost,
authorization.CanAccess("update", "Post", "postId"))
```
###### 🔧 Legacy Syntax
ResourceAuthMiddleware
```
// Legacy syntax (still supported)
router.PUT("/posts/:id", c.Update,
authorization.ResourceAuthMiddleware("post", "update", "id"))
router.DELETE("/posts/:id", c.Delete,
authorization.ResourceAuthMiddleware("post", "delete", "id"))
```
##### 3. HasRole() - Role-Based Protection
Ensures the user has a specific role (Owner, Administrator, Member, External).
###### ✨ New Syntax (Recommended)
HasRole() Function
```
// ✨ Clean role checking
adminGroup := router.Group("/admin",
authorization.HasRole("Administrator"))
ownerGroup := router.Group("/owner",
authorization.HasRole("Owner"))
// Individual routes
router.GET("/admin/settings", c.GetSettings,
authorization.HasRole("Administrator"))
// Multiple role checks
router.POST("/posts", c.Create,
authorization.HasRole("Member"),
authorization.Can("create", "Post"))
```
###### 🔧 Legacy Syntax
RequireRole
```
// Legacy syntax (still supported)
adminGroup := router.Group("/admin",
authorization.RequireRole("Administrator"))
ownerGroup := router.Group("/owner",
authorization.RequireRole("Owner"))
```
##### 4. Advanced Permission Checks
Additional flexible authorization functions for complex scenarios.
###### 🔄 CanAny() - Any Permission
OR Logic
```
// User needs ANY of these permissions
router.GET("/posts/moderate", c.ModeratePosts,
authorization.CanAny([]string{
"moderate:Post",
"delete:Post",
"update:Post",
}))
// Editor or admin access
router.POST("/content", c.CreateContent,
authorization.CanAny([]string{
"create:Content",
"edit:Content",
}))
```
###### 🔒 CanAll() - All Permissions
AND Logic
```
// User needs ALL of these permissions
router.POST("/posts/publish", c.PublishPost,
authorization.CanAll([]string{
"create:Post",
"publish:Post",
"read:Post",
}))
// Requires multiple related permissions
router.DELETE("/users/:id", c.DeleteUser,
authorization.CanAll([]string{
"delete:User",
"manage:User",
}))
```
Authorization API Endpoints
The authorization module exposes REST API endpoints for managing roles, permissions, and checking access:
##### Role Management
GET /authorization/roles
List
GET /authorization/roles/:id
Get
POST /authorization/roles
Create
PUT /authorization/roles/:id
Update
DELETE /authorization/roles/:id
Delete
##### Permission Management
GET /authorization/roles/:id/permissions
List
POST /authorization/roles/:id/permissions
Assign
DELETE /authorization/roles/:id/permissions/:permId
Revoke
POST /authorization/check
Check
Frontend Permission Checking
Check user permissions from your frontend using the authorization API endpoints:
##### Vue.js Composable Example
useAuth.js Composable
```
// composables/useAuth.js
import { ref, computed } from 'vue'
const user = ref(null)
const token = ref(localStorage.getItem('auth_token'))
export function useAuth() {
const isAuthenticated = computed(() => !!token.value && !!user.value)
// Check if user has permission for a resource type and action
const checkPermission = async (resourceType, action, resourceId = null) => {
if (!isAuthenticated.value) return false
try {
const response = await fetch('/api/authorization/check', {
method: 'POST',
headers: {
'Content-Type': 'application/json',
'Authorization': `Bearer ${token.value}`,
'Base-Orgid': getCurrentOrgId() // Your org ID
},
body: JSON.stringify({
user_id: user.value.id,
organization_id: getCurrentOrgId(),
resource_type: resourceType,
action: action,
resource_id: resourceId
})
})
const data = await response.json()
return data.has_permission || false
} catch (error) {
console.error('Permission check failed:', error)
return false
}
}
// Check if user has a specific role
const hasRole = (roleName) => {
return user.value?.roles?.some(role => role.name === roleName) || false
}
return {
user,
token,
isAuthenticated,
checkPermission,
hasRole
}
}
```
##### Vue Component Usage
Component Example
```
&lt;template&gt;
&lt;div class="post-actions"&gt;
&lt;!-- Show edit button only if user can update posts --&gt;
&lt;button
v-if="canEditPost"
@click="editPost"
class="btn btn-primary"&gt;
Edit Post
&lt;/button&gt;
&lt;!-- Show delete button only if user can delete this specific post --&gt;
&lt;button
v-if="canDeletePost"
@click="deletePost"
class="btn btn-danger"&gt;
Delete Post
&lt;/button&gt;
&lt;!-- Show admin panel only for administrators --&gt;
&lt;router-link
v-if="isAdmin"
to="/admin"
class="btn btn-secondary"&gt;
Admin Panel
&lt;/router-link&gt;
&lt;/div&gt;
&lt;/template&gt;
&lt;script setup&gt;
import { ref, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
const props = defineProps(['post'])
const { checkPermission, hasRole } = useAuth()
const canEditPost = ref(false)
const canDeletePost = ref(false)
const isAdmin = ref(false)
onMounted(async () => {
// Check general permission to update posts
canEditPost.value = await checkPermission('post', 'update')
// Check permission to delete this specific post
canDeletePost.value = await checkPermission('post', 'delete', props.post.id)
// Check if user has Administrator role
isAdmin.value = hasRole('Administrator')
})
&lt;/script&gt;
```
##### JavaScript/React Example
React Hook
```
// hooks/usePermissions.js
import { useState, useEffect } from 'react'
import { useAuth } from './useAuth'
export function usePermissions() {
const { user, token, isAuthenticated } = useAuth()
const checkPermission = async (resourceType, action, resourceId = null) => {
if (!isAuthenticated) return false
try {
const response = await fetch('/api/authorization/check', {
method: 'POST',
headers: {
'Content-Type': 'application/json',
'Authorization': `Bearer ${token}`,
'Base-Orgid': getCurrentOrgId()
},
body: JSON.stringify({
user_id: user.id,
organization_id: getCurrentOrgId(),
resource_type: resourceType,
action: action,
resource_id: resourceId
})
})
const data = await response.json()
return data.has_permission || false
} catch (error) {
console.error('Permission check failed:', error)
return false
}
}
return { checkPermission }
}
// Component usage
function PostComponent({ post }) {
const [canEdit, setCanEdit] = useState(false)
const [canDelete, setCanDelete] = useState(false)
const { checkPermission } = usePermissions()
useEffect(() => {
const checkPermissions = async () => {
setCanEdit(await checkPermission('post', 'update'))
setCanDelete(await checkPermission('post', 'delete', post.id))
}
checkPermissions()
}, [post.id])
return (
&lt;div className="post-actions"&gt;
{canEdit && &lt;button onClick={editPost}&gt;Edit&lt;/button&gt;}
{canDelete && &lt;button onClick={deletePost}&gt;Delete&lt;/button&gt;}
&lt;/div&gt;
)
}
```
## Integration Examples
Protecting Generated Modules
When you generate modules with Base Framework, you can easily add authentication to protect your endpoints:
Generated Controller with Auth
```
func (c *PostController) Routes(router *router.RouterGroup) {
// Public endpoints (no auth required)
router.GET("/posts", c.List)                    // Public post listing
router.GET("/posts/:id", c.Get)               // Public post view
// Protected endpoints - using new intuitive syntax
router.POST("/posts", c.Create, authorization.Can("create", "Post"))
router.PUT("/posts/:id", c.Update, authorization.CanAccess("update", "Post", "id"))
router.DELETE("/posts/:id", c.Delete, authorization.CanAccess("delete", "Post", "id"))
// Admin endpoints (role-based protection)
adminGroup := router.Group("/admin/posts", authorization.HasRole("Administrator"))
{
adminGroup.DELETE("/:id", c.ForceDelete)      // Admin force delete
adminGroup.PUT("/:id/featured", c.SetFeatured) // Admin feature post
adminGroup.POST("/bulk-delete", c.BulkDelete)  // Admin bulk operations
}
// Mixed protection: role + permission
router.POST("/posts/publish", c.PublishPost,
authorization.HasRole("Member"),
authorization.Can("publish", "Post"))
}
```