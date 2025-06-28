# Blog Application - Complete API Endpoints Reference

## Overview
Complete API specification for the production-ready multi-author blog platform with role-based access control, content management, and advanced features.

---

## **Authentication & User Management**

### **Public Authentication Endpoints**

#### **POST** `/api/v1/auth/register` ✅
**Description**: Register a new user account  
**Access**: Public  
**Request Body**:
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```
**Response**: `201 Created`
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### **POST** `/api/v1/auth/login` ✅
**Description**: Authenticate user and set JWT cookies  
**Access**: Public  
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```
**Response**: `200 OK` + HTTP-only cookies set
```json
{
  "message": "Login successful",
  "data": {
    "user": { /* user object */ },
    "roles": ["author"]
  }
}
```

#### **POST** `/api/v1/auth/logout` ✅
**Description**: Clear authentication cookies  
**Access**: Public  
**Response**: `200 OK`

#### **POST** `/api/v1/auth/refresh` ✅
**Description**: Refresh access token using refresh token  
**Access**: Public (requires refresh token cookie)  
**Response**: `200 OK` + new access token cookie

---

## **User Profile Management**

### **Protected User Endpoints**

#### **GET** `/api/v1/profile` 
**Description**: Get current user's profile  
**Access**: Authenticated users  
**Response**: `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Content creator and blogger",
    "avatar_url": "https://s3.amazonaws.com/bucket/avatar.jpg",
    "roles": ["author"],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### **PUT** `/api/v1/profile`
**Description**: Update current user's profile  
**Access**: Authenticated users  
**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "bio": "Updated bio",
  "avatar_url": "https://s3.amazonaws.com/bucket/new-avatar.jpg"
}
```
**Response**: `200 OK`

#### **PUT** `/api/v1/profile/password`
**Description**: Change user password  
**Access**: Authenticated users  
**Request Body**:
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```
**Response**: `200 OK`

---

## **Role & Permission Management**

### **Role Management (Admin Only)**

#### **GET** `/api/v1/roles`
**Description**: Get all roles with their permissions  
**Access**: `users.manage` permission  
**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "admin",
      "description": "System administrator",
      "permissions": ["posts.create", "posts.delete", "users.manage"],
      "user_count": 5,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### **POST** `/api/v1/roles`
**Description**: Create a new role  
**Access**: `system.admin` permission  
**Request Body**:
```json
{
  "name": "editor",
  "description": "Content editor with advanced permissions",
  "permissions": ["posts.create", "posts.update", "posts.publish"]
}
```
**Response**: `201 Created`

#### **GET** `/api/v1/roles/:roleId/permissions`
**Description**: Get permissions for a specific role  
**Access**: `users.manage` permission  
**Response**: `200 OK`

#### **POST** `/api/v1/roles/:roleId/permissions`
**Description**: Assign permission to role  
**Access**: `system.admin` permission  
**Request Body**:
```json
{
  "permission": "posts.publish"
}
```
**Response**: `200 OK`

#### **DELETE** `/api/v1/roles/:roleId/permissions/:permissionId`
**Description**: Remove permission from role  
**Access**: `system.admin` permission  
**Response**: `200 OK`

### **User Role Management**

#### **GET** `/api/v1/users/:userId/roles`
**Description**: Get roles assigned to a user  
**Access**: `users.manage` permission or own profile  
**Response**: `200 OK`
```json
{
  "data": {
    "user_id": "uuid",
    "roles": [
      {
        "id": "uuid",
        "name": "author",
        "assigned_at": "2025-01-01T00:00:00Z",
        "assigned_by": "uuid"
      }
    ]
  }
}
```

#### **POST** `/api/v1/users/:userId/roles`
**Description**: Assign role to user  
**Access**: `users.manage` permission  
**Request Body**:
```json
{
  "role_name": "admin"
}
```
**Response**: `200 OK`

#### **DELETE** `/api/v1/users/:userId/roles/:roleId`
**Description**: Remove role from user  
**Access**: `users.manage` permission  
**Response**: `200 OK`

### **Permission Reference**

#### **GET** `/api/v1/permissions`
**Description**: Get all available permissions  
**Access**: `users.manage` permission  
**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "posts.create",
      "resource": "posts",
      "action": "create",
      "description": "Create new blog posts"
    }
  ]
}
```

---

## **Content Management**

### **Categories Management**

#### **GET** `/api/v1/categories`
**Description**: Get all categories  
**Access**: Public  
**Query Parameters**: `?include_count=true` (include post count)  
**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Technology",
      "slug": "technology",
      "description": "Tech-related posts",
      "color": "#FF5733",
      "post_count": 42,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### **POST** `/api/v1/categories`
**Description**: Create a new category  
**Access**: `categories.manage` permission  
**Request Body**:
```json
{
  "name": "Technology",
  "description": "Technology and programming posts",
  "color": "#FF5733"
}
```
**Response**: `201 Created`

#### **GET** `/api/v1/categories/:categoryId`
**Description**: Get category by ID  
**Access**: Public  
**Response**: `200 OK`

#### **PUT** `/api/v1/categories/:categoryId`
**Description**: Update category  
**Access**: `categories.manage` permission  
**Response**: `200 OK`

#### **DELETE** `/api/v1/categories/:categoryId`
**Description**: Delete category  
**Access**: `categories.manage` permission  
**Response**: `200 OK`

### **Tags Management**

#### **GET** `/api/v1/tags`
**Description**: Get all tags  
**Access**: Public  
**Query Parameters**: `?search=keyword&limit=50`  
**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "javascript",
      "slug": "javascript",
      "post_count": 15,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### **POST** `/api/v1/tags`
**Description**: Create a new tag  
**Access**: `tags.manage` permission  
**Request Body**:
```json
{
  "name": "JavaScript"
}
```
**Response**: `201 Created`

#### **GET** `/api/v1/tags/:tagId`
**Description**: Get tag by ID with related posts  
**Access**: Public  
**Response**: `200 OK`

#### **PUT** `/api/v1/tags/:tagId`
**Description**: Update tag  
**Access**: `tags.manage` permission  
**Response**: `200 OK`

#### **DELETE** `/api/v1/tags/:tagId`
**Description**: Delete tag  
**Access**: `tags.manage` permission  
**Response**: `200 OK`

---

## **Blog Post Management**

### **Public Post Endpoints**

#### **GET** `/api/v1/posts`
**Description**: Get published posts with filtering and pagination  
**Access**: Public  
**Query Parameters**:
- `?page=1&limit=10` - Pagination
- `?category=technology` - Filter by category slug
- `?tag=javascript` - Filter by tag slug
- `?author=johndoe` - Filter by author username
- `?search=keyword` - Search in title and content
- `?sort=created_at&order=desc` - Sorting options

**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Getting Started with Go",
      "slug": "getting-started-with-go",
      "excerpt": "Learn the basics of Go programming...",
      "featured_image_url": "https://s3.amazonaws.com/bucket/image.jpg",
      "status": "published",
      "author": {
        "id": "uuid",
        "username": "johndoe",
        "first_name": "John",
        "last_name": "Doe",
        "avatar_url": "https://s3.amazonaws.com/bucket/avatar.jpg"
      },
      "category": {
        "id": "uuid",
        "name": "Technology",
        "slug": "technology"
      },
      "tags": [
        {"id": "uuid", "name": "go", "slug": "go"},
        {"id": "uuid", "name": "programming", "slug": "programming"}
      ],
      "published_at": "2025-01-01T00:00:00Z",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false
  }
}
```

#### **GET** `/api/v1/posts/:slug`
**Description**: Get single post by slug (full content)  
**Access**: Public  
**Response**: `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "title": "Getting Started with Go",
    "slug": "getting-started-with-go",
    "content": "Full HTML content of the post...",
    "excerpt": "Learn the basics...",
    "featured_image_url": "https://s3.amazonaws.com/bucket/image.jpg",
    "status": "published",
    "author": { /* author object */ },
    "category": { /* category object */ },
    "tags": [ /* tag objects */ ],
    "published_at": "2025-01-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

### **Author Post Management**

#### **GET** `/api/v1/posts/drafts`
**Description**: Get current user's draft posts  
**Access**: `posts.create` permission  
**Response**: `200 OK`

#### **POST** `/api/v1/posts`
**Description**: Create a new post  
**Access**: `posts.create` permission  
**Request Body**:
```json
{
  "title": "My New Blog Post",
  "content": "<p>Rich HTML content...</p>",
  "excerpt": "Brief description of the post",
  "featured_image_url": "https://s3.amazonaws.com/bucket/image.jpg",
  "category_id": "uuid",
  "tag_ids": ["uuid1", "uuid2"],
  "status": "draft"
}
```
**Response**: `201 Created`

#### **GET** `/api/v1/posts/:postId`
**Description**: Get post by ID (includes drafts for authors)  
**Access**: Public for published, `posts.view_own` for own drafts, `posts.view_all` for all  
**Response**: `200 OK`

#### **PUT** `/api/v1/posts/:postId`
**Description**: Update post  
**Access**: `posts.update` permission + ownership or `posts.update_all`  
**Response**: `200 OK`

#### **DELETE** `/api/v1/posts/:postId`
**Description**: Delete post  
**Access**: `posts.delete` permission + ownership or `posts.delete_all`  
**Response**: `200 OK`

#### **PATCH** `/api/v1/posts/:postId/publish`
**Description**: Publish a draft post  
**Access**: `posts.publish` permission  
**Response**: `200 OK`

#### **PATCH** `/api/v1/posts/:postId/unpublish`
**Description**: Unpublish a post (revert to draft)  
**Access**: `posts.publish` permission  
**Response**: `200 OK`

---

## **File Upload Management**

### **Image Upload**

#### **POST** `/api/v1/upload/image`
**Description**: Upload image to S3  
**Access**: Authenticated users  
**Request**: `multipart/form-data`
- `image`: Image file (jpg, png, webp)
- `type`: "avatar" | "post_featured" | "post_content"

**Response**: `201 Created`
```json
{
  "data": {
    "id": "uuid",
    "filename": "image-uuid.jpg",
    "original_name": "my-image.jpg",
    "url": "https://s3.amazonaws.com/bucket/images/image-uuid.jpg",
    "size": 1024000,
    "mime_type": "image/jpeg",
    "uploaded_by": "uuid",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### **DELETE** `/api/v1/upload/:fileId`
**Description**: Delete uploaded file  
**Access**: File owner or admin  
**Response**: `200 OK`

#### **GET** `/api/v1/uploads`
**Description**: Get user's uploaded files  
**Access**: Authenticated users  
**Query Parameters**: `?type=post_featured&page=1&limit=20`  
**Response**: `200 OK`

---

## **Search & Analytics**

### **Search**

#### **GET** `/api/v1/search`
**Description**: Global search across posts, categories, tags  
**Access**: Public  
**Query Parameters**:
- `?q=search+terms` - Search query (required)
- `?type=posts,categories,tags` - Search types
- `?limit=20` - Results limit

**Response**: `200 OK`
```json
{
  "data": {
    "posts": [
      {
        "id": "uuid",
        "title": "Highlighted title with search terms",
        "excerpt": "Highlighted excerpt...",
        "slug": "post-slug",
        "relevance_score": 0.95
      }
    ],
    "categories": [ /* matching categories */ ],
    "tags": [ /* matching tags */ ]
  },
  "query": "search terms",
  "total_results": 25,
  "search_time_ms": 45
}
```

### **Analytics (Admin)**

#### **GET** `/api/v1/analytics/overview`
**Description**: Dashboard analytics overview  
**Access**: `system.admin` permission  
**Response**: `200 OK`
```json
{
  "data": {
    "total_posts": 150,
    "total_published": 120,
    "total_drafts": 30,
    "total_users": 25,
    "total_authors": 8,
    "posts_this_month": 15,
    "new_users_this_month": 3
  }
}
```

#### **GET** `/api/v1/analytics/posts`
**Description**: Post analytics and statistics  
**Access**: `posts.view_all` permission  
**Query Parameters**: `?period=30d&author=uuid`  
**Response**: `200 OK`

---

## **System & Health**

### **Health Checks**

#### **GET** `/health`
**Description**: System health check  
**Access**: Public  
**Response**: `200 OK`
```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T00:00:00Z",
  "version": "1.0.0",
  "checks": {
    "database": "healthy",
    "s3": "healthy",
    "email_service": "healthy"
  }
}
```

#### **GET** `/health/detailed`
**Description**: Detailed system health  
**Access**: `system.admin` permission  
**Response**: `200 OK`

### **System Configuration**

#### **GET** `/api/v1/system/config`
**Description**: Get public system configuration  
**Access**: Public  
**Response**: `200 OK`
```json
{
  "data": {
    "site_name": "My Blog",
    "max_upload_size": 10485760,
    "allowed_image_types": ["jpg", "jpeg", "png", "webp"],
    "pagination_default": 10,
    "pagination_max": 100
  }
}
```

---

## **Error Responses**

### **Standard Error Format**
```json
{
  "error": "Error message",
  "code": "ERROR_CODE", 
  "details": {
    "field": "additional context"
  },
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### **Common HTTP Status Codes**
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (email exists, etc.)
- `422 Unprocessable Entity` - Validation errors
- `500 Internal Server Error` - Server error

### **Common Error Codes**
- `VALIDATION_FAILED` - Input validation errors
- `EMAIL_EXISTS` - Email already registered
- `USERNAME_EXISTS` - Username already taken
- `INVALID_CREDENTIALS` - Login failed
- `TOKEN_EXPIRED` - JWT token expired
- `INSUFFICIENT_PERMISSIONS` - Missing required permission
- `RESOURCE_NOT_FOUND` - Requested resource doesn't exist
- `RATE_LIMIT_EXCEEDED` - Too many requests

---

## **Authentication & Authorization**

### **Authentication Methods**
- **JWT Tokens**: Access (15 min) + Refresh (7 days) tokens
- **HTTP-Only Cookies**: Secure token storage
- **Bearer Header**: Alternative to cookies for API clients

### **Permission System**
- **Resource-based**: `posts.create`, `users.manage`
- **Ownership**: Users can edit their own content
- **Role-based**: Admin, Author roles with different permissions

### **Rate Limiting**
- **Authentication**: 5 attempts per 15 minutes per IP
- **API calls**: 1000 requests per hour per user
- **Upload**: 10 uploads per minute per user