package controllers

import (
	"com.tcs.mobile-pharmacy/user.microservice/utils"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ChangePassword(c *gin.Context) {

	type Body struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	var body Body
	var UserPasswoed string

	userId, _ := utils.GetUserId(c)
	utils.ParseBody(c, &body)

	sqlStatement := `SELECT password FROM users WHERE id=?`
	row := db.QueryRow(sqlStatement, userId)
	err := row.Scan(&UserPasswoed)
	if err != nil {
		if err.Error() == constant.NoDataFound {
			utils.SuccessResponse(c, http.StatusOK, "User not found!", nil)
			return
		}
		utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(UserPasswoed), []byte(body.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Password!", "status": 0})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)

	fmt.Println(string(hashedPassword))

	_, err = db.Exec(`UPDATE users SET password=? where id=?`, string(hashedPassword), userId)
	if err != nil {
		// need to central error handling
		fmt.Println(err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully.", body)

}
