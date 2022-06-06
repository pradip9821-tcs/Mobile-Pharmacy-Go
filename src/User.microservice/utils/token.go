package utils

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"os"
	"strconv"
	"strings"
)

func init() {
	gotenv.Load()
}

var db *sql.DB

func ExtractToken(c *gin.Context) string {

	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func TokenValid(c *gin.Context) error {

	tokenString := ExtractToken(c)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return err
	}

	return nil
}

func GetUserId(c *gin.Context) (int, error) {

	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		auth0Id := fmt.Sprint(claims["sub"])

		uid, err := strconv.ParseInt(strings.Split(auth0Id, "|")[1], 10, 32)
		if err != nil {
			spew.Dump(err)
			return 0, err
		}
		return int(uid), nil
	}
	return 0, nil
}
