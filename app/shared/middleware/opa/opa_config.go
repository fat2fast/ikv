package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OPAConfig chứa cấu hình cho OPA middleware
type OPAConfig struct {
	// URL của OPA server
	// Ví dụ: "http://localhost:8181"
	URL string

	// Query là câu query để đánh giá policy
	// Ví dụ: "data.policy.allow"
	Query string

	// InputCreationMethod là function để tạo input data từ gin.Context
	// Input này sẽ được gửi đến OPA để đánh giá
	InputCreationMethod func(c *gin.Context) (map[string]interface{}, error)

	// ExceptedResult là kết quả mong đợi từ OPA query
	// Nếu kết quả khác ExceptedResult thì request sẽ bị từ chối
	ExceptedResult interface{}

	// DeniedStatusCode là HTTP status code trả về khi request bị từ chối
	// Mặc định là 403 (Forbidden)
	DeniedStatusCode int

	// DeniedMessage là message trả về khi request bị từ chối
	DeniedMessage string

	// Context cho các request đến OPA server
	Context context.Context

	// HTTPClient để gọi remote OPA server
	HTTPClient *http.Client
}

// SetDefaults thiết lập các giá trị mặc định cho config
func (c *OPAConfig) SetDefaults() {
	if c.DeniedStatusCode == 0 {
		c.DeniedStatusCode = http.StatusForbidden
	}
	if c.DeniedMessage == "" {
		c.DeniedMessage = "Access denied by policy"
	}
	if c.Context == nil {
		c.Context = context.Background()
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	if c.ExceptedResult == nil {
		c.ExceptedResult = true
	}
}
