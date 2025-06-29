package bookhttpgin

import (
	"net/http"

	bookmodel "fat2fast/ikv/modules/book/model"
	bookservice "fat2fast/ikv/modules/book/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActionUpdateBook cập nhật book - PUT /:id
func (c *BookHTTPController) ActionUpdateBook(ctx *gin.Context) {
	// Parse và validate ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug("Invalid book ID format"))
	}

	var requestBodyData bookmodel.UpdateBookRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
	}

	// Tạo command
	cmd := bookservice.UpdateBookCommand{
		ID:  id,
		Dto: requestBodyData,
	}

	// Thực thi command
	err = c.updateCmdHdl.Execute(ctx.Request.Context(), &cmd)
	if err != nil {
		panic(err)
	}

	// Trả về response
	ctx.JSON(http.StatusOK, datatype.ResponseSuccess(gin.H{
		"message": "Book updated successfully",
	}))
}
