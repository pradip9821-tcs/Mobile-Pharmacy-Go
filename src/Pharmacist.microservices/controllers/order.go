package controllers

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetOfflineOrder(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only pharmacist can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Order struct {
		Id             int64   `json:"id"`
		PaymentMethod  int64   `json:"payment_method"`
		Status         int64   `json:"status"`
		CheckoutType   int64   `json:"checkout_type"`
		DeliveryCharge float64 `json:"delivery_charge"`
		Amount         float64 `json:"amount"`
		CreatedAt      string  `json:"created_at"`
		UserId         int64   `json:"user_id"`
		StoreId        int64   `json:"store_id"`
		QuoteId        int64   `json:"quote_id"`
	}

	var storeId int64
	var orders []Order

	user, err := db.Query(`SELECT id FROM stores WHERE user_id=?`, userId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch store_id!", "Something went wrong!", err)
		return
	}
	for user.Next() {
		user.Scan(&storeId)
	}

	order, err := db.Query(`SELECT id, payment_method, status, checkout_type, delivery_charge, amount, createdAt, user_id, store_id, quote_id FROM orders WHERE store_id=? and payment_method=1 ORDER BY id DESC`, storeId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch order!", "Something went wrong!", err)
		return
	}
	for order.Next() {
		var ord Order
		order.Scan(&ord.Id, &ord.PaymentMethod, &ord.Status, &ord.CheckoutType, &ord.DeliveryCharge, &ord.Amount, &ord.CreatedAt, &ord.UserId, &ord.StoreId, &ord.QuoteId)
		orders = append(orders, ord)
	}

	utils.SuccessResponse(c, http.StatusOK, "Orders fetched successfully..", orders)

}

func CollectPaymentOffline(c *gin.Context) {

	_, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only pharmacist can access this routes!", "status": constant.FailedStatus})
		return
	}

	QueryParams := c.Request.URL.Query()

	orderId := QueryParams.Get("order_id")

	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "order_id is required in query parameter!", "status": 0})
		return
	}

	order, err := db.Exec(`UPDATE orders SET status=3 WHERE id=?`, orderId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch order!", "Something went wrong!", err)
		return
	}
	count, _ := order.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Order not found!", "status": 0})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment collected successfully..", "status": 1})
}

func ChangeOrderStatus(c *gin.Context) {

	_, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only pharmacist can access this routes!", "status": constant.FailedStatus})
		return
	}

	QueryParams := c.Request.URL.Query()

	orderId := QueryParams.Get("order_id")
	status, _ := strconv.ParseInt(QueryParams.Get("status"), 10, 32)

	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "order_id is required in query parameter!", "status": 0})
		return
	}
	if status < 2 || status > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid status in query parameter!", "status": 0})
		return
	}

	order, err := db.Exec(`UPDATE orders SET status=? WHERE id=?`, status, orderId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch order!", "Something went wrong!", err)
		return
	}
	count, _ := order.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "May be order not found or order is up-to date!", "status": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status changed successfully..", "status": 1})
}
