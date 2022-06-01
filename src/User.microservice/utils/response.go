package utils

import (
	"com.tcs.mobile-pharmacy/user.microservice/models"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondWithError(c *gin.Context, message, err, code, description string, statusCode int) {
	var error models.Error
	error.Message = message
	error.Error = err
	error.Code = code
	error.ErrorDescription = description
	if error.Error == constant.NilString {
		error.Error = constant.BadRequestError
	}
	if error.Code == "" {
		error.Code = constant.EmptyData
	}
	if error.ErrorDescription == constant.NilString {
		error.ErrorDescription = constant.EmptyData
	}
	if error.ErrorDescription == "Unauthorized" {
		statusCode = http.StatusUnauthorized
	}
	error.Status = constant.FailedStatus
	c.JSON(statusCode, error)
	return
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
