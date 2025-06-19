package shared

import (
	"net/http"

	bookv1 "fat2fast/ikv/modules/book/urls/v1"

	"github.com/gin-gonic/gin"
)

// RouteRegistry quản lý tất cả modules
type RouteRegistry struct {
	modules map[string]func() []gin.RouteInfo
}

// NewRouteRegistry tạo instance mới
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		modules: make(map[string]func() []gin.RouteInfo),
	}
}

// Register đăng ký module với prefix
func (rr *RouteRegistry) Register(prefix string, getRoutes func() []gin.RouteInfo) {
	rr.modules[prefix] = getRoutes
}

// GetAllRoutes lấy tất cả routes với prefix
func (rr *RouteRegistry) GetAllRoutes() []gin.RouteInfo {
	var allRoutes []gin.RouteInfo

	for prefix, getRoutes := range rr.modules {
		routes := getRoutes()

		// Thêm prefix vào path của mỗi route
		for _, route := range routes {
			route.Path = "/" + prefix + route.Path
			allRoutes = append(allRoutes, route)
		}
	}

	return allRoutes
}

// ActionPing handler mặc định cho ping
func ActionPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"status":  "ok",
	})
}

// GetUrl trả về tất cả routes từ scan shared/urls directory
func GetUrl() []gin.RouteInfo {
	// Tạo registry
	registry := NewRouteRegistry()
	groupV1 := "v1"
	// registry.Register(groupV1+"/user", userv1.GetRoutes)
	registry.Register(groupV1+"/book", bookv1.GetRoutes)

	// Lấy tất cả routes
	allRoutes := registry.GetAllRoutes()

	// Thêm các route cơ bản
	baseRoutes := []gin.RouteInfo{
		{
			Method:      http.MethodGet,
			Path:        "/health",
			HandlerFunc: ActionPing,
		},
	}

	// Kết hợp tất cả routes
	return append(baseRoutes, allRoutes...)
}
