package v1

import (
	"net/http"

	userhttpgin "fat2fast/ikv/modules/user/infras/controller/http-gin"

	"github.com/gin-gonic/gin"
)

// GetUserRoutes trả về danh sách routes cho user module v1
func GetRoutes(controller *userhttpgin.UserHTTPController) []gin.RouteInfo {

	return []gin.RouteInfo{
		{
			Method:      http.MethodPost,
			Path:        "/authenticate",
			HandlerFunc: controller.ActionAuthenticate,
		},
		{
			Method:      http.MethodPost,
			Path:        "/register",
			HandlerFunc: controller.ActionRegister,
		},
		{
			Method:      http.MethodGet,
			Path:        "/profile/:id",
			HandlerFunc: controller.ActionGetProfile,
		},
		{
			Method:      http.MethodPut,
			Path:        "/profile/:id",
			HandlerFunc: controller.ActionUpdateProfile,
		},
	}
}
