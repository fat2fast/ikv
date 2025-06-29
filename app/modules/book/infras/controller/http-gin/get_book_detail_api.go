package bookhttpgin

import (
	"net/http"

	bookservice "fat2fast/ikv/modules/book/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActionGetBookDetail lấy chi tiết book - GET /:id
func (c *BookHTTPController) ActionGetBookDetail(ctx *gin.Context) {
	// Parse và validate ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug("Invalid book ID format"))
	}

	// Tạo query
	query := &bookservice.GetBookDetailQuery{ID: id}

	// Thực thi query
	response, err := c.getDetailQryHdl.Execute(ctx.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	// Trả về response
	ctx.JSON(http.StatusOK, datatype.ResponseSuccess(response))
}
