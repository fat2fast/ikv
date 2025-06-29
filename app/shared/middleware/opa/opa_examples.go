package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Ví dụ các InputCreationMethod thông dụng

// BasicInputCreator tạo input cơ bản với method và path
func BasicInputCreator(c *gin.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}, nil
}

// UserInputCreator tạo input với thông tin user từ JWT claims
func UserInputCreator(c *gin.Context) (map[string]interface{}, error) {
	input := map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}

	// Lấy user từ JWT claims (giả sử đã được set bởi JWT middleware)
	if user, exists := c.Get("user"); exists {
		input["user"] = user
	}

	return input, nil
}

// DetailedInputCreator tạo input chi tiết với headers và query params
func DetailedInputCreator(c *gin.Context) (map[string]interface{}, error) {
	input := map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"query":  c.Request.URL.RawQuery,
	}

	// Thêm headers quan trọng
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			// Chỉ lấy các header quan trọng để tránh input quá lớn
			lowerKey := strings.ToLower(key)
			if lowerKey == "authorization" || lowerKey == "user-agent" || lowerKey == "x-real-ip" {
				headers[lowerKey] = values[0]
			}
		}
	}
	input["headers"] = headers

	// Lấy user từ context nếu có
	if user, exists := c.Get("user"); exists {
		input["user"] = user
	}

	// Thêm thông tin client IP
	input["client_ip"] = c.ClientIP()

	return input, nil
}

// Ví dụ cách sử dụng trong module

/*
// Trong file module.go hoặc router setup:

// 1. Sử dụng với remote OPA server - Basic
func setupBasicOPAMiddleware() gin.HandlerFunc {
	opaConfig := &middleware.OPAConfig{
		URL:                 "http://localhost:8181",
		Query:               "data.policy.allow",
		InputCreationMethod: middleware.BasicInputCreator,
		ExceptedResult:      true,
		DeniedMessage:       "Access denied by policy",
	}

	opaMiddleware, err := middleware.NewOPAMiddleware(opaConfig)
	if err != nil {
		panic(err)
	}

	return opaMiddleware.Use()
}

// 2. Sử dụng với user context
func setupUserOPAMiddleware() gin.HandlerFunc {
	opaConfig := &middleware.OPAConfig{
		URL:                 "http://localhost:8181",
		Query:               "data.authz.allow",
		InputCreationMethod: middleware.UserInputCreator,
		ExceptedResult:      true,
		DeniedMessage:       "Access denied by authorization policy",
	}

	opaMiddleware, err := middleware.NewOPAMiddleware(opaConfig)
	if err != nil {
		panic(err)
	}

	return opaMiddleware.Use()
}

// 3. Sử dụng với detailed input
func setupDetailedOPAMiddleware() gin.HandlerFunc {
	opaConfig := &middleware.OPAConfig{
		URL:                 "http://localhost:8181",
		Query:               "data.rbac.allow",
		InputCreationMethod: middleware.DetailedInputCreator,
		ExceptedResult:      true,
		DeniedMessage:       "Insufficient permissions",
	}

	opaMiddleware, err := middleware.NewOPAMiddleware(opaConfig)
	if err != nil {
		panic(err)
	}

	return opaMiddleware.Use()
}

// Cách apply middleware vào route group:
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Apply OPA middleware cho tất cả routes trong group
	api.Use(setupUserOPAMiddleware())

	api.GET("/users/:id", getUserHandler)
	api.PUT("/users/:id/profile", updateProfileHandler)
	api.GET("/admin/users", listUsersHandler)
}

// Custom Input Creator cho specific use case:
func CustomOrgInputCreator(c *gin.Context) (map[string]interface{}, error) {
	input := map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}

	// Thêm organization ID từ header
	if orgID := c.GetHeader("X-Org-ID"); orgID != "" {
		input["organization_id"] = orgID
	}

	// Thêm user context
	if user, exists := c.Get("user"); exists {
		input["user"] = user
	}

	// Thêm resource ID từ URL params
	if resourceID := c.Param("id"); resourceID != "" {
		input["resource_id"] = resourceID
	}

	return input, nil
}
*/
