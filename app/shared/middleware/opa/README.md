# OPA Gin Middleware (Remote Mode)

Open Policy Agent (OPA) middleware cho Gin framework, hỗ trợ remote OPA server.

## Tính năng

- ✅ Hỗ trợ remote OPA server
- ✅ Cấu hình linh hoạt cho input creation
- ✅ Error handling tích hợp với hệ thống datatype của project
- ✅ Examples và templates sẵn có
- ✅ Tích hợp JWT claims và user context
- ✅ Performance tối ưu với HTTP client pooling

## Cài đặt

Dependency đã được thêm vào project:
```bash
go get github.com/open-policy-agent/opa@latest
```

## Cấu hình

### OPAConfig struct

```go
type OPAConfig struct {
    // URL: Địa chỉ OPA server  
    URL string
    
    // Query: Câu query để đánh giá (ví dụ: "data.policy.allow")
    Query string
    
    // InputCreationMethod: Function tạo input từ gin.Context
    InputCreationMethod func(c *gin.Context) (map[string]interface{}, error)
    
    // ExceptedResult: Kết quả mong đợi (mặc định: true)
    ExceptedResult interface{}
    
    // DeniedStatusCode: HTTP status khi bị từ chối (mặc định: 403)
    DeniedStatusCode int
    
    // DeniedMessage: Message khi bị từ chối
    DeniedMessage string
    
    // Context: Context cho OPA operations
    Context context.Context
    
    // HTTPClient: HTTP client cho remote calls
    HTTPClient *http.Client
}
```

## Cách sử dụng

### 1. Basic Setup

```go
import "fat2fast/ikv/shared/middleware"

func setupOPAMiddleware() gin.HandlerFunc {
    opaConfig := &middleware.OPAConfig{
        URL: "http://localhost:8181",
        Query: "data.policy.allow",
        InputCreationMethod: middleware.BasicInputCreator,
        ExceptedResult: true,
        DeniedMessage: "Access denied by policy",
    }
    
    opaMiddleware, err := middleware.NewOPAMiddleware(opaConfig)
    if err != nil {
        panic(err)
    }
    
    return opaMiddleware.Use()
}
```

### 2. With User Context

```go
func setupUserOPA() gin.HandlerFunc {
    opaConfig := &middleware.OPAConfig{
        URL: "http://localhost:8181",
        Query: "data.authz.allow",
        InputCreationMethod: middleware.UserInputCreator,
        ExceptedResult: true,
        DeniedMessage: "Access denied by authorization service",
    }
    
    opaMiddleware, err := middleware.NewOPAMiddleware(opaConfig)
    if err != nil {
        panic(err)
    }
    
    return opaMiddleware.Use()
}
```

### 3. Apply vào Routes

```go
func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")
    
    // Apply OPA middleware cho tất cả routes
    api.Use(setupUserOPA())
    
    api.GET("/users", listUsers)
    api.POST("/users", createUser)
    api.GET("/users/:id", getUser)
}
```

## Input Creation Methods có sẵn

### BasicInputCreator
Tạo input cơ bản với method và path:
```go
middleware.BasicInputCreator(c)
// Output: {"method": "GET", "path": "/api/v1/users"}
```

### UserInputCreator  
Tạo input với thông tin user từ JWT:
```go
middleware.UserInputCreator(c)
// Output: {"method": "GET", "path": "/api/v1/users", "user": {...}}
```

### DetailedInputCreator
Tạo input chi tiết với headers, IP, query params:
```go
middleware.DetailedInputCreator(c)
// Output: {"method": "GET", "path": "/api/v1/users", "user": {...}, "headers": {...}, "client_ip": "127.0.0.1"}
```

## Custom Input Creator

Tạo custom input creator:

```go
func MyCustomInputCreator(c *gin.Context) (map[string]interface{}, error) {
    input := map[string]interface{}{
        "method": c.Request.Method,
        "path":   c.Request.URL.Path,
    }
    
    // Thêm custom fields
    if orgID := c.GetHeader("X-Org-ID"); orgID != "" {
        input["organization_id"] = orgID
    }
    
    if user, exists := c.Get("user"); exists {
        input["user"] = user
    }
    
    // Thêm resource ID từ URL params
    if resourceID := c.Param("id"); resourceID != "" {
        input["resource_id"] = resourceID
    }
    
    return input, nil
}
```

## Error Handling

Middleware tự động xử lý lỗi và trả về response theo format của project:

- **500 Internal Server Error**: Khi có lỗi trong quá trình đánh giá policy
- **403 Forbidden**: Khi policy từ chối request

Response format:
```json
{
    "code": 403,
    "status": "Forbidden", 
    "message": "The requested action was forbidden",
    "reason": "Access denied by policy",
    "debug": "Policy result: false"
}
```

## Setup OPA Server

### 1. Tạo Policy Files

```bash
# policy.rego - Basic authorization
cat > policy.rego << EOF
package policy

default allow = false

# Allow GET requests for all users
allow {
    input.method = "GET"
}

# Allow all methods for admin
allow {
    input.user.role = "admin"
}
EOF
```

```bash
# authz.rego - Role-based authorization
cat > authz.rego << EOF
package authz

default allow = false

# Admin có thể làm tất cả
allow {
    input.user.role = "admin"
}

# User chỉ có thể đọc dữ liệu của mình
allow {
    input.user.role = "user"
    input.method = "GET"
    startswith(input.path, sprintf("/api/v1/users/%s", [input.user.id]))
}

# User có thể update profile của mình
allow {
    input.user.role = "user" 
    input.method = "PUT"
    input.path = sprintf("/api/v1/users/%s/profile", [input.user.id])
}
EOF
```

### 2. Chạy OPA Server

```bash
# Chạy OPA server với policies
opa run --server --addr localhost:8181 policy.rego authz.rego

# Hoặc từ thư mục chứa policies
opa run --server --addr localhost:8181 *.rego
```

### 3. Test Policy

```bash
# Test với curl
curl -X POST http://localhost:8181/v1/data/policy/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "method": "GET",
      "path": "/api/v1/users",
      "user": {"role": "user", "id": "123"}
    }
  }'
```

## Best Practices

1. **Policy Organization**: Tổ chức policies theo modules (authz, rbac, etc.)
2. **Input Size**: Giữ input size nhỏ để tránh overhead
3. **Caching**: Sử dụng HTTP client với connection pooling
4. **Monitoring**: Monitor OPA server health và response time
5. **Testing**: Test policies trước khi deploy

## OPA Server Production Setup

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  opa:
    image: openpolicyagent/opa:latest-envoy
    ports:
      - "8181:8181"
    volumes:
      - ./policies:/policies
    command:
      - "run"
      - "--server"
      - "--addr=0.0.0.0:8181"
      - "/policies"
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: opa-server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: opa-server
  template:
    metadata:
      labels:
        app: opa-server
    spec:
      containers:
      - name: opa
        image: openpolicyagent/opa:latest
        ports:
        - containerPort: 8181
        command:
        - "opa"
        - "run"
        - "--server"
        - "--addr=0.0.0.0:8181"
        - "/policies"
        volumeMounts:
        - name: policies
          mountPath: /policies
      volumes:
      - name: policies
        configMap:
          name: opa-policies
```

## Troubleshooting

### OPA server connection failed  
- Kiểm tra OPA server đang chạy: `curl http://localhost:8181/health`
- Verify URL và network connectivity
- Check firewall rules

### Policy evaluation failed
- Verify policy syntax với `opa fmt policies/`
- Test query với `opa eval -d policies/ "data.policy.allow"`
- Check logs của OPA server

### Input data không đúng
- Log input data để debug
- Use DetailedInputCreator để xem full input
- Validate input structure với policy expectations

### Performance issues
- Monitor OPA server metrics
- Optimize policy complexity
- Use appropriate HTTP client timeouts
- Consider OPA server clustering cho high availability 