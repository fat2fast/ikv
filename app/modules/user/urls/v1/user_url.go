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
