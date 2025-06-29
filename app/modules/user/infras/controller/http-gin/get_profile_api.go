package userhttpgin

import (
	"net/http"

	userservice "fat2fast/ikv/modules/user/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActionGetProfile xử lý GET /profile/:id - Lấy thông tin profile user
func (c *UserHTTPController) ActionGetProfile(ctx *gin.Context) {

	// Parse user ID
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithError("Invalid user ID format"))
	}

	// Tạo query
	query := &userservice.GetProfileQuery{
		UserID: userID,
	}

	// Thực thi query
	profile, err := c.getProfileQryHdl.Execute(ctx.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	// Trả về thành công
	ctx.JSON(http.StatusOK, datatype.ResponseSuccess(profile))
}
