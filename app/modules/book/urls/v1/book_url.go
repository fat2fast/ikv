package v1

import (
	"net/http"

	bookhttpgin "fat2fast/ikv/modules/book/infras/controller/http-gin"

	"github.com/gin-gonic/gin"
)

// GetRoutes trả về danh sách routes cho book module v1
func GetRoutes(controller *bookhttpgin.BookHTTPController) []gin.RouteInfo {
	return []gin.RouteInfo{
		// GET / - Lấy danh sách books
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: controller.ActionListBooks,
		},
		// GET /:id - Lấy chi tiết book
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: controller.ActionGetBookDetail,
		},
		// POST / - Tạo book mới
		{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: controller.ActionCreateBook,
		},
		// PUT /:id - Cập nhật book
		{
			Method:      http.MethodPut,
			Path:        "/:id",
			HandlerFunc: controller.ActionUpdateBook,
		},
		// DELETE /:id - Xóa book
		{
			Method:      http.MethodDelete,
			Path:        "/:id",
			HandlerFunc: controller.ActionDeleteBook,
		},
	}
}
