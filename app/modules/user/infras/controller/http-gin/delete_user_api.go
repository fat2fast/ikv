package userhttpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uc *UserHTTPController) ActionDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User deleted v1"})
}
