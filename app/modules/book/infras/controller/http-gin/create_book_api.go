package bookhttpgin

import (
	"net/http"

	bookmodel "fat2fast/ikv/modules/book/model"
	bookservice "fat2fast/ikv/modules/book/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
)

// ActionCreateBook xử lý tạo book mới - POST /
func (c *BookHTTPController) ActionCreateBook(ctx *gin.Context) {
	var requestBodyData bookmodel.CreateBookRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
	}

	// Tạo command
	cmd := bookservice.CreateBookCommand{Dto: requestBodyData}

	// Thực thi command
	response, err := c.createCmdHdl.Execute(ctx.Request.Context(), &cmd)
	if err != nil {
		panic(err)
	}

	// Trả về response
	ctx.JSON(http.StatusCreated, datatype.ResponseSuccess(response))
}
