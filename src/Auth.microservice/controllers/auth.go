package controllers

import (
	"com.tcs.mobile-pharmacy/auth.microservice/models"
	"com.tcs.mobile-pharmacy/auth.microservice/services"
	"com.tcs.mobile-pharmacy/auth.microservice/utils"
	"com.tcs.mobile-pharmacy/auth.microservice/utils/constant"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"net/http"
	"os"
	"strings"
)

func init() {
	gotenv.Load()
}

var db *sql.DB

func Register(c *gin.Context) {

	type Response struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var error models.Error
	var response Response

	db = services.ConnectDB()

	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role")
	gender := c.PostForm("gender")
	country_code := c.PostForm("country_code")
	phone := c.PostForm("phone")
	//store_name := c.PostForm("store_name")
	//license_id := c.PostForm("license_id")

	str := `email=` + email +
		`&password=` + password +
		`&name=` + name +
		`&client_id=` + os.Getenv("CLIENT_ID") +
		`&connection=` + os.Getenv("CONNECTION")

	body := utils.APIHandler(constant.Auth0SignupAPI, str)

	_ = json.Unmarshal([]byte(string(body)), &error)

	if error.Code != constant.NilString {
		utils.RespondWithError(c, constant.FailedToSignup, error.Error, error.Code, error.ErrorDescription, http.StatusBadRequest)
		return
	}

	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	file, data, err := c.Request.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed to upload",
		})
		return
	}

	result, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		PublicID: strings.Split(data.Filename, ".")[0],
	})

	if err != nil {
		c.String(http.StatusConflict, "Upload to cloudinary failed")
		return
	}

	picture := result.URL

	fmt.Println(string(body))
	_ = json.Unmarshal([]byte(string(body)), &response)
	fmt.Println(response)

	update, err := db.Query(`UPDATE users SET role = ?, gender = ?, country_code = ?, phone = ?, picture = ? WHERE id=?`, role, gender, country_code, phone, picture, response.Id)

	if err != nil {
		utils.RespondWithError(c, constant.DatabaseError, constant.BadRequestError, constant.EmptyData, constant.InternalError, http.StatusInternalServerError)
		return
	}

	//update, err = db.Query(`INSERT INTO `)

	defer update.Close()

	utils.SuccessResponse(c, http.StatusOK, constant.SignupSuccess, response)
}

func Login(c *gin.Context) {

	type Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Message      string `json:"message"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		User         models.User
		Status       int `json:"status"`
	}

	var response Response
	var error models.Error
	var reqBody Body
	var user models.User

	utils.ParseBody(c, &reqBody)

	db = services.ConnectDB()

	str := `grant_type=password` +
		`&username=` + reqBody.Email +
		`&password=` + reqBody.Password +
		`&scope=offline_access` +
		`&audience=` + os.Getenv("AUDIENCE") +
		`&client_id=` + os.Getenv("CLIENT_ID") +
		`&client_secret=` + os.Getenv("CLIENT_SECRET")

	body := utils.APIHandler(constant.Auth0GetTokenAPI, str)

	_ = json.Unmarshal([]byte(string(body)), &error)

	if error.Error != constant.NilString {
		utils.RespondWithError(c, constant.FailedToLogin, error.Error, error.Code, error.ErrorDescription, http.StatusBadRequest)
		return
	}

	sqlStatement := `SELECT id, name, email, role, gender, picture, country_code, phone, is_test, is_active, is_verify, is_delete, createdAt, updatedAt FROM users WHERE email=?`
	row := db.QueryRow(sqlStatement, reqBody.Email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Gender, &user.Picture, &user.CountryCode, &user.Phone, &user.IsTest, &user.IsActive, &user.IsVerify, &user.IsDelete, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		utils.RespondWithError(c, constant.DatabaseError, constant.BadRequestError, constant.EmptyData, constant.InternalError, http.StatusInternalServerError)
		return
	}

	_ = json.Unmarshal([]byte(string(body)), &response)
	response.Message = constant.LoginSuccess
	response.Status = constant.SuccessStatus
	response.User = user
	c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {

	type Body struct {
		RefreshToken string `json:"refresh_token"`
	}

	type Response struct {
		Message     string `json:"message"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Status      int    `json:"status"`
	}

	var reqBody Body
	var error models.Error
	var response Response

	utils.ParseBody(c, &reqBody)

	str := `grant_type=refresh_token` +
		`&client_id=` + os.Getenv("CLIENT_ID") +
		`&client_secret=` + os.Getenv("CLIENT_SECRET") +
		`&refresh_token=` + reqBody.RefreshToken

	body := utils.APIHandler(constant.Auth0GetTokenAPI, str)

	_ = json.Unmarshal([]byte(string(body)), &error)

	if error.Error != constant.NilString {
		utils.RespondWithError(c, constant.FailedToFetchToken, error.Error, error.Code, error.ErrorDescription, http.StatusBadRequest)
		return
	}

	_ = json.Unmarshal([]byte(string(body)), &response)

	response.Message = constant.GetTokenSuccess
	response.Status = constant.SuccessStatus
	c.JSON(http.StatusOK, response)

}

func ForgotPassword(c *gin.Context) {

	type Body struct {
		Email string `json:"email"`
	}

	var reqBody Body

	utils.ParseBody(c, &reqBody)

	db = services.ConnectDB()

	sqlStatement := `SELECT email FROM users WHERE email=?`
	row := db.QueryRow(sqlStatement, reqBody.Email)
	err := row.Scan(&reqBody.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.RespondWithError(c, constant.UserNotFound, constant.InvalidAccess, constant.EmptyData, "Entered email not found in database!", http.StatusNotFound)
			return
		}
		utils.RespondWithError(c, constant.DatabaseError, err.Error(), constant.EmptyData, constant.InternalError, http.StatusInternalServerError)
		return
	}

	str := `username=` + reqBody.Email +
		`&client_id=` + os.Getenv("CLIENT_ID") +
		`&connection=` + os.Getenv("CONNECTION")

	body := utils.APIHandler(constant.Auth0ForgotPasswordAPI, str)

	fmt.Println(string(body))

	if string(body) == `{"error":"connection is required."}` {
		utils.RespondWithError(c, constant.InternalError, "connection is required.", constant.EmptyData, constant.BadRequestError, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": constant.SendLinkSuccess, "status": constant.SuccessStatus})

}
