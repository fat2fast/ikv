package userhttpgin

import (
	"net/http"

	usermodel "fat2fast/ikv/modules/user/model"
	userservice "fat2fast/ikv/modules/user/service"
	"fat2fast/ikv/shared/datatype"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionAuthenticate(c *gin.Context) {
	var requestBodyData usermodel.LoginForm

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err))
	}
	cmd := &userservice.AuthenticateCommand{
		Dto: requestBodyData,
	}
	result, err := uc.authenticateCmdHdl.Execute(c, cmd)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(result))
}
