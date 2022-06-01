package controllers

import (
	"com.tcs.mobile-pharmacy/user.microservice/models"
	"com.tcs.mobile-pharmacy/user.microservice/services"
	"com.tcs.mobile-pharmacy/user.microservice/utils"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"database/sql"
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
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
		Name        string `json:"name"`
		Email       string `json:"email"`
		Gender      int    `json:"gender"`
		Picture     string `json:"picture"`
		CountryCode string `json:"country_code"`
		Phone       int    `json:"phone"`
	}

	var reqBody Body
	var user models.User

	userId, _, _ := utils.GetUserId(c)

	sqlGetStat := `SELECT name, email, gender, picture, country_code, phone FROM users WHERE id=?`
	row := db.QueryRow(sqlGetStat, userId)
	err := row.Scan(&reqBody.Name, &reqBody.Email, &reqBody.Gender, &reqBody.Picture, &reqBody.CountryCode, &reqBody.Phone)
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

	_, err = db.Exec("UPDATE users SET name = ? where id=?", reqBody.Name, userId)
	if err != nil {
		fmt.Println("Data update failed!")
		return
	}

	sqlStatement := `SELECT name, email, gender, picture, country_code, phone FROM users WHERE id=?`
	row = db.QueryRow(sqlStatement, userId)
	err = row.Scan(&user.Name, &user.Email, &user.Gender, &user.Picture, &user.CountryCode, &user.Phone)
	if err != nil {
		fmt.Println(err)
		if err.Error() == constant.NoDataFound {
			utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
			return
		}
		utils.RespondWithError(c, constant.DatabaseError, constant.BadRequestError, err.Error(), constant.InternalError, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data updated successfully.", "user": reqBody, "status": constant.SuccessStatus})
}

func UploadImage(c *gin.Context) {

	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	file, data, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed to upload",
		})
	}

	result, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		PublicID: strings.Split(data.Filename, ".")[0],
	})

	if err != nil {
		c.String(http.StatusConflict, "Upload to cloudinary failed")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Successfully uploaded the file",
		"secureURL": result.SecureURL,
		"publicURL": result.URL,
	})
}

func AddAddress() {

}
