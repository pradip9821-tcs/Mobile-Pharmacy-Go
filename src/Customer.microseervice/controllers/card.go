package controllers

import (
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/card"
	"net/http"
	"os"
)

func AddCard(c *gin.Context) {

	stripe.Key = os.Getenv("STRIPE_SK")

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Body struct {
		CVC      string `json:"cvc" binding:"required"`
		ExpMonth string `json:"exp_month" binding:"required"`
		ExpYear  string `json:"exp_year"  binding:"required"`
		Name     string `json:"name"  binding:"required"`
		Number   string `json:"number"  binding:"required"`
		StripeID string `json:"stripe_id"`
	}

	type Response struct {
		CardId string `json:"card_id"`
	}

	var body Body

	err := utils.ParseBody(c, &body)
	if err != nil {
		return
	}

	user, err := db.Query(`SELECT stripe_id FROM users WHERE id=?`, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Card creation failed!", "Database Error", err)
		return
	}
	for user.Next() {
		user.Scan(&body.StripeID)
	}

	params := &stripe.CardParams{
		Customer: stripe.String(body.StripeID),
		CVC:      stripe.String(body.CVC),
		ExpMonth: stripe.String(body.ExpMonth),
		ExpYear:  stripe.String(body.ExpYear),
		Name:     stripe.String(body.Name),
		Number:   stripe.String(body.Number),
	}
	car, _ := card.New(params)

	_, err = db.Exec(`INSERT INTO cards (card_id,user_id) VALUES (?,?)`, car.ID, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Card creation failed!", "Data insertion failed!", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Card Added Successfully..", gin.H{"card_id": car.ID})
}

func GetCards(c *gin.Context) {

	stripe.Key = os.Getenv("STRIPE_SK")

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Card struct {
		Id       string `json:"id"`
		Brand    string `json:"brand"`
		Country  string `json:"country"`
		CVCCheck string `json:"cvc_check"`
		ExpMonth uint8  `json:"exp_month"`
		ExpYear  uint16 `json:"exp_year"`
		Name     string `json:"name"`
		Last4    string `json:"last_4"`
	}

	var StripeID string
	var cardsData []Card

	user, err := db.Query(`SELECT stripe_id FROM users WHERE id=?`, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cards!", "Database Error", err)
		return
	}
	for user.Next() {
		user.Scan(&StripeID)
	}

	cardId, err := db.Query(`SELECT card_id FROM cards WHERE user_id=?`, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cards!", "Database Error", err)
		return
	}
	for cardId.Next() {

		var id string
		var cards Card

		cardId.Scan(&id)

		params := &stripe.CardParams{
			Customer: stripe.String(StripeID),
		}
		car, _ := card.Get(
			id,
			params,
		)

		cards.Id = car.ID
		cards.Brand = string(car.Brand)
		cards.CVCCheck = string(car.CVCCheck)
		cards.ExpYear = car.ExpYear
		cards.ExpMonth = car.ExpMonth
		cards.Country = car.Country
		cards.Name = car.Name
		cards.Last4 = car.Last4

		cardsData = append(cardsData, cards)
	}

	utils.SuccessResponse(c, http.StatusOK, "Cards fetched successfully..", cardsData)
}
