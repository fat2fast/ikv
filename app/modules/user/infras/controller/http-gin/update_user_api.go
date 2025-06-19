package userhttpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User updated v1"})
}
