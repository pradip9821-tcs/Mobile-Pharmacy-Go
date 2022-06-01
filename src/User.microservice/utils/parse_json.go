package utils

import (
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseBody(c *gin.Context, reqBody interface{}) error {

	if err := c.BindJSON(&reqBody); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"ErrorMessage": err.Error(), "Status": constant.FailedStatus})
		return err
	}
	return nil
}
