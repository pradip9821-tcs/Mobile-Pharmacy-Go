package main

import (
	"com.tcs.mobile-pharmacy/customer.microservice/controllers"
	"com.tcs.mobile-pharmacy/customer.microservice/middlewares"
	"com.tcs.mobile-pharmacy/customer.microservice/services"
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
		c.JSON(http.StatusOK, gin.H{"Message": "Welcome to Customer microservice of Mobile Pharmacy App..", "Status": constant.SuccessStatus})
	})

	router.Use(middlewares.SetMiddlewareAuthentication())

	router.GET("/get-nearby-pharmacy", controllers.GetNearByPharmacy)

	router.GET("/get-nearby-pharmacy/v2", controllers.GetNearByPharmacyV2)

	router.POST("/create-prescription", controllers.CreatePrescription)

	router.DELETE("/delete-prescription", controllers.DeletePrescription)

	router.GET("/get-prescription", controllers.GetPrescription)

	router.POST("/add-card", controllers.AddCard)

	router.GET("/get-cards", controllers.GetCards)

	router.POST("/checkout", controllers.Checkout)

	router.Run(":8082")
}
