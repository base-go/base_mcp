---
title: Storage
description: Storage documentation for Base Framework.
---

# Storage
Flexible file storage system with multiple provider support and automatic file handling
## Overview
🗂️ Multi-Provider Storage
Base Framework's storage system provides a unified interface for file uploads with support for local storage, Amazon S3, and Cloudflare R2. It includes automatic file validation, size limits, and MIME type detection.
Local Storage
Store files directly on your server's filesystem with automatic directory creation.
Amazon S3
Scalable cloud storage with S3-compatible APIs and CDN support.
Cloudflare R2
Zero egress fee storage with global CDN and S3 compatibility.
## Configuration
Environment Variables
.env Configuration
```
# Storage Provider (local, s3, r2)
STORAGE_PROVIDER=local
# Local Storage
STORAGE_PATH=./uploads
STORAGE_BASE_URL=http://localhost:3000/uploads
# Amazon S3
STORAGE_S3_KEY=your-access-key
STORAGE_S3_SECRET=your-secret-key
STORAGE_S3_BUCKET=your-bucket
STORAGE_S3_REGION=us-east-1
STORAGE_S3_CDN=https://cdn.example.com
# Cloudflare R2
STORAGE_R2_ACCOUNT_ID=your-account-id
STORAGE_R2_KEY=your-access-key
STORAGE_R2_SECRET=your-secret-key
STORAGE_R2_BUCKET=your-bucket
STORAGE_R2_CDN=https://pub-xxx.r2.dev
```
## File Types & Validation
Base Framework automatically validates file types and sizes based on the field type in your models.
Image Type
Max Size:
5MB
Extensions:
.jpg
.jpeg
.png
.gif
.webp
File Type
Max Size:
50MB
Extensions:
.pdf
.doc
.docx
.xls
.xlsx
.ppt
Attachment Type
Max Size:
10MB
Extensions:
.zip
.rar
.7z
.tar
+ all file types
Custom Types
You can define custom file types with specific extensions and size limits in your models.
## Usage Examples
Model Definition
models/product.go
```
package models
import "base/core/storage"
type Product struct {
ID          uint                 `gorm:"primaryKey"`
Name        string               `gorm:"not null"`
Description string               `gorm:"type:text"`
Price       float64              `gorm:"not null"`
// Single image field
Thumbnail   storage.Attachment   `gorm:"foreignKey:ModelId;references:Id" field_type:"image"`
// Multiple images
Gallery     []storage.Attachment `gorm:"foreignKey:ModelId;references:Id" field_type:"image" multiple:"true"`
// File attachments
Manual      storage.Attachment   `gorm:"foreignKey:ModelId;references:Id"`
Downloads   []storage.Attachment `gorm:"foreignKey:ModelId;references:Id" multiple:"true"`
}
```
File Upload in Service
services/product_service.go
```
func (s *ProductService) Create(data CreateProductRequest) (*Product, error) {
product := &Product{
Name:        data.Name,
Description: data.Description,
Price:       data.Price,
}
// Handle thumbnail upload
if data.Thumbnail != nil {
attachment, err := s.Storage.Upload(data.Thumbnail, "thumbnail", product)
if err != nil {
return nil, err
}
product.Thumbnail = *attachment
}
// Handle multiple gallery images
if len(data.Gallery) > 0 {
for _, file := range data.Gallery {
attachment, err := s.Storage.Upload(file, "gallery", product)
if err != nil {
continue // or handle error
}
product.Gallery = append(product.Gallery, *attachment)
}
}
// Save product with attachments
if err := s.DB.Create(product).Error; err != nil {
return nil, err
}
// Emit upload event
s.Emitter.Emit("product.uploaded", product)
return product, nil
}
```
File Deletion
Deleting Files
```
// Delete single file
if product.Thumbnail.Path != "" {
err := s.Storage.Delete(product.Thumbnail.Path)
if err != nil {
s.Logger.Error("Failed to delete thumbnail", "error", err)
}
}
// Delete multiple files
for _, attachment := range product.Gallery {
if err := s.Storage.Delete(attachment.Path); err != nil {
s.Logger.Error("Failed to delete gallery image", "error", err)
}
}
// Clear from database
product.Thumbnail = storage.Attachment{}
product.Gallery = []storage.Attachment{}
s.DB.Save(product)
```
## API Response
JSON Response with Attachments
GET /api/products/1
```
{
"id": 1,
"name": "Premium Laptop",
"description": "High-performance laptop for professionals",
"price": 1299.99,
"thumbnail": {
"id": 1,
"filename": "laptop-thumbnail.jpg",
"path": "products/thumbnails/laptop-thumbnail.jpg",
"size": 245632,
"url": "https://cdn.example.com/products/thumbnails/laptop-thumbnail.jpg",
"created_at": "2024-01-15T10:30:00Z"
},
"gallery": [
{
"id": 2,
"filename": "laptop-angle1.jpg",
"url": "https://cdn.example.com/products/gallery/laptop-angle1.jpg"
},
{
"id": 3,
"filename": "laptop-angle2.jpg",
"url": "https://cdn.example.com/products/gallery/laptop-angle2.jpg"
}
],
"manual": {
"id": 4,
"filename": "user-manual.pdf",
"size": 2458976,
"url": "https://cdn.example.com/products/manuals/user-manual.pdf"
}
}
```
## Storage Events
The storage system emits events that you can listen to for additional processing:
Upload Events
`{model}.{field}.uploaded` - File uploaded
`{model}.{field}.processing` - File processing
`{model}.{field}.completed` - Upload complete
Deletion Events
`{model}.{field}.deleted` - File deleted
`{model}.{field}.cleanup` - Cleanup performed
`{model}.{field}.error` - Error occurred
## Best Practices
✅ Do's
- Validate file types and sizes before upload
- Use CDN URLs for production environments
- Clean up orphaned files when records are deleted
- Use transactions when uploading multiple files
- Implement virus scanning for user uploads
❌ Don'ts
- Don't store sensitive files without encryption
- Don't allow unlimited file sizes
- Don't trust client-side file validation alone
- Don't expose internal file paths in APIs
- Don't forget to set proper CORS headers for CDN