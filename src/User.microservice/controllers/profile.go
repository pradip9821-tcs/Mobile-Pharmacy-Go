package controllers

import (
	"com.tcs.mobile-pharmacy/user.microservice/models"
	"com.tcs.mobile-pharmacy/user.microservice/services"
	"com.tcs.mobile-pharmacy/user.microservice/utils"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var db *sql.DB = services.ConnectDB()

func GetProfile(c *gin.Context) {

	type Response struct {
		Message string `json:"message"`
		User    models.User
		Status  int `json:"status"`
	}

	var response Response
	var user models.User

	userId, _, err := utils.GetUserId(c)

	if err != nil {
		fmt.Println("Can't get user id!")
		return
	}

	sqlStatement := `SELECT id, name, email, role, gender, picture, country_code, phone, is_test, is_active, is_verify, is_delete, createdAt, updatedAt FROM users WHERE id=?`
	row := db.QueryRow(sqlStatement, userId)
	err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Gender, &user.Picture, &user.CountryCode, &user.Phone, &user.IsTest, &user.IsActive, &user.IsVerify, &user.IsDelete, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		fmt.Println(err)
		if err.Error() == constant.NoDataFound {
			utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
			return
		}
		utils.RespondWithError(c, constant.DatabaseError, constant.BadRequestError, err.Error(), constant.InternalError, http.StatusInternalServerError)
		return
	}

	response.Message = constant.GetProfileSuccess
	response.User = user
	response.Status = constant.SuccessStatus

	c.JSON(http.StatusOK, response)

}

func UpdateProfile(c *gin.Context) {

	type Body struct {
		Id          int    `json:"id"`
		Role        int    `json:"role"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		Gender      int    `json:"gender"`
		Picture     string `json:"picture"`
		CountryCode string `json:"country_code"`
		Phone       int    `json:"phone"`
	}

	var reqBody Body

	userId, _, _ := utils.GetUserId(c)

	sqlGetStat := `SELECT id, role,name, email, gender, picture, country_code, phone FROM users WHERE id=?`
	row := db.QueryRow(sqlGetStat, userId)
	err := row.Scan(&reqBody.Id, &reqBody.Role, &reqBody.Name, &reqBody.Email, &reqBody.Gender, &reqBody.Picture, &reqBody.CountryCode, &reqBody.Phone)
	if err != nil {
		fmt.Println(err)
		if err.Error() == constant.NoDataFound {
			utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
			return
		}
		utils.RespondWithError(c, constant.DatabaseError, constant.BadRequestError, err.Error(), constant.InternalError, http.StatusInternalServerError)
		return
	}

	utils.ParseBody(c, &reqBody)

	_, err = db.Exec(`UPDATE users SET name=?,email=?,gender=?,picture=?,country_code=?,phone=? where id=?`, reqBody.Name, reqBody.Email, reqBody.Gender, reqBody.Picture, reqBody.CountryCode, reqBody.Phone, userId)
	if err != nil {
		fmt.Println("Data update failed!")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile updated successfully", reqBody)
}
