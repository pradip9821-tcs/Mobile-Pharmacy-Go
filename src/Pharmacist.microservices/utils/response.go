package utils

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
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

func SuccessResponse(c *gin.Context, statuscode int, message string, data interface{}) {

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Status  int         `json:"status"`
	}{
		Message: message,
		Data:    data,
		Status:  constant.SuccessStatus,
	}

	c.JSON(statuscode, response)
}

func ErrorResponse(c *gin.Context, statuscode int, message string, description string, err error) {
	response := struct {
		ErrorMessage string `json:"error_message"`
		Error        error  `json:"error"`
		Description  string `json:"description"`
		Status       int    `json:"status"`
	}{
		ErrorMessage: message,
		Error:        err,
		Description:  description,
		Status:       constant.FailedStatus,
	}

	c.JSON(statuscode, response)
}
