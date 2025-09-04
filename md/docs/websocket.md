---
title: Websocket
description: Websocket documentation for Base Framework.
---

# WebSocket
Real-time communication with WebSocket support. Built-in hub management, room-based messaging, and collaborative features.
## Overview
### Multi-Room Support
Organize users into rooms with automatic user list management and room-specific broadcasting.
### Real-time Messaging
Instant message delivery with system notifications and typing indicators.
### Collaborative Features
Cursor tracking, collaborative drawing, code editing, and shared whiteboards.
Base Framework includes a production-ready WebSocket implementation with automatic hub initialization, connection management, and message broadcasting. The system supports multiple rooms, user management, and various collaborative features out of the box.
WebSocket Endpoint
The WebSocket server runs automatically when you start Base Framework and is available at `/api/ws`
## Connecting to WebSocket
### Basic Connection
```javascript
// Connect to WebSocket server
const socket = new WebSocket('ws://localhost:8100/api/ws?id=user123&nickname=John&room=general');
socket.onopen = function(event) {
console.log('Connected to WebSocket server');
};
socket.onmessage = function(event) {
const message = JSON.parse(event.data);
console.log('Received message:', message);
};
socket.onclose = function(event) {
console.log('Disconnected from WebSocket server');
};
socket.onerror = function(error) {
console.error('WebSocket error:', error);
};
```
### Connection Parameters
Parameter
Required
Description
id
Optional
Unique client identifier
nickname
Optional
Display name for the user
room
Optional
Room to join (default: general)
## Message Types
### Chat Messages
```javascript
// Send a chat message
socket.send(JSON.stringify({
type: 'message',
content: 'Hello, everyone!',
room: 'general',
nickname: 'John'
}));
// Receive chat message
{
"type": "message",
"content": "Hello, everyone!",
"room": "general",
"nickname": "John"
}
```
### System Messages
```javascript
// System notifications (received only)
{
"type": "system",
"content": "John joined the room",
"room": "general",
"nickname": "System"
}
// User list updates
{
"type": "users_update",
"content": ["John", "Jane", "Bob"],
"room": "general"
}
```
### Collaborative Features
```javascript
// Drawing data
socket.send(JSON.stringify({
type: 'draw',
content: {
fromX: 100, fromY: 150,
toX: 200, toY: 250,
color: '#FF0000'
}
}));
// Cursor movement
socket.send(JSON.stringify({
type: 'cursor_move',
content: { x: 150, y: 200 }
}));
// Code updates
socket.send(JSON.stringify({
type: 'code_update',
content: {
code: 'console.log("Hello World");',
language: 'javascript'
}
}));
```
## Live Examples
Base Framework includes several built-in WebSocket examples that you can access when running your server. These examples demonstrate different real-time features and can be used as templates for your own applications.
### Chat Application
Multi-room chat with user lists, system notifications, and real-time messaging.
`/static/chat.html`
View Example
### Collaborative Drawing
Real-time drawing canvas with cursor tracking and collaborative sketching.
`/static/draw.html`
View Example
### Kanban Board
Collaborative task management with drag-and-drop and real-time updates.
`/static/kanban.html`
View Example
### Code Editor
Collaborative code editing with syntax highlighting and real-time changes.
`/static/editor.html`
View Example
### System Monitor
Real-time system metrics dashboard with live charts and monitoring.
`/static/monitor.html`
View Example
### Collaborative Spreadsheet
Real-time spreadsheet editing with cell updates and collaborative features.
`/static/spreadsheet.html`
View Example
Access Examples
When your Base Framework server is running, visit `http://localhost:8100/static/` to see all available examples. Each example includes full source code you can adapt for your needs.
## Integration with Base Framework
### Configuration
WebSocket support can be enabled or disabled using the `WS_ENABLED` environment variable. By default, WebSocket is enabled.
```bash
# Enable WebSocket support (default)
WS_ENABLED=true
# Disable WebSocket support
WS_ENABLED=false
```
Conditional Initialization
When `WS_ENABLED=false`, the WebSocket hub is not initialized and the `/api/ws` endpoint is not registered, reducing memory usage and startup time.
### Automatic Initialization
When WebSocket is enabled, it's automatically initialized during application startup. The WebSocket hub starts automatically and registers the `/api/ws` endpoint.
### Broadcasting from Controllers
```go
// Access WebSocket hub in your controllers
func (c *YourController) BroadcastUpdate(ctx *router.Context) error {
// Get the WebSocket hub (this would need to be injected)
// hub.BroadcastMessage("update", map[string]any{
//     "type": "data_update",
//     "payload": yourData,
// })
return ctx.JSON(200, map[string]any{"status": "broadcasted"})
}
```
### Event Integration
WebSocket integrates seamlessly with Base Framework's event system. You can listen for model changes and automatically broadcast updates to connected clients.
## Security & Best Practices
CORS Configuration
The WebSocket upgrader currently allows all origins for development. In production, configure the CheckOrigin function to validate allowed origins.
### Authentication
WebSocket connections can be authenticated using query parameters, headers, or by validating tokens during the connection upgrade process.
### Rate Limiting
Implement rate limiting for WebSocket messages to prevent spam and abuse. Consider limiting messages per user per second.
### Message Validation
Always validate incoming WebSocket messages on the server side. Never trust client-provided data without proper validation and sanitization.
## Architecture
### Hub Pattern
Base Framework uses a centralized hub pattern for WebSocket management. The hub maintains all active connections, handles room management, and coordinates message broadcasting between clients.
### Connection Lifecycle
- **Connect:** Client upgrades HTTP connection to WebSocket
- **Register:** Client is added to the hub and assigned to a room
- **Message Handling:** Incoming messages are processed and broadcast
- **Disconnect:** Client is removed from hub and room notifications sent
### Scalability
For high-scale applications, consider using Redis or other message brokers to enable WebSocket clustering across multiple server instances.