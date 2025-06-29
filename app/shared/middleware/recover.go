package middleware

import (
	"fat2fast/ikv/shared/datatype"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type CanGetStatusCode interface {
	StatusCode() int
}

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			isProduction := os.Getenv("ENV") == "prod" || os.Getenv("GIN_MODE") == "release"

			if r := recover(); r != nil {
				if appError, ok := r.(CanGetStatusCode); ok {
					c.JSON(appError.StatusCode(), appError)

					if !isProduction {
						log.Printf("Error: %+v", appError)
						panic(r)
					}
					return
				}

				appError := datatype.ErrInternalServerError

				if isProduction {
					c.JSON(appError.StatusCode(), appError.WithDebug(""))
				} else {
					c.JSON(appError.StatusCode(), appError.WithDebug(fmt.Sprintf("%s", r)))
					panic(r)
				}
			}
		}()

		c.Next()
	}
}
