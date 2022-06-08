package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	card "github.com/stripe/stripe-go/v72/card"
	"net/http"
	"os"
)

func Create(c *gin.Context) {

	stripe.Key = os.Getenv("STRIPE_SK")

	params := &stripe.CardParams{
		Customer: stripe.String("cus_LpwXww1HUq5psz"),
		CVC:      stripe.String("123"),
		ExpMonth: stripe.String("04"),
		ExpYear:  stripe.String("2024"),
		Name:     stripe.String("don"),
		Number:   stripe.String("4242424242424242"),
	}
	car, _ := card.New(params)

	c.JSON(http.StatusOK, car)
}
