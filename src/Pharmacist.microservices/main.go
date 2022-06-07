package main

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/controllers"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/middlewares"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()

	router.Use(middlewares.CorsMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Welcome to Pharmacist microservice of Mobile Pharmacy App..", "Status": constant.SuccessStatus})
	})

	router.Use(middlewares.SetMiddlewareAuthentication())

	router.GET("/get-requests", controllers.GetRequests)

	router.POST("/add-quote", controllers.AddQuotes)

	router.Run(":8083")
}
