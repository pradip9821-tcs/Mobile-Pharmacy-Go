package controllers

import (
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Checkout(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Body struct {
		QuoteId        int     `json:"quote_id" binding:"required"`
		PaymentMethod  int     `json:"payment_method"  binding:"required"`
		CheckoutType   int     `json:"checkout_type" binding:"required"`
		DeliveryCharge float64 `json:"delivery_charge"`
		Amount         float64 `json:"amount"`
		Status         int64   `json:"status"`
		UserId         int     `json:"user_id"`
		StoreId        int64   `json:"store_id"`
	}

	var body Body

	body.UserId = userId

	err := utils.ParseBody(c, &body)
	if err != nil {
		return
	}

	quote, err := db.Query(`SELECT price, store_id FROM quotes WHERE id=?`, body.QuoteId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to quote!", "Database Error", err)
		return
	}
	for quote.Next() {
		quote.Scan(&body.Amount, &body.StoreId)
	}

	body.Amount = body.DeliveryCharge + body.Amount

	c.JSON(http.StatusOK, body)

}
