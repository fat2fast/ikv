package userhttpgin

import (
	"net/http"

	usermodel "fat2fast/ikv/modules/user/model"
	userservice "fat2fast/ikv/modules/user/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActionUpdateProfile xử lý PUT /profile/:id - Cập nhật thông tin profile user
func (c *UserHTTPController) ActionUpdateProfile(ctx *gin.Context) {

	// Parse user ID
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithError("Invalid user ID format"))
	}

	// Bind JSON request body
	var requestData usermodel.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithError("Invalid request data"))
	}

	// Tạo command
	cmd := &userservice.UpdateProfileCommand{
		UserID: userID,
		Dto:    requestData,
	}

	// Thực thi command
	err = c.updateProfileCmdHdl.Execute(ctx.Request.Context(), cmd)
	if err != nil {
		panic(err)
	}

	// Trả về thành công
	ctx.JSON(http.StatusOK, gin.H{"data": true})
}
