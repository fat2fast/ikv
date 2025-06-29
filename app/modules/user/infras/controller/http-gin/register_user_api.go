package userhttpgin

import (
	"net/http"

	usermodel "fat2fast/ikv/modules/user/model"
	usersevice "fat2fast/ikv/modules/user/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionRegister(c *gin.Context) {
	var requestBodyData usermodel.RegisterForm

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
	}

	// call business logic in service
	cmd := usersevice.CreateCommand{Dto: requestBodyData}
	user, err := uc.createCmdHdl.Execute(c.Request.Context(), &cmd)
	if err != nil {
		panic(datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error()))
	}
	c.JSON(http.StatusCreated, datatype.ResponseSuccess(user))
}
