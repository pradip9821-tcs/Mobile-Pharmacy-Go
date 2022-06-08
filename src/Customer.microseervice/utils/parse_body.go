package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseBody(c *gin.Context, reqBody interface{}) error {

	if err := c.BindJSON(&reqBody); err != nil {
		fmt.Println(err)
		ErrorResponse(c, http.StatusBadRequest, "Failed To Parse Body", err.Error(), nil)
		return err
	}
	return nil
}
