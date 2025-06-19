package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController handlers cho user routes
type BookController struct{}

// ActionGetUsers lấy danh sách users
func (uc *BookController) ActionGetList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users v1",
		"data":    []string{"user1", "user2"},
	})
}

// ActionGetUser lấy thông tin user theo ID
func (uc *BookController) ActionGetDetail(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Get user v1",
		"user_id": userID,
	})
}

// ActionCreateUser tạo user mới
func (uc *BookController) ActionCreate(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created v1",
	})
}

// UpdateUser cập nhật user
func (uc *BookController) ActionUpdate(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated v1",
		"user_id": userID,
	})
}

// DeleteUser xóa user
func (uc *BookController) ActionDelete(c *gin.Context) {
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted v1",
		"user_id": userID,
	})
}

// GetUserRoutes trả về danh sách routes cho user module v1
func GetRoutes() []gin.RouteInfo {
	controller := &BookController{}

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
