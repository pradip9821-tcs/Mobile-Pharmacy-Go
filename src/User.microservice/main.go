package main

import (
	"com.tcs.mobile-pharmacy/user.microservice/controllers"
	"com.tcs.mobile-pharmacy/user.microservice/middlewares"
	"com.tcs.mobile-pharmacy/user.microservice/services"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {

	db = services.ConnectDB()
	router := gin.Default()

	router.Use(middlewares.CorsMiddleware())
	router.Use(middlewares.SetMiddlewareAuthentication())

	router.GET("/get-profile", controllers.GetProfile)

	router.POST("/update-profile", controllers.UpdateProfile)

	router.POST("/upload", controllers.UploadImage)

	router.POST("/add-address", controllers.AddAddress)

	router.Run(":8081")

}
