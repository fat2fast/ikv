package userhttpgin

import (
	"net/http"

	usermodel "fat2fast/ikv/modules/user/model"
	usersevice "fat2fast/ikv/modules/user/service"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionCreate(c *gin.Context) {
	var requestBodyData usermodel.User

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call business logic in service
	cmd := usersevice.CreateCommand{Dto: requestBodyData}
	uc.createCmdHdl.Execute(c.Request.Context(), &cmd)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created v1",
	})
}
