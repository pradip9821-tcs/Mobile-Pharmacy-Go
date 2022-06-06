package utils

import (
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"github.com/gin-gonic/gin"
)

func SuccessResponseWithCount(c *gin.Context, statuscode int, message string, data interface{}, count int) {

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Count   int         `json:"count"`
		Status  int         `json:"status"`
	}{
		Message: message,
		Data:    data,
		Count:   count,
		Status:  constant.SuccessStatus,
	}

	c.JSON(statuscode, response)
}
