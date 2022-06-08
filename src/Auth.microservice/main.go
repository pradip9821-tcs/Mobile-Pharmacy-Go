package main

import (
	"com.tcs.mobile-pharmacy/auth.microservice/controllers"
	"com.tcs.mobile-pharmacy/auth.microservice/middlewares"
	"com.tcs.mobile-pharmacy/auth.microservice/services"
	"com.tcs.mobile-pharmacy/auth.microservice/utils/constant"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"net/http"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {

	db = services.ConnectDB()

	router := gin.Default()

	router.Use(middlewares.CorsMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Welcome to Auth microservice of Mobile Pharmacy App..", "Status": constant.SuccessStatus})
	})

	router.POST("/login", controllers.Login)

	router.POST("/signup", controllers.Register)

	router.POST("/refresh-token", controllers.RefreshToken)

	router.POST("/forgot-password", controllers.ForgotPassword)

	router.POST("/create", controllers.Create)

	router.Run(":8080")
}
