package userhttpgin

import (
	userservice "fat2fast/ikv/modules/user/service"
	"fat2fast/ikv/shared/datatype"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (uc *UserHTTPController) ActionGetDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	query := userservice.GetDetailsQuery{Id: id}

	user, err := uc.getDetailQryHdl.Execute(c.Request.Context(), &query)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
	c.JSON(http.StatusOK, gin.H{"message": "User detail v1"})
}
