---
title: Scheduler
description: Scheduler documentation for Base Framework.
---

# Task Scheduler
Built-in task scheduling system with support for interval-based and cron-based scheduling
## Overview
Base Framework includes a powerful task scheduling system that allows you to run background
tasks at specific intervals or using cron expressions. The scheduler supports two types of
scheduling:
- **Interval-based scheduling** - Run tasks at regular intervals (daily, monthly,
custom intervals)
- **Cron-based scheduling** - Use cron expressions for complex scheduling
patterns
## Features
### Multiple Schedule Types
Support for daily, monthly, interval, and cron-based scheduling patterns
### Task Management
Enable/disable tasks, run tasks manually, and monitor execution statistics
### Error Handling
Built-in error tracking, timeout protection, and comprehensive logging
### Statistics
Track run counts, error counts, last run times, and next scheduled runs
### HTTP API
RESTful API endpoints for managing tasks programmatically
### Concurrent Execution
Tasks run concurrently with proper context handling and timeout management
## Quick Start
### 1. Generate a Task with CLI
```bash
# Generate a new scheduled task
base scheduler generate posts cleanup-old-posts
base scheduler g users send-weekly-digest
# Works with different naming conventions
base scheduler g Post publish  # Finds 'posts' module
base scheduler g posts publish # Direct match
```
### 2. Implement Your Task Logic
```go
// Generated task file: cleanup_old_posts_task.go
func (t *CleanupOldPostsTask) execute(ctx context.Context) error {
t.logger.Info("Starting cleanup old posts task")
// Check for cancellation
select {
case
### 3. Register and Manage Tasks
```go
// Register in your module initialization
func (m *PostsModule) Start() error {
cleanupTask := NewCleanupOldPostsTask(m.Logger)
cleanupTask.RegisterTask(m.Scheduler)
return nil
}
// Or use CLI to manage tasks
// base scheduler list --api-key=your-key
// base scheduler run cleanup-old-posts --api-key=your-key
// base scheduler enable cleanup-old-posts --api-key=your-key
```
## Schedule Types
### Daily Schedule
Run tasks daily at a specific time.
```go
schedule := &scheduler.DailySchedule{
Hour: 14, // 2 PM
Minute: 30, // 30 minutes
}
// Runs every day at 2:30 PM
```
### Monthly Schedule
Run tasks monthly on a specific day and time.
```go
schedule := &scheduler.MonthlySchedule{
Day: 1, // 1st of the month
Hour: 9, // 9 AM
Minute: 0, // 0 minutes
}
// Runs on the 1st of each month at 9:00 AM
```
### Interval Schedule
Run tasks at regular intervals.
```go
schedule := &scheduler.IntervalSchedule{
Interval: 30 * time.Minute, // Every 30 minutes
}
// Other examples:
// 5 * time.Second // Every 5 seconds
// time.Hour // Every hour
// 24 * time.Hour // Every 24 hours
```
### Cron Schedule
Use cron expressions for complex scheduling patterns.
```go
// Using CronScheduler for cron
expressions
cronTask := &scheduler.CronTask{
Name: "weekly-report",
Description: "Generate weekly reports",
CronExpr: "0 0 9 * * MON", // Every Monday at 9:00 AM
Handler: func(ctx context.Context) error {
// Generate reports
return nil
},
Enabled: true,
}
cronScheduler.RegisterTask(cronTask)
```
#### Common Cron Expressions:
`0 0 * * * *`
Every hour
`0 0 0 * * *`
Every day at midnight
`0 0 9 * * MON-FRI`
Weekdays at 9 AM
`0 */15 * * * *`
Every 15 minutes
## Task Management
### Enable/Disable Tasks
```go
// Enable a task
err := scheduler.EnableTask("task-name")
// Disable a task
err := scheduler.DisableTask("task-name")
```
### Run Tasks Manually
```go
// Run a task immediately (bypasses
schedule)
err := scheduler.RunTaskNow("task-name")
if err != nil {
log.Error("Failed to run task", err)
}
```
### Get Task Information
```go
// Get a specific task
task, exists := scheduler.GetTask("task-name")
if exists {
fmt.Printf("Task: %s, Enabled: %t\n", task.Name, task.Enabled)
}
// Get all tasks
tasks := scheduler.GetAllTasks()
for name, task := range tasks {
fmt.Printf("Task: %s, Run Count: %d\n", name, task.RunCount)
}
```
## CLI Management
Use the Base CLI to generate and manage scheduled tasks from the command line.
### Task Generation
```bash
# Generate task in specific module
base scheduler generate posts cleanup-old-posts
base scheduler g users send-weekly-digest
base scheduler g core backup-database
# Smart module detection (case-insensitive, pluralization)
base scheduler g Post publish    # Finds 'posts' module
base scheduler g User reminder   # Finds 'users' module
```
### Task Management Commands
```bash
# List all tasks
base scheduler list --api-key=your-key
# Run task immediately
base scheduler run cleanup-old-posts --api-key=your-key
# Enable/disable tasks
base scheduler enable cleanup-old-posts --api-key=your-key
base scheduler disable cleanup-old-posts --api-key=your-key
# Get scheduler status
base scheduler status --api-key=your-key
```
### Generated Task Structure
Each generated task includes a complete, ready-to-use structure:
- • **Task struct** with logger integration
- • **Constructor function** for dependency injection
- • **RegisterTask method** for interval-based scheduling
- • **RegisterCronTask method** for cron-based scheduling
- • **Execute function** with context cancellation support
- • **GetTaskInfo method** for metadata
## HTTP API
The scheduler provides RESTful API endpoints for managing tasks programmatically.
### Available Endpoints
Method
Endpoint
Description
GET
/api/scheduler/status
Get scheduler status and statistics
GET
/api/scheduler/tasks
List all registered tasks
GET
/api/scheduler/tasks/{name}
Get specific task details
POST
/api/scheduler/tasks/{name}/run
Run a task immediately
PUT
/api/scheduler/tasks/{name}/enable
Enable a task
PUT
/api/scheduler/tasks/{name}/disable
Disable a task
GET
/api/scheduler/stats
Get detailed scheduler statistics
### API Usage Examples
#### Get All Tasks
```bash
curl -X GET
"http://localhost:8100/api/scheduler/tasks" \
-H "X-Api-Key: your-api-key"
```
#### Run Task Immediately
```bash
curl -X POST
"http://localhost:8100/api/scheduler/tasks/daily-cleanup/run" \
-H "X-Api-Key: your-api-key"
```
#### Enable/Disable Task
```bash
# Enable task
curl -X PUT "http://localhost:8100/api/scheduler/tasks/daily-cleanup/enable" \
-H "X-Api-Key: your-api-key"
# Disable task
curl -X PUT "http://localhost:8100/api/scheduler/tasks/daily-cleanup/disable" \
-H "X-Api-Key: your-api-key"
```
## Best Practices
### Task Design
- • Keep tasks idempotent - they should be safe to run multiple times
- • Use meaningful task names and descriptions
- • Handle errors gracefully and log important information
- • Avoid long-running tasks (use timeouts)
### Error Handling
```go
Handler: func(ctx context.Context) error
{
// Check if context is cancelled
select {
case
### Monitoring
- • Monitor task execution statistics regularly
- • Set up alerts for tasks with high error rates
- • Use the HTTP API to integrate with monitoring systems
- • Check logs for task execution details
### Performance
- • Tasks run concurrently, so design for thread safety
- • Use appropriate check intervals for the scheduler
- • Consider using cron scheduler for complex patterns
- • Disable tasks that are no longer needed
## Complete Example
### Email Notification Service with Scheduled Tasks
```go
package main
import (
"context"
"fmt"
"time"
"base/core/scheduler"
"base/core/logger"
)
type EmailService struct {
scheduler *scheduler.Scheduler
logger logger.Logger
}
func NewEmailService(s *scheduler.Scheduler, log logger.Logger) *EmailService {
service := &EmailService{
scheduler: s,
logger: log,
}
// Register tasks
service.registerTasks()
return service
}
func (e *EmailService) registerTasks() {
// Daily digest email
dailyDigest := &scheduler.Task{
Name: "daily-digest",
Description: "Send daily digest emails to users",
Schedule: &scheduler.DailySchedule{Hour: 8, Minute: 0}, // 8:00 AM
Handler: e.sendDailyDigest,
Enabled: true,
}
// Weekly report
weeklyReport := &scheduler.Task{
Name: "weekly-report",
Description: "Generate and send weekly reports",
Schedule: &scheduler.DailySchedule{Hour: 9, Minute: 0}, // 9:00 AM on Mondays
Handler: e.sendWeeklyReport,
Enabled: true,
}
// Cleanup old emails every 6 hours
cleanup := &scheduler.Task{
Name: "email-cleanup",
Description: "Clean up old email logs",
Schedule: &scheduler.IntervalSchedule{Interval: 6 * time.Hour},
Handler: e.cleanupOldEmails,
Enabled: true,
}
// Register all tasks
e.scheduler.RegisterTask(dailyDigest)
e.scheduler.RegisterTask(weeklyReport)
e.scheduler.RegisterTask(cleanup)
}
func (e *EmailService) sendDailyDigest(ctx context.Context) error {
e.logger.Info("Sending daily digest emails...")
// Check for cancellation
select {
case