package main

import (
	"com.tcs.mobile-pharmacy/customer.microservice/middlewares"
	"com.tcs.mobile-pharmacy/customer.microservice/services"
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

var db *sql.DB

func main() {

	db = services.ConnectDB()

	router := gin.Default()

	router.Use(middlewares.CorsMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Welcome to User microservice of Mobile Pharmacy App..", "Status": constant.SuccessStatus})
	})

	router.Use(middlewares.SetMiddlewareAuthentication())

	router.GET("/protected", func(c *gin.Context) {
		userId, role, _ := utils.GetUserId(c)

		IsAuth := utils.Authorization(role)
		if !IsAuth {
			c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Message": "Welcome to User microservice of Mobile Pharmacy App..",
			"Status":  constant.SuccessStatus,
			"UserId":  userId,
			"Role":    role,
		})
	})

	router.Run(":8082")
}
