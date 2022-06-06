package main

import (
	"com.tcs.mobile-pharmacy/user.microservice/controllers"
	"com.tcs.mobile-pharmacy/user.microservice/middlewares"
	"com.tcs.mobile-pharmacy/user.microservice/services"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
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

	router.GET("/get-profile", controllers.GetProfile)

	router.PUT("/update-profile", controllers.UpdateProfile)

	router.POST("/upload", controllers.UploadImage)

	router.POST("/add-address", controllers.AddAddress)

	router.GET("/get-address", controllers.GetAddress)

	router.PUT("/update-address", controllers.UpdateAddress)

	router.DELETE("/delete-address", controllers.DeleteAddress)

	router.PUT("/change-password", controllers.ChangePassword)

	router.Run(":8081")

}
