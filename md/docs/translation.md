---
title: Translation
description: Translation documentation for Base Framework.
---

# Translation
Comprehensive internationalization (i18n) system for multi-language applications with flexible field translation
## Overview
🌍 Multi-Language Support
Base Framework's translation system provides a flexible, database-driven internationalization solution. It supports both static text translation and dynamic field translation at the model level, allowing you to create truly multi-language applications.
Database-Driven
Store translations directly in your database with automatic loading and caching for optimal performance.
Field-Level Translation
Translate specific model fields with automatic fallback to original values when translations are missing.
Bulk Operations
Efficiently update multiple translations at once with transaction-based bulk operations.
REST API
Complete REST API for managing translations with filtering, pagination, and language-specific endpoints.
## Key Features
Translation Field Type
- Custom Field type that handles multiple languages automatically
- JSON marshaling that adapts based on available translations
- Database storage with original value preservation
- Automatic translation loading and caching
Module Integration
- Translatable interface for modules to define translated fields
- Helper functions for retrieving module translations
- Automatic response enrichment with translated content
- Service layer integration for seamless usage
Advanced Operations
- Bulk translation updates with transaction support
- Language discovery and supported language listing
- Model-specific translation management
- Translation history and versioning support
## CLI Generation
🚀 Quick Setup
Use the Base CLI to automatically generate models with translation fields. The CLI handles all the boilerplate code for you, including automatic translation loading in service methods.
Generate Model with Translation Fields
Terminal
```
# Generate a Post model with translatable title field
./cmd g Post title:translation desc:text feat:image
# Generate a Product model with translatable fields
./cmd g Product name:translation description:translation price:float category_id:uint
# Generate a Page model with multiple translatable fields
./cmd g Page title:translation content:translation meta_title:translation meta_description:translation slug:string
```
Generated Model Structure
When you use `title:translation`, the CLI automatically generates:
models/post.go
```
type Post struct {
Id        uint                `json:"id" gorm:"primarykey"`
CreatedAt time.Time           `json:"created_at"`
UpdatedAt time.Time           `json:"updated_at"`
DeletedAt gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
Title     translation.Field   `json:"title"`  // Translatable field
Desc      string              `json:"desc"`
FeatId    uint                `json:"feat_id,omitempty"`
Feat      *storage.Attachment `json:"feat,omitempty"`
}
// Request types use string for input
type CreatePostRequest struct {
Title string `json:"title" binding:"required"`
Desc  string `json:"desc"`
}
// Response types return translation.Field with full translation object
type PostResponse struct {
Id        uint              `json:"id"`
CreatedAt time.Time         `json:"created_at"`
UpdatedAt time.Time         `json:"updated_at"`
DeletedAt gorm.DeletedAt    `json:"deleted_at"`
Title     translation.Field `json:"title"`  // Returns translation object
Desc      string            `json:"desc"`
}
```
Automatic Service Integration
The CLI automatically generates service methods that load translations:
posts/service.go - Auto-generated
```
func (s *PostService) GetById(id uint) (*models.Post, error) {
var item models.Post
if err := s.DB.Preload(clause.Associations).First(&item, id).Error; err != nil {
return nil, err
}
// Automatically load translations for translation fields
s.loadTranslationsForItem(&item)
return &item, nil
}
func (s *PostService) loadTranslationsForItem(item *models.Post) error {
if item == nil {
return nil
}
modelName := item.GetModelName()
modelId := item.GetId()
// Load translations for Title field
if err := s.TranslationHelper.Service.LoadTranslationsForField(&item.Title, modelName, modelId, "title"); err != nil {
s.Logger.Error("Failed to load translations for Title", logger.String("error", err.Error()))
}
return nil
}
```
## Translation Field Usage
Model Definition
models/post.go
```
package models
import (
"base/core/translation"
"time"
"gorm.io/gorm"
)
type Post struct {
ID          uint                   `gorm:"primarykey"`
CreatedAt   time.Time
UpdatedAt   time.Time
DeletedAt   gorm.DeletedAt         `gorm:"index"`
// Translatable fields
Title       translation.Field      `gorm:"type:text"`
Content     translation.Field      `gorm:"type:text"`
Description translation.Field      `gorm:"type:text"`
// Regular fields
AuthorID    uint                   `gorm:"not null"`
Published   bool                   `gorm:"default:false"`
ViewCount   int                    `gorm:"default:0"`
}
```
Working with Translation Fields
Basic Field Operations
```
// Create a new post with original text
post := &Post{
Title:       translation.NewField("Hello World"),
Content:     translation.NewField("This is the original content"),
Description: translation.NewField("Original description"),
}
// Set translations for different languages
post.Title.SetTranslation("es", "Hola Mundo")
post.Title.SetTranslation("fr", "Bonjour le Monde")
// Get translation for specific language with fallback
titleInSpanish := post.Title.GetTranslationOrOriginal("es")  // Returns: "Hola Mundo"
titleInGerman := post.Title.GetTranslationOrOriginal("de")   // Returns: "Hello World" (fallback)
// Check if translation exists
hasSpanish := post.Title.HasTranslation("es")  // Returns: true
hasGerman := post.Title.HasTranslation("de")   // Returns: false
// Get all available languages for a field
languages := post.Title.GetAvailableLanguages()  // Returns: ["es", "fr"]
```
## Service Integration
Using Translation Service
Generated Service Pattern
```
// All generated services follow this pattern
type YourEntityService struct {
DB                *gorm.DB
Emitter           *emitter.Emitter
Storage           *storage.ActiveStorage
Logger            logger.Logger
TranslationHelper *translation.Helper  // Automatically injected
}
func (s *YourEntityService) Create(req *models.CreateYourEntityRequest) (*models.YourEntity, error) {
item := &models.YourEntity{
// For translation fields, use translation.NewField()
TranslatableField: translation.NewField(req.TranslatableField),
// For regular fields, use directly
RegularField: req.RegularField,
}
if err := s.DB.Create(item).Error; err != nil {
return nil, err
}
return s.GetById(item.Id)  // Automatically loads translations
}
func (s *YourEntityService) GetById(id uint) (*models.YourEntity, error) {
item := &models.YourEntity{}
if err := s.DB.First(item, id).Error; err != nil {
return nil, err
}
// Automatically load translations for all translatable fields
if err := s.loadTranslationsForItem(item); err != nil {
s.Logger.Error("Failed to load translations", logger.String("error", err.Error()))
}
return item, nil
}
```
Loading Translations
Auto-Generated Translation Loading
```
// This method is automatically generated for each entity
func (s *YourEntityService) loadTranslationsForItem(item *models.YourEntity) error {
if item == nil {
return nil
}
modelName := item.GetModelName()
modelId := item.GetId()
// Load translations for each translatable field
// This loops through all fields marked with 'translation' type
if err := s.TranslationHelper.Service.LoadTranslationsForField(&item.TranslatableField1, modelName, modelId, "translatable_field1"); err != nil {
s.Logger.Error("Failed to load translations", logger.String("error", err.Error()))
}
if err := s.TranslationHelper.Service.LoadTranslationsForField(&item.TranslatableField2, modelName, modelId, "translatable_field2"); err != nil {
s.Logger.Error("Failed to load translations", logger.String("error", err.Error()))
}
return nil
}
// Using the Translation Helper directly
func (s *YourEntityService) AddTranslation(id uint, field, language, value string) error {
modelName := "your_entity"  // matches your model name
return s.TranslationHelper.SetTranslation(modelName, id, field, value, language)
}
// Bulk set multiple translations
func (s *YourEntityService) BulkAddTranslations(id uint, language string, translations map[string]string) error {
modelName := "your_entity"
return s.TranslationHelper.BulkSetTranslations(modelName, id, language, translations)
}
```
## Module Integration
Module Integration (Auto-Generated)
Generated Module Pattern
```
// All generated modules follow this pattern
func Init(deps module.Dependencies) module.Module {
// Translation service and helper are automatically created
translationService := translation.NewTranslationService(deps.DB, deps.Emitter, deps.Storage, deps.Logger)
translationHelper := translation.NewHelper(translationService)
// Service gets translation helper injected
service := NewYourEntityService(deps.DB, deps.Emitter, deps.Storage, deps.Logger, translationHelper)
controller := NewYourEntityController(service, deps.Storage)
mod := &Module{
DB:                deps.DB,
Service:           service,
Controller:        controller,
TranslationHelper: translationHelper,
}
return mod
}
// Module structure includes translation helper
type Module struct {
DB                *gorm.DB
Service           *YourEntityService
Controller        *YourEntityController
TranslationHelper *translation.Helper
}
```
Using Translation Helper
Using Helper Functions
```
// Frontend usage - accessing translation values
// From API response: {"title": {"original": "Title", "de": "Titel", "sq": "Titulli"}}
// Get specific language or fallback to original
function getTranslation(field, language) {
return field[language] || field.original;
}
// Usage examples
const title_de = getTranslation(post.title, "de");  // "Titel"
const title_en = getTranslation(post.title, "en");  // "Title" (fallback to original)
// Backend helper usage
helper := translation.NewHelper(translationService)
// Set a single translation
err = helper.SetTranslation("your_entity", 1, "field_name", "Translated Value", "de")
// Bulk set multiple translations
translationMap := map[string]string{
"title": "Deutscher Titel",
"description": "Deutsche Beschreibung",
}
err = helper.BulkSetTranslations("your_entity", 1, "de", translationMap)
```
## API Reference
Translation Endpoints
GET
`/translations`
List all translations with pagination and filtering
POST
`/translations`
Create a new translation
GET
`/translations/by-id/{id}`
Get a specific translation by ID
PUT
`/translations/by-id/{id}`
Update an existing translation
DELETE
`/translations/by-id/{id}`
Delete a translation
POST
`/translations/bulk`
Bulk update multiple translations
GET
`/translations/languages`
Get supported languages
GET
`/translations/models/{model}/{id}`
Get all translations for a model
API Request Examples
Create Translation
```
POST /api/translations
Content-Type: application/json
{
"key": "title",
"value": "Hola Mundo",
"model": "post",
"model_id": 1,
"language": "es"
}
// Bulk Update Translations
POST /api/translations/bulk
Content-Type: application/json
{
"model": "post",
"model_id": 1,
"language": "es",
"translations": {
"title": "Título en Español",
"content": "Contenido en Español",
"description": "Descripción en Español"
}
}
// Get Model Translations
GET /api/translations/models/post/1/es
// List with Filters
GET /api/translations?model=post&model_id=1&page=1&limit=20
```
## JSON Response Format
Translation Field JSON Output
Adaptive JSON Marshaling
```
// When no translations are loaded (original value only)
{
"title": "Hello World",
"content": "This is the original content"
}
// When translations are loaded (translation map)
{
"title": {
"es": "Hola Mundo",
"fr": "Bonjour le Monde",
"de": "Hallo Welt"
},
"content": {
"es": "Este es el contenido original",
"fr": "Ceci est le contenu original"
}
}
// Actual API response (returns full translation field structure)
GET /api/posts/1
{
"id": 1,
"created_at": "2025-08-28T20:37:28.986523+02:00",
"updated_at": "2025-08-28T20:37:28.986523+02:00",
"deleted_at": null,
"title": {
"original": "Title",
"de": "Titel",
"sq": "Titulli"
},
"desc": "description" // original value since no translation
}
```
## Best Practices
✅ Do's
- Always provide meaningful original values in your default language
- Use standard language codes (ISO 639-1) like "en", "es", "fr"
- Implement fallback logic to handle missing translations gracefully
- Use bulk operations for updating multiple translations efficiently
- Implement the Translatable interface for modules with translated fields
- Load translations only when needed to optimize performance
- Use transactions when updating multiple translations together
⚠️ Considerations
- Translation fields store original values in the main table for performance
- Translations are stored separately and loaded on demand
- Consider caching strategies for frequently accessed translations
- Plan your database indexing strategy for translation queries
- Be mindful of character limits and text expansion in translations
❌ Don'ts
- Don't rely solely on translations without original fallbacks
- Don't use translation fields for data that doesn't need localization
- Don't forget to handle cases where translations might be empty
- Don't load all translations when only specific languages are needed
- Don't hardcode language codes in your application logic
- Don't forget to clean up orphaned translations when deleting records
## Advanced Usage
Querying with Translations
Advanced Query Examples
```
// Search in original values (faster)
var posts []Post
db.Where("title LIKE ?", "%search%").Find(&posts)
// Search in translations (requires join)
db.Joins("LEFT JOIN translations t ON t.model = 'post' AND t.model_id = posts.id").
Where("t.key = 'title' AND t.language = 'es' AND t.value LIKE ?", "%búsqueda%").
Find(&posts)
// Get posts with specific translation availability
db.Joins("INNER JOIN translations t ON t.model = 'post' AND t.model_id = posts.id").
Where("t.language = ?", "es").
Distinct("posts.*").
Find(&posts)
// Count translations by language
type LanguageCount struct {
Language string
Count    int64
}
var counts []LanguageCount
db.Model(&Translation{}).
Select("language, count(*) as count").
Where("model = ?", "post").
Group("language").
Find(&counts)
```
Migration and Setup
Database Migration
```
// The Translation module automatically creates this table:
CREATE TABLE translations (
id SERIAL PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
deleted_at TIMESTAMP,
`key` VARCHAR(255) NOT NULL,
value TEXT NOT NULL,
model VARCHAR(255) NOT NULL,
model_id INTEGER NOT NULL,
language VARCHAR(5) NOT NULL,
INDEX idx_translation_lookup (`key`, model, model_id, language),
INDEX idx_deleted_at (deleted_at)
);
// Register the translation module in app/init.go
func RegisterModules(deps module.Dependencies) error {
// Core modules
translationModule := translation.NewTranslationModule(
deps.DB, deps.Router, deps.Logger, deps.Emitter, deps.Storage,
)
module.RegisterModule("translation", translationModule)
// Your app modules...
return nil
}
```