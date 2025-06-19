package userhttpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionGetList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User list v1"})
}
