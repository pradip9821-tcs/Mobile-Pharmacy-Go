package middlewares

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetMiddlewareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": constant.Unauthorized, "status": constant.FailedStatus})
			c.Abort()
			return
		}
		c.Next()
	}
}
