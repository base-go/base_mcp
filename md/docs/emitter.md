---
title: Emitter
description: Emitter documentation for Base Framework.
---

# Emitter
Thread-safe async event system with context support for decoupled module communication
## Overview
🔗 Event-Driven Architecture
Base Framework's Emitter provides a thread-safe, async event system that enables decoupled communication between structures. It supports synchronous and asynchronous event emission, context cancellation, timeouts, and automatic panic recovery for robust event handling.
Thread-Safe
Concurrent-safe operations with read-write mutexes ensuring data integrity across goroutines.
Async Events
Support for both blocking and non-blocking event emission with goroutine-based execution.
Context Support
Built-in context cancellation and timeout support for better control over event processing.
## Basic Usage
Creating and Using an Emitter
Emitter Initialization and Basic Events
```
import (
"base/core/emitter"
"base/core/logger"
)
// Initialize emitter
e := emitter.New()
// Register event listeners
e.On("user.created", func(data any) {
if user, ok := data.(*User); ok {
log.Info("New user registered",
logger.String("email", user.Email),
logger.Int("id", int(user.ID)))
}
})
e.On("user.created", func(data any) {
// Multiple listeners can handle the same event
if user, ok := data.(*User); ok {
// Send welcome email, update analytics, etc.
sendWelcomeEmail(user)
}
})
// Emit events
user := &User{ID: 1, Email: "john@example.com"}
e.Emit("user.created", user) // Blocks until all listeners complete
// Service integration example
type UserService struct {
DB      *gorm.DB
Emitter *emitter.Emitter
Logger  logger.Logger
}
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
user := &User{
Email: req.Email,
Name:  req.Name,
}
if err := s.DB.Create(user).Error; err != nil {
return nil, err
}
// Emit event after successful creation
s.Emitter.Emit("user.created", user)
return user, nil
}
```
## API Reference
Core Methods
#### New() *Emitter
Creates a new emitter instance with initialized listener map.
e := emitter.New()
#### On(event string, listener func(any))
Registers an event listener for the specified event name. Multiple listeners can be registered for the same event.
e.On("user.login", func(data any) { ... })
#### Emit(event string, data any)
Emits an event synchronously, blocking until all listeners complete. Uses WaitGroup internally to ensure all goroutines finish.
e.Emit("order.completed", order)
#### EmitAsync(event string, data any)
Emits an event asynchronously without blocking. Listeners run in separate goroutines and the method returns immediately.
e.EmitAsync("analytics.track", eventData)
#### EmitWithContext(ctx context.Context, event string, data any) error
Emits an event with context support. Returns an error if the context is cancelled before all listeners complete.
err := e.EmitWithContext(ctx, "file.process", fileData)
#### EmitWithTimeout(event string, data any, timeout time.Duration) error
Emits an event with a timeout. Returns an error if listeners don't complete within the specified duration.
err := e.EmitWithTimeout("heavy.task", data, 30*time.Second)
#### ListenerCount(event string) int
Returns the number of listeners registered for a specific event.
count := e.ListenerCount("user.created")
#### EventNames() []string
Returns all registered event names as a slice of strings.
events := e.EventNames()
#### Clear()
Removes all registered listeners for all events. Useful for testing or cleanup.
e.Clear()
## Async Events & Context Support
Asynchronous Event Patterns
Non-blocking Event Emission
```
// Non-blocking events for analytics, logging, notifications
e.EmitAsync("analytics.page_view", analyticsData)
e.EmitAsync("audit.log", auditEvent)
e.EmitAsync("notification.send", notificationData)
// Async event handlers
e.On("file.uploaded", func(data any) {
if file, ok := data.(*UploadedFile); ok {
// Heavy processing that shouldn't block the request
go processImageThumbnails(file)
go scanForVirus(file)
go updateSearchIndex(file)
}
})
// Fire-and-forget pattern for non-critical operations
func (s *PostService) PublishPost(post *Post) error {
if err := s.DB.Save(post).Error; err != nil {
return err
}
// These operations shouldn't block the HTTP response
s.Emitter.EmitAsync("post.published", post)
s.Emitter.EmitAsync("analytics.content_created", post)
s.Emitter.EmitAsync("social.auto_share", post)
return nil
}
// Mixed patterns: critical vs non-critical events
func (s *OrderService) ProcessOrder(order *Order) error {
// Critical event - wait for completion
s.Emitter.Emit("order.validated", order)
if err := s.DB.Save(order).Error; err != nil {
return err
}
// Critical notifications
s.Emitter.Emit("order.confirmed", order)
// Non-critical analytics and marketing
s.Emitter.EmitAsync("analytics.order", order)
s.Emitter.EmitAsync("marketing.customer_activity", order)
return nil
}
```
Context and Timeout Handling
Context-Aware Event Processing
```
import (
"context"
"time"
)
// Context cancellation example
func (s *FileService) ProcessLargeFile(ctx context.Context, file *File) error {
// Process file with context support
if err := s.EmitWithContext(ctx, "file.processing", file); err != nil {
if errors.Is(err, context.Canceled) {
s.Logger.Info("File processing cancelled", logger.String("file", file.Name))
return err
}
return err
}
return nil
}
// Timeout pattern for slow operations
func (s *ReportService) GenerateReport(data *ReportData) error {
// Set reasonable timeout for report generation
err := s.Emitter.EmitWithTimeout("report.generate", data, 5*time.Minute)
if err != nil {
if errors.Is(err, context.DeadlineExceeded) {
s.Logger.Error("Report generation timed out")
return errors.New("report generation took too long")
}
return err
}
return nil
}
// HTTP request context integration
func (c *PostController) CreatePost(ctx *router.Context) error {
var req CreatePostRequest
if err := ctx.Bind(&req); err != nil {
return ctx.JSON(400, map[string]string{"error": "Invalid request"})
}
post := &Post{Title: req.Title, Content: req.Content}
// Use HTTP request context for event emission
requestCtx := ctx.Request.Context()
if err := c.Service.EmitWithContext(requestCtx, "post.creating", post); err != nil {
return ctx.JSON(500, map[string]string{"error": "Failed to process post"})
}
return ctx.JSON(201, post)
}
// Advanced timeout with cleanup
func (s *EmailService) SendBulkEmail(emails []EmailData) error {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()
// Create a channel to track progress
done := make(chan struct{})
go func() {
defer close(done)
for _, email := range emails {
if err := s.Emitter.EmitWithContext(ctx, "email.send", email); err != nil {
s.Logger.Error("Failed to send email", logger.String("error", err.Error()))
return
}
}
}()
select {
case "Bulk email sending completed")
return nil
case "Bulk email sending timed out")
return ctx.Err()
}
}
```
## Module Integration Patterns
Cross-Module Communication
Decoupled Module Design
```
// User module emits events
type UserService struct {
DB      *gorm.DB
Emitter *emitter.Emitter
Logger  logger.Logger
}
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
user := &User{
Email: req.Email,
Name:  req.Name,
}
if err := s.DB.Create(user).Error; err != nil {
return nil, err
}
// Emit events for other structures to handle
s.Emitter.Emit("user.created", user)
s.Emitter.EmitAsync("analytics.user_registered", user)
return user, nil
}
// Notification module listens for user events
type NotificationService struct {
Emitter *emitter.Emitter
Logger  logger.Logger
Email   email.Sender
}
func (s *NotificationService) Init() {
// Listen for user events without direct coupling
s.Emitter.On("user.created", func(data any) {
if user, ok := data.(*User); ok {
s.sendWelcomeEmail(user)
}
})
s.Emitter.On("user.password_reset", func(data any) {
if event, ok := data.(*PasswordResetEvent); ok {
s.sendPasswordResetEmail(event.User, event.Token)
}
})
}
// Analytics module tracks user behavior
type AnalyticsService struct {
Emitter *emitter.Emitter
Logger  logger.Logger
DB      *gorm.DB
}
func (s *AnalyticsService) Init() {
// Track various user events
s.Emitter.On("user.created", s.trackUserRegistration)
s.Emitter.On("user.login", s.trackUserLogin)
s.Emitter.On("post.created", s.trackContentCreation)
s.Emitter.On("order.completed", s.trackPurchase)
}
func (s *AnalyticsService) trackUserRegistration(data any) {
if user, ok := data.(*User); ok {
event := &AnalyticsEvent{
Type:      "user_registration",
UserID:    user.ID,
Timestamp: time.Now(),
Data:      map[string]any{"email": user.Email},
}
if err := s.DB.Create(event).Error; err != nil {
s.Logger.Error("Failed to track user registration",
logger.String("error", err.Error()))
}
}
}
// Post module with file handling
type PostService struct {
DB      *gorm.DB
Emitter *emitter.Emitter
Storage storage.Storage
Logger  logger.Logger
}
func (s *PostService) Init() {
// Listen for file upload events from storage
s.Emitter.On("post.featured_image.uploaded", func(data any) {
if post, ok := data.(*Post); ok {
s.Logger.Info("Featured image uploaded",
logger.Int("post_id", int(post.ID)),
logger.String("image", post.FeaturedImageURL))
}
})
// Listen for file deletion events
s.Emitter.On("post.featured_image.deleted", func(data any) {
if post, ok := data.(*Post); ok {
s.Logger.Info("Featured image deleted",
logger.Int("post_id", int(post.ID)))
}
})
}
func (s *PostService) CreatePost(req CreatePostRequest) (*Post, error) {
post := &Post{
Title:   req.Title,
Content: req.Content,
UserID:  req.UserID,
}
if err := s.DB.Create(post).Error; err != nil {
return nil, err
}
// Emit creation event for other structures
s.Emitter.Emit("post.created", post)
return post, nil
}
```
Module Registration with Emitter
Dependency Injection Pattern
```
// Module structure with emitter injection
type PostModule struct {
Service    *PostService
Controller *PostController
Router     *router.RouterGroup
Emitter    *emitter.Emitter
Logger     logger.Logger
}
func NewPostModule(
db *gorm.DB,
router *router.RouterGroup,
log logger.Logger,
emitter *emitter.Emitter,
storage *storage.ActiveStorage,
) module.Module {
service := &PostService{
DB:      db,
Emitter: emitter,
Logger:  log,
Storage: storage,
}
controller := &PostController{
Service: service,
Logger:  log,
}
// Initialize event listeners
service.Init()
// Register routes
router.GET("/posts", controller.GetPosts)
router.POST("/posts", controller.CreatePost)
router.PUT("/posts/:id", controller.UpdatePost)
router.DELETE("/posts/:id", controller.DeletePost)
return &PostModule{
Service:    service,
Controller: controller,
Router:     router,
Emitter:    emitter,
Logger:     log,
}
}
// App initialization with emitter
func (app *App) initModules() *App {
// Create shared emitter instance
emitter := emitter.New()
// Register structures with emitter dependency injection
authModule := authentication.NewAuthenticationModule(
app.db,
app.router.Group("/auth"),
app.emailSender,
app.logger,
emitter,
)
postModule := NewPostModule(
app.db,
app.router.Group("/api/v1"),
app.logger,
emitter,
app.storage,
)
notificationModule := notification.NewNotificationModule(
app.router.Group("/notifications"),
app.logger,
emitter,
app.emailSender,
)
// Store structures
app.modules = []module.Module{
authModule,
postModule,
notificationModule,
}
return app
}
// Module initialization helper
type ModuleInitializer struct {
DB          *gorm.DB
Router      *router.RouterGroup
Logger      logger.Logger
Emitter     *emitter.Emitter
Storage     *storage.ActiveStorage
EmailSender email.Sender
}
func (m *ModuleInitializer) InitializeAll() {
// Initialize all structures with shared dependencies
modules := []module.Module{
authentication.NewAuthenticationModule(m.DB, m.Router.Group("/auth"), m.EmailSender, m.Logger, m.Emitter),
media.NewMediaModule(m.DB, m.Router.Group("/media"), m.Storage, m.Emitter, m.Logger),
translation.NewTranslationModule(m.DB, m.Router.Group("/translations"), m.Logger, m.Emitter, m.Storage),
}
for _, mod := range modules {
if initializer, ok := mod.(interface{ Init() }); ok {
initializer.Init()
}
}
}
```
## Error Handling & Best Practices
Panic Recovery & Error Resilience
Robust Event Handling
```
// The emitter automatically recovers from panics in listeners
func (s *PostService) Init() {
s.Emitter.On("post.created", func(data any) {
// Even if this panics, other listeners will still execute
panic("something went wrong")
})
s.Emitter.On("post.created", func(data any) {
// This will still run despite the panic above
if post, ok := data.(*Post); ok {
s.Logger.Info("Post created", logger.Int("id", int(post.ID)))
}
})
}
// Best practice: Handle errors gracefully in listeners
func (s *EmailService) Init() {
s.Emitter.On("user.created", func(data any) {
user, ok := data.(*User)
if !ok {
s.Logger.Error("Invalid data type for user.created event")
return
}
if err := s.sendWelcomeEmail(user); err != nil {
s.Logger.Error("Failed to send welcome email",
logger.String("error", err.Error()),
logger.String("user_email", user.Email))
// Don't panic - log and continue
return
}
s.Logger.Info("Welcome email sent", logger.String("email", user.Email))
})
}
// Defensive event emission with error checking
func (s *OrderService) ProcessOrder(order *Order) error {
// Validate order before processing
if order == nil {
return errors.New("order cannot be nil")
}
// Process order logic
if err := s.DB.Save(order).Error; err != nil {
return fmt.Errorf("failed to save order: %w", err)
}
// Safe event emission with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := s.Emitter.EmitWithContext(ctx, "order.processed", order); err != nil {
// Log error but don't fail the entire operation
s.Logger.Error("Failed to emit order.processed event",
logger.String("error", err.Error()),
logger.Int("order_id", int(order.ID)))
}
return nil
}
// Testing event listeners with mock data
func TestPostServiceEvents(t *testing.T) {
emitter := emitter.New()
service := &PostService{Emitter: emitter}
var receivedPost *Post
emitter.On("post.created", func(data any) {
if post, ok := data.(*Post); ok {
receivedPost = post
}
})
testPost := &Post{ID: 1, Title: "Test Post"}
emitter.Emit("post.created", testPost)
assert.Equal(t, testPost, receivedPost)
}
// Event listener cleanup for tests
func (s *TestSuite) TearDown() {
// Clear all listeners after each test
s.emitter.Clear()
}
// Monitoring event listener performance
func (s *AnalyticsService) trackEventPerformance() {
s.Emitter.On("performance.track", func(data any) {
start := time.Now()
defer func() {
duration := time.Since(start)
s.Logger.Debug("Event processing time",
logger.String("duration", duration.String()))
}()
// Process analytics data
s.processAnalyticsData(data)
})
}
```
✅ Best Practices
- Use EmitAsync for non-critical operations like analytics
- Always validate data types in event listeners
- Handle errors gracefully - don't panic in listeners
- Use meaningful event names following module.action pattern
- Implement proper logging in event handlers
- Use context cancellation for long-running operations
- Clear listeners in tests to avoid interference
❌ Anti-Patterns
- Don't emit events for every database operation
- Avoid deep event chains that create circular dependencies
- Don't use events for synchronous data exchange
- Never assume listeners will complete successfully
- Don't emit sensitive data without sanitization
- Avoid blocking operations in async event handlers
- Don't ignore context cancellation in listeners
## Performance & Thread Safety
Performance Features
Thread-safe with read-write mutexes
Goroutine-based parallel listener execution
Automatic panic recovery in listeners
Context cancellation and timeout support
Zero allocation for listener management
Thread Safety
Concurrent emissions:
✓ Safe
Listener registration:
✓ Safe
Event introspection:
✓ Safe
Memory consistency:
✓ Guaranteed
Goroutine safety:
✓ Full support
Performance Monitoring Example
Event Performance Tracking
```
// Performance monitoring service
type EventMonitor struct {
emitter *emitter.Emitter
logger  logger.Logger
metrics map[string]*EventMetrics
mutex   sync.RWMutex
}
type EventMetrics struct {
TotalEmissions  int64
TotalListeners  int64
AverageLatency  time.Duration
ErrorCount      int64
}
func (m *EventMonitor) Init() {
m.metrics = make(map[string]*EventMetrics)
// Monitor all events by wrapping emitter methods
originalEmit := m.emitter.Emit
m.emitter.Emit = func(event string, data any) {
start := time.Now()
originalEmit(event, data)
m.recordMetrics(event, time.Since(start))
}
}
func (m *EventMonitor) recordMetrics(event string, duration time.Duration) {
m.mutex.Lock()
defer m.mutex.Unlock()
if _, exists := m.metrics[event]; !exists {
m.metrics[event] = &EventMetrics{}
}
metrics := m.metrics[event]
metrics.TotalEmissions++
metrics.TotalListeners = int64(m.emitter.ListenerCount(event))
// Calculate running average
if metrics.TotalEmissions == 1 {
metrics.AverageLatency = duration
} else {
metrics.AverageLatency = time.Duration(
(int64(metrics.AverageLatency)*metrics.TotalEmissions + int64(duration)) /
(metrics.TotalEmissions + 1),
)
}
}
// Get performance stats
func (m *EventMonitor) GetStats() map[string]*EventMetrics {
m.mutex.RLock()
defer m.mutex.RUnlock()
stats := make(map[string]*EventMetrics)
for event, metrics := range m.metrics {
stats[event] = &EventMetrics{
TotalEmissions:  metrics.TotalEmissions,
TotalListeners:  metrics.TotalListeners,
AverageLatency:  metrics.AverageLatency,
ErrorCount:      metrics.ErrorCount,
}
}
return stats
}
// Concurrent event emission test
func BenchmarkConcurrentEmit(b *testing.B) {
emitter := emitter.New()
// Register multiple listeners
for i := 0; i 10; i++ {
emitter.On("benchmark.event", func(data any) {
time.Sleep(time.Microsecond) // Simulate work
})
}
b.ResetTimer()
b.RunParallel(func(pb *testing.PB) {
for pb.Next() {
emitter.Emit("benchmark.event", "test data")
}
})
}
```
## Common Event Patterns
Event Naming Conventions
#### Entity Lifecycle Events
- • `user.created`
- • `user.updated`
- • `user.deleted`
- • `post.published`
- • `order.completed`
#### Action Events
- • `user.login_attempt`
- • `email.sent`
- • `file.uploaded`
- • `payment.processed`
- • `cache.invalidated`
#### File Events
- • `post.image.uploaded`
- • `user.avatar.deleted`
- • `document.file.processed`
#### System Events
- • `system.startup`
- • `analytics.track`
- • `audit.log`
💡 Event Design Tips
- Use dot notation for hierarchical event names (module.action.detail)
- Keep event data immutable - pass copies, not references to mutable objects
- Include sufficient context in event data for listeners to work independently
- Use past tense for completed actions (created, updated) and present for ongoing (creating, updating)
- Consider event versioning for breaking changes (user.created.v2)
- Document your events and their data structures for other developers