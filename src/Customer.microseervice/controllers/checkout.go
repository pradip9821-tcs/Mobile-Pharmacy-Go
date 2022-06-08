package controllers

import (
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"net/http"
	"os"
)

func Checkout(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Body struct {
		PaymentMethod  int     `json:"payment_method"`
		CheckoutType   int     `json:"checkout_type" binding:"required"`
		DeliveryCharge float64 `json:"delivery_charge"`
		Amount         float64 `json:"amount"`
		Status         int64   `json:"status"`
		QuoteId        int     `json:"quote_id" binding:"required"`
		UserId         int     `json:"user_id"`
		StoreId        int64   `json:"store_id"`
		CardId         string  `json:"card_id"`
	}

	type Response struct {
		Id            int64       `json:"id"`
		OrderData     interface{} `json:"order_data"`
		PaymentIntent string      `json:"payment_intent"`
		PaymentId     int64       `json:"payment_id"`
	}

	var body Body
	var response Response

	body.UserId = userId

	err := utils.ParseBody(c, &body)
	if err != nil {
		return
	}

	if body.PaymentMethod != 0 && body.PaymentMethod != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid payment method!", "status": 0})
		return
	}
	if body.CheckoutType > 4 || body.CheckoutType < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid checkout type!", "status": 0})
		return
	}
	if body.PaymentMethod == 0 && body.CardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Card_id is required!", "status": 0})
		return
	}

	result, err := db.Query(`SELECT price, store_id FROM quotes WHERE id=?`, body.QuoteId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Failed to fetch quote!", err)
		return
	}

	if !result.Next() {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Quote not found!", "status": 0})
		return
	}

	var stripeId string

	user, err := db.Query(`SELECT stripe_id FROM users WHERE id=?`, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Failed to fetch quote!", err)
		return
	}
	for user.Next() {
		user.Scan(&stripeId)
	}

	quote, err := db.Query(`SELECT price, store_id FROM quotes WHERE id=?`, body.QuoteId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Failed to fetch quote!", err)
		return
	}
	for quote.Next() {
		quote.Scan(&body.Amount, &body.StoreId)
	}

	body.Amount = body.DeliveryCharge + body.Amount

	order, err := db.Exec(`INSERT INTO orders (payment_method, checkout_type, delivery_charge, amount, user_id, store_id, quote_id) VALUES
			(?,?,?,?,?,?,?)`, body.PaymentMethod, body.CheckoutType, body.DeliveryCharge, body.Amount, body.UserId, body.StoreId, body.QuoteId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Order data insertion failed", err)
		return
	}

	if body.PaymentMethod == 1 {
		response.Id, _ = order.LastInsertId()
		response.OrderData = body

		//payment, err := db.Exec(`INSERT INTO payments (payment_intent,transaction_id, amount, user_id, store_id, order_id) VALUES
		//			(?,?,?,?,?,?)`, "OFFLINE", "OFFLINE", body.Amount, userId, body.StoreId, response.Id)
		//if err != nil {
		//	fmt.Println(err)
		//	utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Payment creation failed", err)
		//	_, err := db.Exec(`DELETE FROM orders WHERE id=?`, response.Id)
		//	if err != nil {
		//		fmt.Println(err)
		//		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Order deletion failed!", err)
		//		return
		//	}
		//	return
		//}
		//response.PaymentId, _ = payment.LastInsertId()

		_, err = db.Exec(`UPDATE orders SET status=1 WHERE id=?`, response.Id)
		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Order Updation Failed", "Database Error!", err)
			return
		}

		utils.SuccessResponse(c, http.StatusOK, "Order Created Successfully..", response)
		return
	}

	stripe.Key = os.Getenv("STRIPE_SK")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(body.Amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Description:   stripe.String("Mobile-Pharmacy Payment"),
		Customer:      stripe.String(stripeId),
		PaymentMethod: stripe.String(body.CardId),
	}
	pi, _ := paymentintent.New(params)

	response.Id, _ = order.LastInsertId()
	response.OrderData = body
	response.PaymentIntent = pi.ClientSecret

	payment, err := db.Exec(`INSERT INTO payments (payment_intent, amount, user_id, store_id, order_id) VALUES 
					(?,?,?,?,?)`, pi.ClientSecret, body.Amount, userId, body.StoreId, response.Id)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Payment creation failed", err)
		_, err := db.Exec(`DELETE FROM orders WHERE id=?`, response.Id)
		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Order creation failed!", "Order deletion failed!", err)
			return
		}
		return
	}
	response.PaymentId, _ = payment.LastInsertId()

	utils.SuccessResponse(c, http.StatusOK, "Order Created Successfully..", response)

}

func ChangePaymentStatus(c *gin.Context) {

	type Body struct {
		OrderId       int    `json:"order_id" binding:"required"`
		PaymentId     int    `json:"payment_id" binding:"required"`
		TransactionId string `json:"transaction_id"`
		Status        string `json:"status" binding:"required"`
	}

	var body Body

	err := utils.ParseBody(c, &body)
	if err != nil {
		return
	}

	if body.Status != "SUCCESS" && body.Status != "FAILED" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Invalid status!", "status": 0})
		return
	}

	if body.Status == "SUCCESS" && body.TransactionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "Transaction_id is required!", "status": 0})
		return
	}

	_, err = db.Exec(`UPDATE payments SET transaction_id=?,status=? WHERE id=?`, body.TransactionId, body.Status, body.PaymentId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Payment Updation Failed", "Database Error!", err)
		return
	}

	_, err = db.Exec(`UPDATE orders SET status=1 WHERE id=?`, body.OrderId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Order Updation Failed", "Database Error!", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment updated successfully..", body)

}

func GetOrderStatus(c *gin.Context) {

	_, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	QueryParams := c.Request.URL.Query()

	orderId := QueryParams.Get("order_id")

	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_message": "order_id is required in query parameter!", "status": 0})
		return
	}

	order, err := db.Query(`SELECT status FROM orders WHERE id=?`, orderId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed To Fetch Order!", "Database Error!", err)
		return
	}
	var status int
	for order.Next() {
		order.Scan(&status)
	}

	utils.SuccessResponse(c, http.StatusOK, "Order status fetched successfully..", gin.H{"order_status": status})
}
