# Quy Chuẩn Module Development

## Tổng quan

Tài liệu này mô tả quy chuẩn để phát triển một module mới trong hệ thống IKV, dựa trên phân tích module User hiện tại. Mỗi module được thiết kế theo kiến trúc Clean Architecture và CQRS pattern để đảm bảo tính độc lập, khả năng mở rộng và bảo trì.

## Nguyên tắc thiết kế

### 1. Module Independence
- Mỗi module hoàn toàn độc lập về database, configuration và business logic
- Không phụ thuộc trực tiếp vào module khác
- Có thể enable/disable riêng biệt

### 2. Clean Architecture
- **Domain Layer**: Entities, business rules (không phụ thuộc gì)
- **Application Layer**: Use cases, business logic
- **Infrastructure Layer**: Controllers, repositories, external services

### 3. CQRS Pattern
- Tách biệt Command (write operations) và Query (read operations)
- Command handlers cho business logic phức tạp
- Query handlers cho read operations tối ưu

## Cấu trúc Module Chuẩn

```
app/modules/{module_name}/
├── module.go                 # Module chính - entry point
├── config.yaml              # Cấu hình module
├── README.md                 # Documentation module
├── model/                    # Domain Layer
│   ├── {entity}.go          # Domain entities
│   ├── dtos.go              # Data Transfer Objects
│   ├── interfaces.go        # Repository interfaces
│   └── error.go             # Custom errors
├── service/                  # Application Layer
│   ├── commands/            # Write operations
│   │   ├── create_{entity}.go
│   │   ├── update_{entity}.go
│   │   └── delete_{entity}.go
│   └── queries/             # Read operations
│       ├── get_{entity}_details.go
│       └── list_{entity}.go
├── infras/                   # Infrastructure Layer
│   ├── controller/          # Presentation layer
│   │   ├── http-gin/        # REST API controllers
│   │   │   ├── base_{entity}_controller.go    # Base controller với dependencies
│   │   │   ├── create_{entity}_api.go         # POST / - Tạo mới
│   │   │   ├── get_{entity}_detail_api.go     # GET /:id - Chi tiết
│   │   │   ├── list_{entity}_api.go           # GET / - Danh sách
│   │   │   ├── update_{entity}_api.go         # PUT /:id - Cập nhật
│   │   │   └── delete_{entity}_api.go         # DELETE /:id - Xóa
│   │   └── grpcctl/         # gRPC controllers
│   │       ├── base_{entity}_grpc_controller.go
│   │       ├── create_{entity}_grpc.go
│   │       ├── get_{entity}_detail_grpc.go
│   │       ├── list_{entity}_grpc.go
│   │       ├── update_{entity}_grpc.go
│   │       └── delete_{entity}_grpc.go
│   └── repository/          # Data access layer
│       └── gorm-pgsql/      # Database implementation
│           ├── repo.go                        # Base repository struct
│           ├── create_{entity}.go             # Insert operations
│           ├── get_{entity}.go                # Get by ID operations
│           ├── list_{entity}.go               # List/Search operations
│           ├── update_{entity}.go             # Update operations
│           └── delete_{entity}.go             # Delete operations
├── urls/                     # API routing
│   └── v1/
│       └── {entity}_url.go
└── migrations/               # Database migrations
    ├── 001_create_{entity}_table.up.sql
    └── 001_create_{entity}_table.down.sql
```

## Thành phần chi tiết

### 1. Module.go - Entry Point

```go
package {module_name}

import (
    "fmt"
    "log"
    "path/filepath"
    "runtime"
    "time"
    
    "{project_name}/ikv/modules/{module_name}/infras/controller/http-gin"
    "{project_name}/ikv/modules/{module_name}/infras/repository/gorm-pgsql"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/modules/{module_name}/urls/v1"
    "{project_name}/ikv/shared"
    "{project_name}/ikv/shared/infras"
    
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Config cấu hình module
type Config struct {
    Module struct {
        Name        string `yaml:"name"`
        Version     string `yaml:"version"`
        Description string `yaml:"description"`
        Enabled     bool   `yaml:"enabled"`
    } `yaml:"module"`
    
    Database struct {
        Connection struct {
            Driver     string `yaml:"driver"`
            Host       string `yaml:"host"`
            Port       string `yaml:"port"`
            Database   string `yaml:"database"`
            Username   string `yaml:"username"`
            Password   string `yaml:"password"`
            SSLMode    string `yaml:"ssl_mode"`
            Timezone   string `yaml:"timezone"`
            Schema     string `yaml:"schema"`
            AutoCreate bool   `yaml:"auto_create"`
        } `yaml:"connection"`
        
        Migration struct {
            Path   string `yaml:"path"`
            Table  string `yaml:"table"`
            Schema string `yaml:"schema"`
        } `yaml:"migration"`
        
        Performance struct {
            MaxOpenConns    int    `yaml:"max_open_conns"`
            MaxIdleConns    int    `yaml:"max_idle_conns"`
            ConnMaxLifetime string `yaml:"conn_max_lifetime"`
        } `yaml:"performance"`
    } `yaml:"database"`
}

// Module đại diện cho module
type Module struct {
    config Config
    DB     *gorm.DB
}

// Interface cho module
type IModule interface {
    GetName() string
    IsEnabled() bool
    Register(router *gin.Engine) error
    RunMigrations() error
    GetDB() *gorm.DB
}

// NewModule khởi tạo module mới
func NewModule() (*Module, error) {
    // Load config và khởi tạo database connection
    // Implementation tương tự module User
}

// Register đăng ký module với router
func (m *Module) Register(router *gin.Engine) error {
    if !m.IsEnabled() {
        return nil
    }
    
    // Dependency injection
    controller := m.Initialize()
    routes := v1.GetRoutes(controller)
    
    // Đăng ký routes
    apiV1 := router.Group("/v1")
    moduleGroup := apiV1.Group("/{module_path}")
    
    for _, route := range routes {
        moduleGroup.Handle(route.Method, route.Path, route.HandlerFunc)
    }
    
    return nil
}

// Initialize dependency injection
func (m *Module) Initialize() *httpgin.{Entity}HTTPController {
    dbCtx := sharedinfras.NewDbContext(m.DB)
    
    // Repository
    repository := repository.New{Entity}Repository(dbCtx)
    
    // Command handlers
    createHandler := service.NewCreate{Entity}Handler(repository)
    updateHandler := service.NewUpdate{Entity}Handler(repository)
    deleteHandler := service.NewDelete{Entity}Handler(repository)
    
    // Query handlers
    getDetailsHandler := service.NewGet{Entity}DetailsHandler(repository)
    listHandler := service.NewList{Entity}Handler(repository)
    
    // Controller
    return httpgin.New{Entity}HTTPController(
        createHandler,
        updateHandler,
        deleteHandler,
        getDetailsHandler,
        listHandler,
    )
}
```

### 2. Config.yaml Template

```yaml
# Module Configuration
module:
  name: "{module_name}"
  version: "1.0.0"
  enabled: true
  description: "{Module description}"

# Database Configuration
database:
  connection:
    driver: "${MODULE_{MODULE_NAME}_DB_DRIVER:postgres}"
    host: "${MODULE_{MODULE_NAME}_DB_HOST:localhost}"
    port: "${MODULE_{MODULE_NAME}_DB_PORT:5432}"
    database: "${MODULE_{MODULE_NAME}_DB_NAME:ikv_{module_name}}"
    username: "${MODULE_{MODULE_NAME}_DB_USER:ikv_{module_name}}"
    password: "${MODULE_{MODULE_NAME}_DB_PASSWORD:ikv_password}"
    schema: "${MODULE_{MODULE_NAME}_DB_SCHEMA:{module_name}_schema}"
    auto_create: ${MODULE_{MODULE_NAME}_DB_AUTO_CREATE:true}
    ssl_mode: "${MODULE_{MODULE_NAME}_DB_SSL_MODE:disable}"
    timezone: "${MODULE_{MODULE_NAME}_DB_TIMEZONE:Asia/Ho_Chi_Minh}"

  migration:
    path: "${MODULE_{MODULE_NAME}_MIGRATION_PATH:/app/modules/{module_name}/migrations}"
    table: "${MODULE_{MODULE_NAME}_MIGRATION_TABLE:{module_name}_migrations}"
    schema: "${MODULE_{MODULE_NAME}_MIGRATION_SCHEMA:public}"

  performance:
    max_open_conns: ${MODULE_{MODULE_NAME}_DB_MAX_OPEN_CONNS:10}
    max_idle_conns: ${MODULE_{MODULE_NAME}_DB_MAX_IDLE_CONNS:2}
    conn_max_lifetime: "${MODULE_{MODULE_NAME}_DB_CONN_MAX_LIFETIME:5m}"
```

## Design Patterns và Interfaces

### 1. Repository Pattern

#### Interface Definition
```go
// model/interfaces.go
package {module_name}model

import (
    "context"
    "github.com/google/uuid"
)

// Repository interfaces tách biệt theo chức năng
type ICreate{Entity}Repository interface {
    Insert(ctx context.Context, entity *{Entity}) error
}

type IRead{Entity}Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*{Entity}, error)
    GetList(ctx context.Context, filter *ListFilter) ([]*{Entity}, int64, error)
}

type IUpdate{Entity}Repository interface {
    Update(ctx context.Context, id uuid.UUID, entity *{Entity}) error
}

type IDelete{Entity}Repository interface {
    Delete(ctx context.Context, id uuid.UUID) error
    SoftDelete(ctx context.Context, id uuid.UUID) error
}

// Composite interface cho full CRUD
type I{Entity}Repository interface {
    ICreate{Entity}Repository
    IRead{Entity}Repository
    IUpdate{Entity}Repository
    IDelete{Entity}Repository
}
```

#### Repository Implementation

##### Base Repository Structure
```go
// infras/repository/gorm-pgsql/repo.go
package repository

import (
    "{project_name}/ikv/modules/{module_name}/model"
    "{project_name}/ikv/shared/infras"
)

type {Entity}Repository struct {
    dbCtx infras.IDbContext
}

func New{Entity}Repository(dbCtx infras.IDbContext) {module_name}model.I{Entity}Repository {
    return &{Entity}Repository{dbCtx: dbCtx}
}

// GetDBContext trả về database context
func (r *{Entity}Repository) GetDBContext() infras.IDbContext {
    return r.dbCtx
}
```

##### Create Operations
```go
// infras/repository/gorm-pgsql/create_{entity}.go
package repository

import (
    "context"
    "{project_name}/ikv/modules/{module_name}/model"
)

// Insert tạo entity mới trong database
func (r *{Entity}Repository) Insert(ctx context.Context, entity *{module_name}model.{Entity}) error {
    db := r.dbCtx.GetDB()
    
    // Validate entity trước khi insert
    if err := r.validateEntity(entity); err != nil {
        return err
    }
    
    // Thực hiện insert
    if err := db.WithContext(ctx).Create(entity).Error; err != nil {
        return {module_name}model.NewDatabaseError("Failed to create {entity}", err)
    }
    
    return nil
}

// BatchInsert tạo nhiều entities cùng lúc
func (r *{Entity}Repository) BatchInsert(ctx context.Context, entities []*{module_name}model.{Entity}) error {
    if len(entities) == 0 {
        return nil
    }
    
    db := r.dbCtx.GetDB()
    
    // Validate tất cả entities
    for _, entity := range entities {
        if err := r.validateEntity(entity); err != nil {
            return err
        }
    }
    
    // Batch insert với transaction
    return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return tx.CreateInBatches(entities, 100).Error
    })
}

// validateEntity kiểm tra tính hợp lệ của entity
func (r *{Entity}Repository) validateEntity(entity *{module_name}model.{Entity}) error {
    if entity.Name == "" {
        return {module_name}model.NewValidationError("name", "Name is required")
    }
    if entity.Email == "" {
        return {module_name}model.NewValidationError("email", "Email is required")
    }
    return nil
}
```

##### Read Operations
```go
// infras/repository/gorm-pgsql/get_{entity}.go
package repository

import (
    "context"
    "errors"
    "{project_name}/ikv/modules/{module_name}/model"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// GetByID lấy entity theo ID
func (r *{Entity}Repository) GetByID(ctx context.Context, id uuid.UUID) (*{module_name}model.{Entity}, error) {
    db := r.dbCtx.GetDB()
    var entity {module_name}model.{Entity}
    
    err := db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, {module_name}model.Err{Entity}NotFound
        }
        return nil, {module_name}model.NewDatabaseError("Failed to get {entity}", err)
    }
    
    return &entity, nil
}

// GetByEmail lấy entity theo email
func (r *{Entity}Repository) GetByEmail(ctx context.Context, email string) (*{module_name}model.{Entity}, error) {
    db := r.dbCtx.GetDB()
    var entity {module_name}model.{Entity}
    
    err := db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&entity).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, {module_name}model.Err{Entity}NotFound
        }
        return nil, {module_name}model.NewDatabaseError("Failed to get {entity} by email", err)
    }
    
    return &entity, nil
}

// Exists kiểm tra entity có tồn tại không
func (r *{Entity}Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
    db := r.dbCtx.GetDB()
    var count int64
    
    err := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
    if err != nil {
        return false, {module_name}model.NewDatabaseError("Failed to check {entity} existence", err)
    }
    
    return count > 0, nil
}

// EmailExists kiểm tra email đã tồn tại chưa
func (r *{Entity}Repository) EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
    db := r.dbCtx.GetDB()
    var count int64
    
    query := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("email = ? AND deleted_at IS NULL", email)
    
    if excludeID != nil {
        query = query.Where("id != ?", *excludeID)
    }
    
    err := query.Count(&count).Error
    if err != nil {
        return false, {module_name}model.NewDatabaseError("Failed to check email existence", err)
    }
    
    return count > 0, nil
}
```

##### List Operations
```go
// infras/repository/gorm-pgsql/list_{entity}.go
package repository

import (
    "context"
    "strings"
    "{project_name}/ikv/modules/{module_name}/model"
    "gorm.io/gorm"
)

// GetList lấy danh sách entities với filter và pagination
func (r *{Entity}Repository) GetList(ctx context.Context, filter *{module_name}model.ListFilter) ([]*{module_name}model.{Entity}, int64, error) {
    db := r.dbCtx.GetDB()
    var entities []*{module_name}model.{Entity}
    var total int64
    
    // Build base query
    query := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("deleted_at IS NULL")
    
    // Apply filters
    query = r.applyFilters(query, filter)
    
    // Count total records
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, {module_name}model.NewDatabaseError("Failed to count {entity}s", err)
    }
    
    // Apply pagination and sorting
    query = r.applyPaginationAndSorting(query, filter)
    
    // Execute query
    if err := query.Find(&entities).Error; err != nil {
        return nil, 0, {module_name}model.NewDatabaseError("Failed to get {entity} list", err)
    }
    
    return entities, total, nil
}

// GetActiveList lấy danh sách entities có status active
func (r *{Entity}Repository) GetActiveList(ctx context.Context, limit int) ([]*{module_name}model.{Entity}, error) {
    db := r.dbCtx.GetDB()
    var entities []*{module_name}model.{Entity}
    
    query := db.WithContext(ctx).Where("status = ? AND deleted_at IS NULL", {module_name}model.StatusActive).
        Order("created_at DESC")
    
    if limit > 0 {
        query = query.Limit(limit)
    }
    
    if err := query.Find(&entities).Error; err != nil {
        return nil, {module_name}model.NewDatabaseError("Failed to get active {entity}s", err)
    }
    
    return entities, nil
}

// applyFilters áp dụng các filter vào query
func (r *{Entity}Repository) applyFilters(query *gorm.DB, filter *{module_name}model.ListFilter) *gorm.DB {
    // Filter by status
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    // Search by name or email
    if filter.Search != "" {
        searchTerm := "%" + strings.ToLower(filter.Search) + "%"
        query = query.Where("(LOWER(name) LIKE ? OR LOWER(email) LIKE ?)", searchTerm, searchTerm)
    }
    
    // Filter by date range
    if !filter.CreatedFrom.IsZero() {
        query = query.Where("created_at >= ?", filter.CreatedFrom)
    }
    if !filter.CreatedTo.IsZero() {
        query = query.Where("created_at <= ?", filter.CreatedTo)
    }
    
    return query
}

// applyPaginationAndSorting áp dụng pagination và sorting
func (r *{Entity}Repository) applyPaginationAndSorting(query *gorm.DB, filter *{module_name}model.ListFilter) *gorm.DB {
    // Sorting
    sortBy := filter.SortBy
    if sortBy == "" {
        sortBy = "created_at"
    }
    
    sortOrder := filter.SortOrder
    if sortOrder == "" {
        sortOrder = "DESC"
    }
    
    query = query.Order(sortBy + " " + sortOrder)
    
    // Pagination
    if filter.Page > 0 && filter.PerPage > 0 {
        offset := (filter.Page - 1) * filter.PerPage
        query = query.Offset(offset).Limit(filter.PerPage)
    }
    
    return query
}
```

##### Update Operations
```go
// infras/repository/gorm-pgsql/update_{entity}.go
package repository

import (
    "context"
    "time"
    "{project_name}/ikv/modules/{module_name}/model"
    "github.com/google/uuid"
)

// Update cập nhật entity theo ID
func (r *{Entity}Repository) Update(ctx context.Context, id uuid.UUID, entity *{module_name}model.{Entity}) error {
    db := r.dbCtx.GetDB()
    
    // Validate entity
    if err := r.validateEntity(entity); err != nil {
        return err
    }
    
    // Set updated_at
    entity.UpdatedAt = time.Now()
    
    // Update entity
    result := db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Updates(entity)
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to update {entity}", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return {module_name}model.Err{Entity}NotFound
    }
    
    return nil
}

// UpdateFields cập nhật các fields cụ thể
func (r *{Entity}Repository) UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
    db := r.dbCtx.GetDB()
    
    // Add updated_at
    fields["updated_at"] = time.Now()
    
    // Update specific fields
    result := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id = ? AND deleted_at IS NULL", id).Updates(fields)
    
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to update {entity} fields", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return {module_name}model.Err{Entity}NotFound
    }
    
    return nil
}

// UpdateStatus cập nhật status của entity
func (r *{Entity}Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
    return r.UpdateFields(ctx, id, map[string]interface{}{
        "status": status,
    })
}

// BulkUpdate cập nhật nhiều entities cùng lúc
func (r *{Entity}Repository) BulkUpdate(ctx context.Context, ids []uuid.UUID, fields map[string]interface{}) error {
    if len(ids) == 0 {
        return nil
    }
    
    db := r.dbCtx.GetDB()
    fields["updated_at"] = time.Now()
    
    result := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id IN ? AND deleted_at IS NULL", ids).Updates(fields)
    
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to bulk update {entity}s", result.Error)
    }
    
    return nil
}
```

##### Delete Operations
```go
// infras/repository/gorm-pgsql/delete_{entity}.go
package repository

import (
    "context"
    "time"
    "{project_name}/ikv/modules/{module_name}/model"
    "github.com/google/uuid"
)

// Delete xóa vĩnh viễn entity
func (r *{Entity}Repository) Delete(ctx context.Context, id uuid.UUID) error {
    db := r.dbCtx.GetDB()
    
    result := db.WithContext(ctx).Where("id = ?", id).Delete(&{module_name}model.{Entity}{})
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to delete {entity}", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return {module_name}model.Err{Entity}NotFound
    }
    
    return nil
}

// SoftDelete xóa mềm entity (đánh dấu deleted_at)
func (r *{Entity}Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
    db := r.dbCtx.GetDB()
    
    now := time.Now()
    result := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id = ? AND deleted_at IS NULL", id).
        Updates(map[string]interface{}{
            "deleted_at": &now,
            "status":     {module_name}model.StatusDeleted,
            "updated_at": now,
        })
    
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to soft delete {entity}", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return {module_name}model.Err{Entity}NotFound
    }
    
    return nil
}

// Restore khôi phục entity đã bị soft delete
func (r *{Entity}Repository) Restore(ctx context.Context, id uuid.UUID) error {
    db := r.dbCtx.GetDB()
    
    result := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id = ? AND deleted_at IS NOT NULL", id).
        Updates(map[string]interface{}{
            "deleted_at": nil,
            "status":     {module_name}model.StatusActive,
            "updated_at": time.Now(),
        })
    
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to restore {entity}", result.Error)
    }
    
    if result.RowsAffected == 0 {
        return {module_name}model.Err{Entity}NotFound
    }
    
    return nil
}

// BulkDelete xóa nhiều entities cùng lúc
func (r *{Entity}Repository) BulkDelete(ctx context.Context, ids []uuid.UUID) error {
    if len(ids) == 0 {
        return nil
    }
    
    db := r.dbCtx.GetDB()
    
    result := db.WithContext(ctx).Where("id IN ?", ids).Delete(&{module_name}model.{Entity}{})
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to bulk delete {entity}s", result.Error)
    }
    
    return nil
}

// BulkSoftDelete xóa mềm nhiều entities cùng lúc
func (r *{Entity}Repository) BulkSoftDelete(ctx context.Context, ids []uuid.UUID) error {
    if len(ids) == 0 {
        return nil
    }
    
    db := r.dbCtx.GetDB()
    now := time.Now()
    
    result := db.WithContext(ctx).Model(&{module_name}model.{Entity}{}).
        Where("id IN ? AND deleted_at IS NULL", ids).
        Updates(map[string]interface{}{
            "deleted_at": &now,
            "status":     {module_name}model.StatusDeleted,
            "updated_at": now,
        })
    
    if result.Error != nil {
        return {module_name}model.NewDatabaseError("Failed to bulk soft delete {entity}s", result.Error)
    }
    
    return nil
}
```

### 2. CQRS Pattern Implementation

#### Command Pattern
```go
// service/commands/create_{entity}.go
package service

import (
    "context"
    "time"
    "github.com/google/uuid"
    "{project_name}/ikv/modules/{module_name}/model"
)

// Create Command
type Create{Entity}Command struct {
    Name        string    `json:"name" validate:"required,min=2,max=100"`
    Email       string    `json:"email" validate:"required,email"`
    Description string    `json:"description"`
    // Các fields khác
}

// Create Command Handler Interface
type ICreate{Entity}Handler interface {
    Execute(ctx context.Context, cmd *Create{Entity}Command) (*uuid.UUID, error)
}

// Create Command Handler Implementation
type Create{Entity}Handler struct {
    repository {module_name}model.ICreate{Entity}Repository
}

func NewCreate{Entity}Handler(repo {module_name}model.ICreate{Entity}Repository) ICreate{Entity}Handler {
    return &Create{Entity}Handler{repository: repo}
}

func (h *Create{Entity}Handler) Execute(ctx context.Context, cmd *Create{Entity}Command) (*uuid.UUID, error) {
    // Business logic validation
    if err := h.validateCreateCommand(cmd); err != nil {
        return nil, err
    }
    
    // Tạo entity từ command
    entity := &{module_name}model.{Entity}{
        ID:          uuid.New(),
        Name:        cmd.Name,
        Email:       cmd.Email,
        Description: cmd.Description,
        Status:      {module_name}model.StatusActive,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Business rules
    if err := h.applyBusinessRules(entity); err != nil {
        return nil, err
    }
    
    // Persist
    if err := h.repository.Insert(ctx, entity); err != nil {
        return nil, err
    }
    
    return &entity.ID, nil
}

func (h *Create{Entity}Handler) validateCreateCommand(cmd *Create{Entity}Command) error {
    // Custom validation logic
    return nil
}

func (h *Create{Entity}Handler) applyBusinessRules(entity *{module_name}model.{Entity}) error {
    // Business rules implementation
    return nil
}
```

#### Query Pattern
```go
// service/queries/get_{entity}_details.go
package service

import (
    "context"
    "github.com/google/uuid"
    "{project_name}/ikv/modules/{module_name}/model"
)

// Get Details Query
type Get{Entity}DetailsQuery struct {
    ID uuid.UUID `json:"id" validate:"required"`
}

// Query Handler Interface
type IGet{Entity}DetailsHandler interface {
    Execute(ctx context.Context, query *Get{Entity}DetailsQuery) (*{module_name}model.{Entity}, error)
}

// Query Handler Implementation
type Get{Entity}DetailsHandler struct {
    repository {module_name}model.IRead{Entity}Repository
}

func NewGet{Entity}DetailsHandler(repo {module_name}model.IRead{Entity}Repository) IGet{Entity}DetailsHandler {
    return &Get{Entity}DetailsHandler{repository: repo}
}

func (h *Get{Entity}DetailsHandler) Execute(ctx context.Context, query *Get{Entity}DetailsQuery) (*{module_name}model.{Entity}, error) {
    // Query validation
    if query.ID == uuid.Nil {
        return nil, {module_name}model.ErrInvalidID
    }
    
    // Fetch data
    entity, err := h.repository.GetByID(ctx, query.ID)
    if err != nil {
        return nil, err
    }
    
    // Business logic for read operations
    if err := h.applyReadRules(entity); err != nil {
        return nil, err
    }
    
    return entity, nil
}

func (h *Get{Entity}DetailsHandler) applyReadRules(entity *{module_name}model.{Entity}) error {
    // Read-specific business rules
    return nil
}
```

## CRUD Implementation chuẩn

### 1. HTTP Controllers

#### Base Controller Structure
```go
// infras/controller/http-gin/base_{entity}_controller.go
package httpgin

import (
    "{project_name}/ikv/modules/{module_name}/service"
)

type {Entity}HTTPController struct {
    // Command handlers
    createHandler service.ICreate{Entity}Handler
    updateHandler service.IUpdate{Entity}Handler
    deleteHandler service.IDelete{Entity}Handler
    
    // Query handlers
    getDetailsHandler service.IGet{Entity}DetailsHandler
    listHandler      service.IList{Entity}Handler
}

func New{Entity}HTTPController(
    createHandler service.ICreate{Entity}Handler,
    updateHandler service.IUpdate{Entity}Handler,
    deleteHandler service.IDelete{Entity}Handler,
    getDetailsHandler service.IGet{Entity}DetailsHandler,
    listHandler service.IList{Entity}Handler,
) *{Entity}HTTPController {
    return &{Entity}HTTPController{
        createHandler:     createHandler,
        updateHandler:     updateHandler,
        deleteHandler:     deleteHandler,
        getDetailsHandler: getDetailsHandler,
        listHandler:       listHandler,
    }
}
```

#### Create API
```go
// infras/controller/http-gin/create_{entity}_api.go
package httpgin

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/shared/datatype"
)

// ActionCreate xử lý tạo {entity} mới - POST /
func (c *{Entity}HTTPController) ActionCreate(ctx *gin.Context) {
    var cmd service.Create{Entity}Command
    
    // Bind JSON request
    if err := ctx.ShouldBindJSON(&cmd); err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "Dữ liệu đầu vào không hợp lệ",
            err,
        ))
        return
    }
    
    // Validate request với custom validation
    if err := c.validateCreateRequest(&cmd); err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeValidation,
            "Validation failed",
            err,
        ))
        return
    }
    
    // Execute command
    id, err := c.createHandler.Execute(ctx.Request.Context(), &cmd)
    if err != nil {
        // Handle specific errors
        switch err.(type) {
        case *{module_name}model.DomainError:
            ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
                datatype.ErrCodeBusiness,
                err.Error(),
                err,
            ))
        default:
            ctx.JSON(http.StatusInternalServerError, datatype.NewAppError(
                datatype.ErrCodeInternal,
                "Không thể tạo {entity}",
                err,
            ))
        }
        return
    }
    
    // Success response
    ctx.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data": gin.H{
            "id": id,
        },
        "message": "{Entity} đã được tạo thành công",
    })
}

// validateCreateRequest kiểm tra custom validation cho create request
func (c *{Entity}HTTPController) validateCreateRequest(cmd *service.Create{Entity}Command) error {
    // Custom business validation logic
    // Ví dụ: kiểm tra email format, name length, etc.
    return nil
}
```

#### Get Detail API
```go
// infras/controller/http-gin/get_{entity}_detail_api.go
package httpgin

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/shared/datatype"
)

// ActionGetDetail lấy chi tiết {entity} - GET /:id
func (c *{Entity}HTTPController) ActionGetDetail(ctx *gin.Context) {
    // Parse và validate ID
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "ID không hợp lệ",
            err,
        ))
        return
    }
    
    // Create query
    query := &service.Get{Entity}DetailsQuery{ID: id}
    
    // Execute query
    entity, err := c.getDetailsHandler.Execute(ctx.Request.Context(), query)
    if err != nil {
        // Handle specific errors
        switch err {
        case {module_name}model.Err{Entity}NotFound:
            ctx.JSON(http.StatusNotFound, datatype.NewAppError(
                datatype.ErrCodeNotFound,
                "Không tìm thấy {entity}",
                err,
            ))
        default:
            ctx.JSON(http.StatusInternalServerError, datatype.NewAppError(
                datatype.ErrCodeInternal,
                "Lỗi khi lấy thông tin {entity}",
                err,
            ))
        }
        return
    }
    
    // Convert to response DTO
    response := entity.ToResponse()
    
    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    response,
    })
}
```

#### List API
```go
// infras/controller/http-gin/list_{entity}_api.go
package httpgin

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    "{project_name}/ikv/modules/{module_name}/model"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/shared/datatype"
)

// ActionGetList lấy danh sách {entity} - GET /
func (c *{Entity}HTTPController) ActionGetList(ctx *gin.Context) {
    // Parse query parameters
    filter, err := c.parseListQueryParams(ctx)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "Tham số truy vấn không hợp lệ",
            err,
        ))
        return
    }
    
    // Create query
    query := &service.List{Entity}Query{
        Filter: filter,
    }
    
    // Execute query
    entities, total, err := c.listHandler.Execute(ctx.Request.Context(), query)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, datatype.NewAppError(
            datatype.ErrCodeInternal,
            "Lỗi khi lấy danh sách {entity}",
            err,
        ))
        return
    }
    
    // Convert to response
    response := {module_name}model.ToListResponse(entities, total, filter.Page, filter.PerPage)
    
    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    response,
    })
}

// parseListQueryParams parse các query parameters cho list API
func (c *{Entity}HTTPController) parseListQueryParams(ctx *gin.Context) (*{module_name}model.ListFilter, error) {
    // Default values
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
    
    // Validate pagination
    if page < 1 {
        page = 1
    }
    if perPage < 1 || perPage > 100 {
        perPage = 10
    }
    
    // Parse dates
    var createdFrom, createdTo time.Time
    if dateStr := ctx.Query("created_from"); dateStr != "" {
        if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
            createdFrom = parsed
        }
    }
    if dateStr := ctx.Query("created_to"); dateStr != "" {
        if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
            createdTo = parsed
        }
    }
    
    return &{module_name}model.ListFilter{
        Page:        page,
        PerPage:     perPage,
        Status:      ctx.Query("status"),
        Search:      ctx.Query("search"),
        SortBy:      ctx.DefaultQuery("sort_by", "created_at"),
        SortOrder:   ctx.DefaultQuery("sort_order", "DESC"),
        CreatedFrom: createdFrom,
        CreatedTo:   createdTo,
    }, nil
}
```

#### Update API
```go
// infras/controller/http-gin/update_{entity}_api.go
package httpgin

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/shared/datatype"
)

// ActionUpdate cập nhật {entity} - PUT /:id
func (c *{Entity}HTTPController) ActionUpdate(ctx *gin.Context) {
    // Parse và validate ID
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "ID không hợp lệ",
            err,
        ))
        return
    }
    
    // Bind JSON request
    var cmd service.Update{Entity}Command
    if err := ctx.ShouldBindJSON(&cmd); err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "Dữ liệu đầu vào không hợp lệ",
            err,
        ))
        return
    }
    
    // Set ID from URL param
    cmd.ID = id
    
    // Validate request
    if err := c.validateUpdateRequest(&cmd); err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeValidation,
            "Validation failed",
            err,
        ))
        return
    }
    
    // Execute command
    err = c.updateHandler.Execute(ctx.Request.Context(), &cmd)
    if err != nil {
        // Handle specific errors
        switch err {
        case {module_name}model.Err{Entity}NotFound:
            ctx.JSON(http.StatusNotFound, datatype.NewAppError(
                datatype.ErrCodeNotFound,
                "Không tìm thấy {entity} để cập nhật",
                err,
            ))
        default:
            switch err.(type) {
            case *{module_name}model.DomainError:
                ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
                    datatype.ErrCodeBusiness,
                    err.Error(),
                    err,
                ))
            default:
                ctx.JSON(http.StatusInternalServerError, datatype.NewAppError(
                    datatype.ErrCodeInternal,
                    "Không thể cập nhật {entity}",
                    err,
                ))
            }
        }
        return
    }
    
    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "{Entity} đã được cập nhật thành công",
    })
}

// validateUpdateRequest kiểm tra custom validation cho update request
func (c *{Entity}HTTPController) validateUpdateRequest(cmd *service.Update{Entity}Command) error {
    // Custom business validation logic
    return nil
}
```

#### Delete API
```go
// infras/controller/http-gin/delete_{entity}_api.go
package httpgin

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "{project_name}/ikv/modules/{module_name}/service"
    "{project_name}/ikv/shared/datatype"
)

// ActionDelete xóa {entity} - DELETE /:id
func (c *{Entity}HTTPController) ActionDelete(ctx *gin.Context) {
    // Parse và validate ID
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
            datatype.ErrCodeInvalidInput,
            "ID không hợp lệ",
            err,
        ))
        return
    }
    
    // Get delete type from query param (soft/hard)
    deleteType := ctx.DefaultQuery("type", "soft")
    
    // Create command based on delete type
    var command interface{}
    if deleteType == "hard" {
        command = &service.HardDelete{Entity}Command{ID: id}
    } else {
        command = &service.SoftDelete{Entity}Command{ID: id}
    }
    
    // Execute appropriate command
    err = c.executeDeleteCommand(ctx, command)
    if err != nil {
        // Handle specific errors
        switch err {
        case {module_name}model.Err{Entity}NotFound:
            ctx.JSON(http.StatusNotFound, datatype.NewAppError(
                datatype.ErrCodeNotFound,
                "Không tìm thấy {entity} để xóa",
                err,
            ))
        default:
            switch err.(type) {
            case *{module_name}model.DomainError:
                ctx.JSON(http.StatusBadRequest, datatype.NewAppError(
                    datatype.ErrCodeBusiness,
                    err.Error(),
                    err,
                ))
            default:
                ctx.JSON(http.StatusInternalServerError, datatype.NewAppError(
                    datatype.ErrCodeInternal,
                    "Không thể xóa {entity}",
                    err,
                ))
            }
        }
        return
    }
    
    // Success response
    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "{Entity} đã được xóa thành công",
    })
}

// executeDeleteCommand thực thi delete command tương ứng
func (c *{Entity}HTTPController) executeDeleteCommand(ctx *gin.Context, command interface{}) error {
    switch cmd := command.(type) {
    case *service.SoftDelete{Entity}Command:
        return c.deleteHandler.ExecuteSoft(ctx.Request.Context(), cmd)
    case *service.HardDelete{Entity}Command:
        return c.deleteHandler.ExecuteHard(ctx.Request.Context(), cmd)
    default:
        return {module_name}model.NewValidationError("delete_type", "Invalid delete command type")
    }
}
```

### 2. URL Routing

```go
// urls/v1/{entity}_url.go
package v1

import (
    "net/http"
    httpgin "fat2fast/ikv/modules/{module_name}/infras/controller/http-gin"
    "github.com/gin-gonic/gin"
)

func GetRoutes(controller *httpgin.{Entity}HTTPController) []gin.RouteInfo {
    return []gin.RouteInfo{
        {
            Method:      http.MethodGet,
            Path:        "",
            HandlerFunc: controller.ActionGetList,
        },
        {
            Method:      http.MethodGet,
            Path:        "/:id",
            HandlerFunc: controller.ActionGetDetail,
        },
        {
            Method:      http.MethodPost,
            Path:        "",
            HandlerFunc: controller.ActionCreate,
        },
        {
            Method:      http.MethodPut,
            Path:        "/:id",
            HandlerFunc: controller.ActionUpdate,
        },
        {
            Method:      http.MethodDelete,
            Path:        "/:id",
            HandlerFunc: controller.ActionDelete,
        },
    }
}
```

## Model và DTOs

### 1. Domain Entity

```go
// model/{entity}.go
package {module_name}model

import (
    "time"
    "github.com/google/uuid"
)

type {Entity} struct {
    ID          uuid.UUID `json:"id" gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
    Name        string    `json:"name" gorm:"column:name;not null"`
    Email       string    `json:"email" gorm:"column:email;uniqueIndex"`
    Description string    `json:"description" gorm:"column:description"`
    Status      string    `json:"status" gorm:"column:status;default:'active'"`
    CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// Table name
func ({Entity}) TableName() string {
    return "{entity_table_name}"
}

// Status constants
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    StatusPending  = "pending"
    StatusDeleted  = "deleted"
)

// Business methods
func (e *{Entity}) IsActive() bool {
    return e.Status == StatusActive
}

func (e *{Entity}) Activate() {
    e.Status = StatusActive
    e.UpdatedAt = time.Now()
}

func (e *{Entity}) Deactivate() {
    e.Status = StatusInactive
    e.UpdatedAt = time.Now()
}
```

### 2. DTOs

```go
// model/dtos.go
package {module_name}model

import (
    "time"
    "github.com/google/uuid"
)

// Request DTOs
type Create{Entity}Request struct {
    Name        string `json:"name" binding:"required,min=2,max=100"`
    Email       string `json:"email" binding:"required,email"`
    Description string `json:"description" binding:"max=500"`
}

type Update{Entity}Request struct {
    Name        string `json:"name" binding:"omitempty,min=2,max=100"`
    Email       string `json:"email" binding:"omitempty,email"`
    Description string `json:"description" binding:"omitempty,max=500"`
    Status      string `json:"status" binding:"omitempty,oneof=active inactive pending"`
}

// Response DTOs
type {Entity}Response struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type {Entity}ListResponse struct {
    Items      []*{Entity}Response `json:"items"`
    TotalCount int64               `json:"total_count"`
    Page       int                 `json:"page"`
    PerPage    int                 `json:"per_page"`
    TotalPages int                 `json:"total_pages"`
}

// Filter DTOs
type ListFilter struct {
    Page    int    `json:"page" binding:"omitempty,min=1"`
    PerPage int    `json:"per_page" binding:"omitempty,min=1,max=100"`
    Status  string `json:"status" binding:"omitempty,oneof=active inactive pending deleted"`
    Search  string `json:"search" binding:"omitempty,max=100"`
}

// Conversion methods
func (e *{Entity}) ToResponse() *{Entity}Response {
    return &{Entity}Response{
        ID:          e.ID,
        Name:        e.Name,
        Email:       e.Email,
        Description: e.Description,
        Status:      e.Status,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,
    }
}

func ToListResponse(entities []*{Entity}, total int64, page, perPage int) *{Entity}ListResponse {
    items := make([]*{Entity}Response, len(entities))
    for i, entity := range entities {
        items[i] = entity.ToResponse()
    }
    
    totalPages := int((total + int64(perPage) - 1) / int64(perPage))
    
    return &{Entity}ListResponse{
        Items:      items,
        TotalCount: total,
        Page:       page,
        PerPage:    perPage,
        TotalPages: totalPages,
    }
}
```

## Migration và Database

### 1. Migration Files

```sql
-- migrations/001_create_{entity}_table.up.sql
CREATE SCHEMA IF NOT EXISTS {module_name}_schema;

CREATE TABLE IF NOT EXISTS {module_name}_schema.{entity_table_name} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'pending', 'deleted')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_{entity_table_name}_status ON {module_name}_schema.{entity_table_name}(status);
CREATE INDEX idx_{entity_table_name}_created_at ON {module_name}_schema.{entity_table_name}(created_at);
CREATE INDEX idx_{entity_table_name}_deleted_at ON {module_name}_schema.{entity_table_name}(deleted_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_{entity_table_name}_updated_at
    BEFORE UPDATE ON {module_name}_schema.{entity_table_name}
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- migrations/001_create_{entity}_table.down.sql
DROP TRIGGER IF EXISTS update_{entity_table_name}_updated_at ON {module_name}_schema.{entity_table_name};
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS {module_name}_schema.{entity_table_name};
DROP SCHEMA IF EXISTS {module_name}_schema CASCADE;
```

## Error Handling

### 1. Custom Errors

```go
// model/error.go
package {module_name}model

import (
    "errors"
    "fmt"
)

// Domain errors
var (
    ErrInvalidID       = errors.New("invalid ID")
    Err{Entity}NotFound = errors.New("{entity} not found")
    ErrDuplicate{Entity} = errors.New("{entity} already exists")
    ErrInvalidStatus   = errors.New("invalid status")
)

// Error types
type DomainError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e DomainError) Error() string {
    if e.Details != "" {
        return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Factory functions
func NewValidationError(field, message string) *DomainError {
    return &DomainError{
        Code:    "VALIDATION_ERROR",
        Message: fmt.Sprintf("Validation failed for field '%s'", field),
        Details: message,
    }
}

func NewNotFoundError(entityType, id string) *DomainError {
    return &DomainError{
        Code:    "NOT_FOUND",
        Message: fmt.Sprintf("%s with ID '%s' not found", entityType, id),
    }
}

func NewConflictError(message string) *DomainError {
    return &DomainError{
        Code:    "CONFLICT",
        Message: message,
    }
}
```

## Testing

### 1. Unit Test Example

```go
// service/commands/create_{entity}_test.go
package service_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "{project_name}/ikv/modules/{module_name}/model"
    "{project_name}/ikv/modules/{module_name}/service"
)

// Mock repository
type Mock{Entity}Repository struct {
    mock.Mock
}

func (m *Mock{Entity}Repository) Insert(ctx context.Context, entity *{module_name}model.{Entity}) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

func TestCreate{Entity}Handler_Execute(t *testing.T) {
    // Arrange
    mockRepo := new(Mock{Entity}Repository)
    handler := service.NewCreate{Entity}Handler(mockRepo)
    
    cmd := &service.Create{Entity}Command{
        Name:        "Test {Entity}",
        Email:       "test@example.com",
        Description: "Test description",
    }
    
    mockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*{module_name}model.{Entity}")).
        Return(nil)
    
    // Act
    id, err := handler.Execute(context.Background(), cmd)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, id)
    assert.NotEqual(t, uuid.Nil, *id)
    mockRepo.AssertExpectations(t)
}
```

## Best Practices

### 1. Naming Conventions
- **Modules**: lowercase, singular (`user`, `product`, `order`)
- **Entities**: PascalCase, singular (`User`, `Product`, `Order`)
- **Interfaces**: Prefix với `I` (`IUserRepository`, `ICreateHandler`)
- **Files**: snake_case (`create_user.go`, `user_repository.go`)

### 2. Dependency Direction
```
Infrastructure → Application → Domain
     ↑              ↑
Controller    ←  Service  ←  Entity
Repository
```

### 3. Error Handling Flow
- **Domain errors**: Business logic violations
- **Application errors**: Use case failures
- **Infrastructure errors**: Database, network issues
- **Presentation errors**: HTTP status codes

### 4. Validation Layers
- **DTO validation**: Struct tags (`binding:"required"`)
- **Business validation**: Command/Query handlers
- **Domain validation**: Entity methods

### 5. Transaction Management
```go
func (h *Create{Entity}Handler) Execute(ctx context.Context, cmd *Create{Entity}Command) (*uuid.UUID, error) {
    // Start transaction if needed
    tx := h.dbCtx.BeginTx()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    // Business logic
    entity := h.buildEntity(cmd)
    
    if err := h.repository.Insert(ctx, entity); err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    return &entity.ID, nil
}
```

## Checklist cho Module mới

### Cấu trúc Files
- [ ] `module.go` - Entry point với dependency injection
- [ ] `config.yaml` - Configuration với environment variables
- [ ] `README.md` - Documentation chi tiết
- [ ] `model/` - Domain entities, DTOs, errors
- [ ] `service/` - CQRS handlers (commands/, queries/)
- [ ] `infras/controller/` - HTTP controllers
- [ ] `infras/repository/` - Database implementations
- [ ] `urls/v1/` - API routing
- [ ] `migrations/` - Database schema

### Implementation
- [ ] Repository interfaces và implementations
- [ ] Command handlers cho write operations
- [ ] Query handlers cho read operations
- [ ] HTTP controllers với proper error handling
- [ ] URL routing với RESTful endpoints
- [ ] Database migrations (up/down)
- [ ] Custom errors và validation

### Testing
- [ ] Unit tests cho services
- [ ] Integration tests cho repositories
- [ ] API tests cho controllers
- [ ] Mock implementations

### Documentation
- [ ] README.md với usage examples
- [ ] API documentation
- [ ] Database schema documentation
- [ ] Deployment instructions

Quy chuẩn này đảm bảo tính nhất quán, khả năng bảo trì và mở rộng cho tất cả modules trong hệ thống IKV.