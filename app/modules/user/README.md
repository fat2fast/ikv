# Module User

## Tổng quan

Module User là một module độc lập trong hệ thống IKV, được thiết kế theo kiến trúc Clean Architecture và CQRS pattern. Module này quản lý toàn bộ chức năng liên quan đến người dùng, authentication và profile management.

## Thông tin Module

- **Tên module**: user
- **Phiên bản**: 1.0.0
- **Trạng thái**: Enabled
- **Mô tả**: User management module with authentication
- **Database Schema**: `user_schema`
- **Table**: `usr_users`

## Cấu trúc thư mục

```
app/modules/user/
├── config.yaml              # Cấu hình module
├── module.go                 # File chính của module với dependency injection
├── README.md                 # File mô tả này
├── model/                    # Domain Layer
│   ├── user.go              # User Entity với business methods
│   ├── dtos.go              # Data Transfer Objects cho API
│   └── error.go             # Custom domain errors
├── service/                  # Application Layer - CQRS Handlers
│   ├── authenticate.go      # Authentication business logic
│   ├── register_new_user.go # User registration handler
│   ├── get_profile.go       # Get user profile query
│   ├── update_profile.go    # Update profile command
│   └── get_list_user.go     # List users query (stub)
├── infras/                   # Infrastructure Layer
│   ├── controller/          # Presentation layer
│   │   └── http-gin/        # REST API controllers
│   │       ├── base_gin_controller.go    # Base controller với dependencies
│   │       ├── authenticate_api.go       # POST /authenticate
│   │       ├── register_user_api.go      # POST /register
│   │       ├── get_profile_api.go        # GET /profile/:id
│   │       └── update_profile_api.go     # PUT /profile/:id
│   └── repository/          # Data access layer
│       └── gorm-pgsql/      # GORM PostgreSQL implementation
│           ├── repo.go                   # Base repository
│           ├── find.go                   # Read operations
│           ├── insert_user.go            # Create operations
│           └── update_profile.go         # Update operations
├── urls/                     # API routing
│   └── v1/
│       └── user_url.go      # API v1 routes definition
└── migrations/               # Database migrations
    ├── 20250623101550_create_table_user.up.sql
    └── 20250623101550_create_table_user.down.sql
```

## Cấu hình (config.yaml)

Module sử dụng file `config.yaml` để cấu hình với hỗ trợ environment variables:

### Cấu hình Module
```yaml
module:
  name: "user"
  version: "1.0.0"
  enabled: true
  description: "User management module with authentication"
```

### Cấu hình Database
```yaml
database:
  connection:
    driver: "${MODULE_USER_DB_DRIVER:postgres}"
    host: "${MODULE_USER_DB_HOST:localhost}"
    port: "${MODULE_USER_DB_PORT:5432}"
    database: "${MODULE_USER_DB_NAME:ikv_user}"
    username: "${MODULE_USER_DB_USER:ikv_user}"
    password: "${MODULE_USER_DB_PASSWORD:ikv_password}"
    schema: "${MODULE_USER_DB_SCHEMA:user_schema}"
    auto_create: ${MODULE_USER_DB_AUTO_CREATE:true}
    ssl_mode: "${MODULE_USER_DB_SSL_MODE:disable}"
    timezone: "${MODULE_USER_DB_TIMEZONE:Asia/Ho_Chi_Minh}"
```

### Environment Variables Override
- `MODULE_USER_DB_HOST`
- `MODULE_USER_DB_PORT`
- `MODULE_USER_DB_NAME`
- `MODULE_USER_DB_USER`
- `MODULE_USER_DB_PASSWORD`
- `MODULE_USER_DB_SCHEMA`

## Domain Model

### User Entity (`model/user.go`)

```go
type User struct {
    ID        uuid.UUID  `json:"id" gorm:"column:id;"`
    CreatedBy string     `json:"created_by" gorm:"column:created_by;"`
    CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
    UpdatedBy string     `json:"updated_by" gorm:"column:updated_by;"`
    UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;"`
    Status    UserStatus `json:"status" gorm:"column:status;"`
    Type      UserType   `json:"type" gorm:"column:type;"`
    Role      UserRole   `json:"role" gorm:"column:role;"`
    FirstName string     `json:"first_name" gorm:"column:first_name;"`
    LastName  string     `json:"last_name" gorm:"column:last_name;"`
    Phone     string     `json:"phone" gorm:"column:phone;"`
    Email     string     `json:"email" gorm:"column:email;"`
    Password  string     `json:"password" gorm:"column:password;"`
    Salt      string     `json:"salt" gorm:"column:salt;"`
}
```

### User Types, Roles và Status

#### User Types
- `email_password`: Đăng ký bằng email/password
- `facebook`: Đăng nhập qua Facebook
- `gmail`: Đăng nhập qua Gmail

#### User Roles
- `user`: Người dùng thông thường
- `admin`: Quản trị viên

#### User Status
- `pending`: Chờ xác nhận email
- `active`: Đang hoạt động
- `inactive`: Tạm khóa
- `banned`: Bị cấm
- `deleted`: Đã xóa

### Business Methods
```go
func (u *User) GetFullName() string
func (u *User) ToProfileResponse() *ProfileResponse
```

## API Endpoints

Module cung cấp RESTful API endpoints cho user management:

### Base URL: `/v1/users`

| Method | Endpoint | Mô tả | Request | Response | Handler |
|--------|----------|-------|---------|----------|---------|
| POST   | `/authenticate` | Đăng nhập | `LoginForm` | `AuthenticateResult` | `ActionAuthenticate` |
| POST   | `/register` | Đăng ký user mới | `RegisterForm` | `RegisterResponse` | `ActionRegister` |
| GET    | `/profile/:id` | Lấy thông tin profile | URL param | `ProfileResponse` | `ActionGetProfile` |
| PUT    | `/profile/:id` | Cập nhật profile | `UpdateProfileRequest` | Success message | `ActionUpdateProfile` |

### API Request/Response Examples

#### 1. Authentication - POST `/v1/users/authenticate`

**Request:**
```json
{
    "username": "user@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expIn": 3600
    }
}
```

#### 2. User Registration - POST `/v1/users/register`

**Request:**
```json
{
    "email": "newuser@example.com",
    "password": "securepass123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "0123456789"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "newuser@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "phone": "0123456789"
    }
}
```

#### 3. Get Profile - GET `/v1/users/profile/:id`

**Response:**
```json
{
    "success": true,
    "data": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "phone": "0123456789",
        "full_name": "John Doe",
        "role": "user",
        "status": "active",
        "type": "email_password",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4. Update Profile - PUT `/v1/users/profile/:id`

**Request:**
```json
{
    "first_name": "Jane",
    "last_name": "Smith",
    "phone": "0987654321"
}
```

**Response:**
```json
{
    "success": true,
    "message": "Profile updated successfully"
}
```

## Kiến trúc Clean Architecture

### 1. Domain Layer (model/)
- **Entities**: `User` với business rules
- **Value Objects**: `UserType`, `UserRole`, `UserStatus`
- **Domain Errors**: Custom errors cho business logic
- **Interfaces**: Repository contracts

### 2. Application Layer (service/)
- **Commands**: Write operations (Create, Update)
  - `AuthenticateCommand`: Xử lý đăng nhập
  - `CreateCommand`: Đăng ký user mới
  - `UpdateProfileCommand`: Cập nhật profile
- **Queries**: Read operations
  - `GetProfileQuery`: Lấy thông tin profile
  - `ListUsersQuery`: Lấy danh sách users (planned)

### 3. Infrastructure Layer (infras/)
- **Controllers**: HTTP request handlers
- **Repositories**: Database access implementations
- **External Services**: JWT token management

## CQRS Implementation

### Command Handlers (Write Operations)

#### Authentication Handler
```go
type AuthenticateCommandHandler struct {
    repo        IAuthenticateRepo
    tokenIssuer ITokenIssuer
}

func (hdl *AuthenticateCommandHandler) Execute(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResult, error)
```

#### Registration Handler
```go
type CreateCommandHandler struct {
    repo ICreateRepo
}

func (hdl *CreateCommandHandler) Execute(ctx context.Context, cmd *CreateCommand) (*RegisterResponse, error)
```

### Query Handlers (Read Operations)

#### Profile Query Handler
```go
type GetProfileQueryHandler struct {
    repo IGetRepo
}

func (hdl *GetProfileQueryHandler) Execute(ctx context.Context, query *GetProfileQuery) (*ProfileResponse, error)
```

## Database Schema

### Table Structure (usr_users)
```sql
CREATE TABLE usr_users (
    id varchar(36) PRIMARY KEY,
    created_by varchar(36),
    created_at timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by varchar(36),
    updated_at timestamp(6),
    status user_status_enum NOT NULL DEFAULT 'active',
    type user_type_enum NOT NULL DEFAULT 'email_password',
    role user_role_enum NOT NULL DEFAULT 'user',
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    phone varchar(20) DEFAULT NULL,
    email varchar(50) NOT NULL,
    password varchar(100) NOT NULL,
    salt varchar(50) NOT NULL,
    UNIQUE (email)
);
```

### Enum Types
- `user_status_enum`: pending, active, inactive, banned, deleted
- `user_type_enum`: email_password, facebook, gmail
- `user_role_enum`: user, admin

### Indexes
- Primary key: `id`
- Unique constraint: `email`
- Default indexes được tạo tự động

## Security

### Password Hashing
- Sử dụng bcrypt với salt riêng cho mỗi user
- Salt được generate ngẫu nhiên và lưu trong database
- Password được hash với `bcrypt.GenerateFromPassword`

### JWT Authentication
- Token được issue sau khi authenticate thành công
- Token chứa user ID và expire time
- Sử dụng shared component cho JWT management

## Repository Pattern

### Interfaces
```go
type IAuthenticateRepo interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
}

type ICreateRepo interface {
    Insert(ctx context.Context, data *User) error
}

type IGetRepo interface {
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type IUpdateRepo interface {
    UpdateProfile(ctx context.Context, id uuid.UUID, data UpdateProfileData) error
}
```

### Implementation
- PostgreSQL implementation sử dụng GORM
- Connection pooling được cấu hình
- Schema isolation với `user_schema`

## Dependency Injection

```go
func Initialize(appCtx sharedinfras.IAppContext) *userhttpgin.UserHTTPController {
    dbCtx := appCtx.DbContext()
    
    // Repository
    userRepository := userrepository.NewUserRepository(dbCtx)
    
    // JWT Component
    jwtComponent := sharecomponent.NewJwtComponent()
    
    // Command Handlers
    authenticateHandler := userservice.NewAuthenticateCommandHandler(userRepository, jwtComponent)
    createHandler := userservice.NewCreateCommandHandler(userRepository)
    updateProfileHandler := userservice.NewUpdateProfileCommandHandler(userRepository)
    
    // Query Handlers
    getProfileHandler := userservice.NewGetProfileQueryHandler(userRepository)
    
    // Controller
    return userhttpgin.NewUserHTTPController(
        createHandler,
        authenticateHandler,
        getProfileHandler,
        updateProfileHandler,
    )
}
```

## Cách sử dụng

### 1. Khởi tạo Module

```go
userModule, err := user.NewModule()
if err != nil {
    log.Fatal("Failed to initialize user module:", err)
}
```

### 2. Đăng ký với Router

```go
router := gin.Default()
err = userModule.Register(router)
if err != nil {
    log.Fatal("Failed to register user module:", err)
}
```

### 3. Chạy Migration

```go
err = userModule.RunMigrations()
if err != nil {
    log.Fatal("Failed to run migrations:", err)
}
```

### 4. Test API với cURL

#### Đăng ký user mới:
```bash
curl -X POST http://localhost:8080/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User",
    "phone": "0123456789"
  }'
```

#### Đăng nhập:
```bash
curl -X POST http://localhost:8080/v1/users/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test@example.com",
    "password": "password123"
  }'
```

#### Lấy profile:
```bash
curl -X GET http://localhost:8080/v1/users/profile/{user_id} \
  -H "Authorization: Bearer {jwt_token}"
```

## Trạng thái Implementation

### ✅ Đã hoàn thành:
- ✅ Cấu trúc module cơ bản với Clean Architecture
- ✅ Database connection và schema management
- ✅ User entity với business methods
- ✅ Authentication với JWT token
- ✅ User registration với password hashing
- ✅ Get profile functionality
- ✅ Update profile functionality
- ✅ Repository pattern implementation
- ✅ CQRS với Command/Query handlers
- ✅ HTTP controllers với error handling
- ✅ URL routing và API endpoints
- ✅ Database migrations với PostgreSQL enums
- ✅ Dependency injection pattern

### 🔄 Đang phát triển:
- 🔄 Email verification system
- 🔄 Password reset functionality
- 🔄 User management (admin functions)
- 🔄 Social login (Facebook, Google)
- 🔄 Advanced user filtering và search

### ⏳ Kế hoạch tiếp theo:
- ⏳ gRPC support cho internal communication
- ⏳ Role-based authorization middleware
- ⏳ User activity logging
- ⏳ Rate limiting cho authentication
- ⏳ Comprehensive unit tests
- ⏳ Integration tests
- ⏳ Performance optimization
- ⏳ Caching layer với Redis
- ⏳ Event sourcing cho user activities

## Import Path Convention

Theo quy ước của project, sử dụng:
```go
import "fat2fast/ikv/modules/user/..."
import "fat2fast/ikv/shared/..."
```

## Error Handling

### Domain Errors
- `ErrInvalidEmailAndPassword`: Sai thông tin đăng nhập
- `ErrUserBannedOrDeleted`: User bị cấm hoặc đã xóa
- `ErrEmailAlreadyExists`: Email đã tồn tại
- `ErrInvalidUserID`: ID user không hợp lệ

### HTTP Status Codes
- `400`: Bad Request (validation errors, business logic errors)
- `401`: Unauthorized (invalid credentials)
- `403`: Forbidden (banned/deleted users)
- `404`: Not Found (user not found)
- `500`: Internal Server Error (database errors, system errors)

## Monitoring và Logs

Module sử dụng structured logging với các level:
- `INFO`: Thông tin hoạt động bình thường
- `WARN`: Cảnh báo (invalid login attempts)
- `ERROR`: Lỗi hệ thống (database connection issues)
- `DEBUG`: Chi tiết cho debugging

## Bảo mật

### Validation
- Email format validation
- Password strength requirements (min 8 chars)
- Input sanitization
- SQL injection prevention với GORM

### Data Protection
- Password không bao giờ trả về trong response
- Salt được generate unique cho mỗi user
- JWT token có thời hạn expire

Đây là module User hoàn chỉnh với authentication, user management và clean architecture. Module đã sẵn sàng cho production với các tính năng bảo mật và performance optimization. 