# Production-Ready Blog Application - Development Roadmap

## Project Overview
Building a multi-author blog platform using Go + Fiber + PostgreSQL + sqlc + go-migrate with JWT + cookie authentication, flexible role-permission system, and production-ready features.

## Technology Stack & Recommended Libraries

### Core Framework & Database
- **Web Framework**: `github.com/gofiber/fiber/v2`
- **Database**: PostgreSQL with Docker
- **Database Migrations**: `github.com/golang-migrate/migrate/v4`
- **Query Builder**: `github.com/sqlc-dev/sqlc`
- **Database Driver**: `github.com/lib/pq`

### Authentication & Security
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **Password Hashing**: `golang.org/x/crypto/bcrypt`
- **Validation**: `github.com/go-playground/validator/v10`
- **UUID**: `github.com/google/uuid`

### Configuration & Environment
- **Config Management**: `github.com/spf13/viper`
- **Environment**: `github.com/joho/godotenv`

### File Upload & Email (Later phases)
- **AWS SDK**: `github.com/aws/aws-sdk-go-v2`
- **Email Service**: Consider Resend, SendGrid, or Mailgun

### Development & Documentation
- **Air (Hot Reload)**: `github.com/cosmtrek/air`
- **Swagger**: `github.com/swaggo/fiber-swagger`

## Recommended Folder Structure

```
blog-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   └── health_handler.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   └── jwt_service.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   ├── role_repository.go
│   │   └── permission_repository.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   └── logger_middleware.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── role.go
│   │   └── response.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrations/
│   └── utils/
│       ├── validator.go
│       ├── response.go
│       └── password.go
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_roles_permissions.up.sql
│   └── ...
├── sqlc/
│   ├── sqlc.yaml
│   ├── queries/
│   │   ├── users.sql
│   │   ├── posts.sql
│   │   └── roles.sql
│   └── schema.sql
├── docker/
│   └── docker-compose.yml
├── docs/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Database Schema Design

### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Roles table (flexible system)
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL, -- 'admin', 'author', 'editor'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL, -- 'posts.create', 'posts.delete', 'users.manage'
    resource VARCHAR(50) NOT NULL, -- 'posts', 'users', 'comments'
    action VARCHAR(50) NOT NULL, -- 'create', 'read', 'update', 'delete'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Role-Permission junction table
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- User-Role junction table
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);

-- Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7), -- Hex color
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tags table
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Posts table
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    featured_image_url TEXT,
    status VARCHAR(20) DEFAULT 'draft', -- 'draft', 'published', 'archived'
    author_id UUID REFERENCES users(id) ON DELETE CASCADE,
    -- category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Post-Tag junction table
CREATE TABLE post_tags (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE post_categories (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);
```

---

# DEVELOPMENT ROADMAP

## **EPIC 1: Project Foundation & Development Environment**
**Goal**: Establish a solid, maintainable codebase foundation

### **Story 1.1**: Project Structure Setup
**Value**: Clean, scalable codebase organization  
**Scope**: 
- Create folder structure as outlined above
- Initialize Go modules
- Setup basic Makefile with common commands
- Configure Air for hot reloading

**Technical Details**:
- Use `go mod init blog-app`
- Create all necessary directories
- Setup basic `.gitignore` for Go projects
- Create `.env.example` with all required environment variables

### **Story 1.2**: Environment Configuration Management
**Value**: Easy environment switching (dev/staging/prod)  
**Scope**: 
- Config struct with Viper
- Environment validation
- Database connection configuration
- JWT secrets management

**Technical Details**:
- Use Viper for config loading from `.env` files
- Implement config validation with struct tags
- Support for different environments (dev, staging, prod)
- Secure handling of sensitive data (JWT secrets, database passwords)

### **Story 1.3**: Database Connection & Health Check
**Value**: Reliable database connectivity with monitoring  
**Scope**: 
- PostgreSQL connection with connection pooling
- Docker Compose setup for local development
- Database health check endpoint
- Migration runner setup

**Technical Details**:
- Docker Compose with PostgreSQL and pgAdmin
- Connection pooling configuration
- Health check endpoint (`/health`) that verifies DB connectivity
- Integration with go-migrate for running migrations

### **Story 1.4**: Basic Fiber Server with Middleware Stack
**Value**: Production-ready server foundation  
**Scope**: 
- Fiber server setup with graceful shutdown
- Essential middleware (CORS, logging, recovery, rate limiting)
- Request/response logging
- Error handling middleware

**Technical Details**:
- Structured logging with configurable levels
- CORS configuration for frontend integration
- Panic recovery with proper error responses
- Rate limiting middleware for API protection
- Custom error handler for consistent error responses

---

## **EPIC 2: User Management & Authentication Foundation**
**Goal**: Secure user registration and authentication system

### **Story 2.1**: User Database Schema & Repository
**Value**: Persistent user data storage  
**Scope**: 
- Users table migration
- sqlc queries for user CRUD operations
- User repository implementation
- Password hashing utilities

**Technical Details**:
- Create migration files for users table
- Define sqlc queries in `queries/users.sql`
- Implement UserRepository interface with methods: Create, GetByID, GetByEmail, Update, Delete
- Bcrypt password hashing with proper cost factor (12)

### **Story 2.2**: User Registration Endpoint
**Value**: New users can join the platform  
**Scope**: 
- Registration request validation
- Email uniqueness check
- Password strength validation
- User creation with default role

**Technical Details**:
- Input validation using validator package
- Email format and uniqueness validation
- Password complexity requirements (min 8 chars, mix of chars)
- Automatic assignment of default "author" role
- Return sanitized user data (no password hash)

### **Story 2.3**: JWT Authentication Service
**Value**: Secure token-based authentication  
**Scope**: 
- JWT token generation and validation
- Access and refresh token logic
- Token expiration handling
- Secure token storage in HTTP-only cookies

**Technical Details**:
- Use RS256 or HS256 algorithm for JWT signing
- Access tokens (15 minutes expiry) and refresh tokens (7 days)
- JWT claims: user_id, email, roles, permissions
- Secure cookie configuration (HttpOnly, Secure, SameSite)

### **Story 2.4**: Login Endpoint with Cookie-based JWT
**Value**: Users can authenticate and stay logged in  
**Scope**: 
- Login credential validation
- JWT token generation
- Cookie setting with security flags
- Login attempt logging

**Technical Details**:
- Email/password validation
- Password verification with bcrypt
- Set access token in HTTP-only cookie
- Set refresh token in separate HTTP-only cookie
- Return user profile data (excluding sensitive info)

### **Story 2.5**: Authentication Middleware
**Value**: Protected routes and user context  
**Scope**: 
- JWT validation from cookies
- User context injection
- Token refresh logic
- Logout endpoint

**Technical Details**:
- Extract JWT from HTTP-only cookies
- Validate token signature and expiration
- Load user data and inject into request context
- Automatic token refresh when access token expires
- Clear cookies on logout

---

## **EPIC 3: Role-Based Access Control System**
**Goal**: Flexible permission management for different user types

### **Story 3.1**: Roles & Permissions Database Schema
**Value**: Flexible, extensible permission management  
**Scope**: 
- Roles, permissions, and junction tables
- Seed data for default roles and permissions
- Migration for role-permission system

**Technical Details**:
- Create migrations for roles, permissions, role_permissions, user_roles tables
- Seed default roles: admin, author
- Seed default permissions: posts.create, posts.update, posts.delete, users.manage, etc.
- Assign permissions to default roles

### **Story 3.2**: Permission Service & Repository
**Value**: Efficient permission checking and role management  
**Scope**: 
- Role and permission repositories
- Permission checking service
- User permission loading

**Technical Details**:
- RoleRepository: CRUD operations for roles
- PermissionRepository: Permission querying and assignment
- Service methods: GetUserPermissions, HasPermission, AssignRole
- Efficient permission caching strategies

### **Story 3.3**: Authorization Middleware
**Value**: Route-level permission control  
**Scope**: 
- Permission-based route protection
- Role-based route guards
- Admin-only endpoints protection

**Technical Details**:
- Middleware factory: RequirePermission("posts.create")
- Role-based middleware: RequireRole("admin")
- Integration with existing auth middleware
- Clear error messages for unauthorized access

### **Story 3.4**: User Profile & Role Management Endpoints
**Value**: Users can manage profiles, admins can assign roles  
**Scope**: 
- User profile CRUD endpoints
- Role assignment endpoints (admin only)
- User role viewing

**Technical Details**:
- GET/PUT `/api/profile` - user profile management
- GET `/api/users/:id/roles` - view user roles
- POST/DELETE `/api/users/:id/roles` - assign/remove roles (admin only)
- Input validation and authorization checks

---

## **EPIC 4: Blog Content Management System**
**Goal**: Complete blog post creation, management, and organization

### **Story 4.1**: Categories & Tags System
**Value**: Content organization and discovery  
**Scope**: 
- Categories and tags database schema
- CRUD operations for categories and tags
- Category/tag management endpoints

**Technical Details**:
- Create migrations for categories and tags tables
- Slug generation for SEO-friendly URLs
- CategoryRepository and TagRepository implementation
- Admin endpoints for managing categories and tags

### **Story 4.2**: Blog Posts Core CRUD
**Value**: Authors can create and manage blog posts  
**Scope**: 
- Posts table and relationships
- Post creation, editing, deletion
- Draft and published status management
- Author association

**Technical Details**:
- Posts table migration with foreign keys
- PostRepository with complex queries (joins with categories, tags, authors)
- Post status workflow: draft → published → archived
- Slug generation and uniqueness validation
- Rich text content storage

### **Story 4.3**: Post-Tag Association System
**Value**: Flexible post tagging for better content discovery  
**Scope**: 
- Many-to-many relationship between posts and tags
- Tag assignment/removal for posts
- Tag-based post querying

**Technical Details**:
- post_tags junction table migration
- Repository methods for tag assignment
- Efficient tag querying with post counts
- Bulk tag operations

### **Story 4.4**: Advanced Post Querying
**Value**: Powerful content filtering and search capabilities  
**Scope**: 
- Post listing with filters (category, tags, author, status)
- Search functionality (title, content, excerpt)
- Pagination implementation
- Sorting options

**Technical Details**:
- Complex SQL queries with sqlc
- Full-text search implementation
- Pagination with cursor-based or offset-based approach
- Query optimization with proper indexing
- Filter combination logic

---

## **EPIC 5: Enhanced Blog Features**
**Goal**: Advanced functionality for better user experience

### **Story 5.1**: Image Upload System
**Value**: Rich content creation with image support  
**Scope**: 
- AWS S3 integration
- Image upload endpoints
- Image processing and optimization
- Featured image support for posts

**Technical Details**:
- AWS SDK v2 configuration
- Image validation (format, size limits)
- Unique filename generation
- S3 bucket configuration with proper permissions
- Integration with post creation/editing

### **Story 5.2**: Email Notification System
**Value**: User engagement through email notifications  
**Scope**: 
- Email service integration (Resend/SendGrid)
- Welcome emails for new users
- Email templates
- Notification preferences

**Technical Details**:
- Email service abstraction layer
- HTML email templates
- Background job processing for emails
- User email preferences management
- Email delivery tracking

### **Story 5.3**: Advanced Search & Filtering
**Value**: Enhanced content discoverability  
**Scope**: 
- Full-text search implementation
- Advanced filtering combinations
- Search result ranking
- Search analytics

**Technical Details**:
- PostgreSQL full-text search or Elasticsearch integration
- Search indexing strategies
- Result relevance scoring
- Search query optimization
- Search result caching

---

## **EPIC 6: API Documentation & Production Readiness**
**Goal**: Complete, documented, and production-ready API

### **Story 6.1**: API Documentation
**Value**: Clear API documentation for frontend developers  
**Scope**: 
- Swagger/OpenAPI documentation
- API endpoint documentation
- Request/response examples
- Authentication documentation

### **Story 6.2**: Production Optimization
**Value**: Performance and reliability improvements  
**Scope**: 
- Response caching strategies
- Database query optimization
- API rate limiting
- Performance monitoring

### **Story 6.3**: Security Hardening
**Value**: Enhanced security for production deployment  
**Scope**: 
- Security headers middleware
- Input sanitization
- SQL injection prevention
- CORS configuration refinement

---

## Development Phases Priority

1. **Phase 1**: Stories 1.1 → 1.4 (Foundation)
2. **Phase 2**: Stories 2.1 → 2.5 (Authentication)
3. **Phase 3**: Stories 3.1 → 3.4 (Authorization)
4. **Phase 4**: Stories 4.1 → 4.4 (Blog Core)
5. **Phase 5**: Stories 5.1 → 5.3 (Enhanced Features)
6. **Phase 6**: Stories 6.1 → 6.3 (Production Readiness)

Each story is designed to deliver tangible value and can be developed incrementally. You can choose to implement them in order or adjust based on your priorities!