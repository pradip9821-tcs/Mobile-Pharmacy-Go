package main

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "Welcome to Pharmacist microservice of Mobile Pharmacy App..", "Status": constant.SuccessStatus})
	})

	router.Run(":8083")
}
