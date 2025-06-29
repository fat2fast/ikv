package bookhttpgin

import (
	"net/http"

	bookservice "fat2fast/ikv/modules/book/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActionDeleteBook xóa book - DELETE /:id
func (c *BookHTTPController) ActionDeleteBook(ctx *gin.Context) {
	// Parse và validate ID
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug("Invalid book ID format"))
	}

	// Get delete type from query param (soft/hard)
	deleteType := ctx.DefaultQuery("type", "soft")
	isSoftDelete := deleteType != "hard"

	// Tạo command
	cmd := bookservice.DeleteBookCommand{
		ID:   id,
		Soft: isSoftDelete,
	}

	// Thực thi command
	err = c.deleteCmdHdl.Execute(ctx.Request.Context(), &cmd)
	if err != nil {
		panic(err)
	}

	// Tạo message phù hợp
	message := "Book deleted successfully"
	if isSoftDelete {
		message = "Book soft deleted successfully"
	} else {
		message = "Book permanently deleted successfully"
	}

	// Trả về response
	ctx.JSON(http.StatusOK, datatype.ResponseSuccess(gin.H{
		"message": message,
	}))
}
