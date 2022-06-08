package controllers

import (
	"com.tcs.mobile-pharmacy/pharmacist.microservice/services"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils"
	"com.tcs.mobile-pharmacy/pharmacist.microservice/utils/constant"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var db *sql.DB = services.ConnectDB()

func GetRequests(c *gin.Context) {

	_, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only pharmacist can access this routes!", "status": constant.FailedStatus})
		return
	}

	QueryParams := c.Request.URL.Query()

	page, _ := strconv.ParseInt(QueryParams.Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(QueryParams.Get("limit"), 10, 64)

	if page == 0 || limit == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "page or limit not provided!", "status": constant.FailedStatus})
		return
	}

	offset := (page - 1) * limit

	type Prescriptions struct {
		UserName           string      `json:"user_name"`
		UserPicture        string      `json:"user_picture"`
		Id                 int64       `json:"id"`
		Name               string      `json:"name"`
		TextNote           string      `json:"text_note"`
		PrescriptionImages interface{} `json:"prescription_images"`
		Medicines          interface{} `json:"medicines"`
		Status             int64       `json:"status"`
		TotalQuotes        int64       `json:"total_quotes"`
		CreatedAt          string      `json:"created_at"`
		UpdatedAt          string      `json:"updated_at"`
		UserId             int64       `json:"user_id"`
	}

	type Image struct {
		Id       int64  `json:"id"`
		Url      string `json:"url"`
		MimeType string `json:"mime_type"`
	}

	type Medicine struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	var prescriptions []Prescriptions

	result, err := db.Query(`SELECT id, name, text_note, status,createdAt, updatedAt, user_id FROM prescriptions where status=0 LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchRequests, "Something went wrong!", err)
		return
	}
	for result.Next() {
		var prescription Prescriptions

		result.Scan(&prescription.Id, &prescription.Name, &prescription.TextNote, &prescription.Status, &prescription.CreatedAt, &prescription.UpdatedAt, &prescription.UserId)

		var images []Image
		var medicines []Medicine

		quotes, err := db.Query(`SELECT COUNT(id) FROM quotes WHERE prescription_id=?`, prescription.Id)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchRequests, "Something went wrong!", err)
			return
		}
		for quotes.Next() {
			quotes.Scan(&prescription.TotalQuotes)
		}

		user, err := db.Query(`SELECT name,picture FROM users WHERE id=?`, prescription.UserId)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchRequests, "Something went wrong!", err)
			return
		}
		for user.Next() {
			user.Scan(&prescription.UserName, &prescription.UserPicture)
		}

		img, err := db.Query(`SELECT id, url, type FROM prescription_images WHERE prescription_id=?`, prescription.Id)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchRequests, "Something went wrong!", err)
			return
		}
		for img.Next() {
			var image Image

			img.Scan(&image.Id, &image.Url, &image.MimeType)

			images = append(images, image)
		}
		prescription.PrescriptionImages = images

		med, err := db.Query(`SELECT id, name FROM medicines WHERE prescription_id=?`, prescription.Id)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchRequests, "Something went wrong!", err)
			return
		}
		for med.Next() {
			var medicine Medicine

			med.Scan(&medicine.Id, &medicine.Name)

			medicines = append(medicines, medicine)
		}
		prescription.Medicines = medicines

		prescriptions = append(prescriptions, prescription)
	}

	utils.SuccessResponse(c, http.StatusOK, "All requests fetched successfully..", prescriptions)

}

func AddQuotes(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only pharmacist can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Quote struct {
		Id             int64   `json:"id"`
		Price          float32 `json:"price" binding:"required"`
		TextNote       string  `json:"text_note" binding:"required"`
		PrescriptionId int64   `json:"prescription_id" binding:"required"`
		StoreName      string  `json:"store_name"`
		StoreId        int64   `json:"store_id"`
	}

	var quote Quote

	if err := c.BindJSON(&quote); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to post quotes", err.Error(), err)
		return
	}

	store := db.QueryRow(`SELECT id,store_name FROM stores WHERE user_id=?`, userId)
	_ = store.Scan(&quote.StoreId, &quote.StoreName)

	insert, err := db.Exec(`INSERT INTO quotes (store_name, price, text_note, store_id, prescription_id) VALUES (?,?,?,?,?)`, quote.StoreName, quote.Price, quote.TextNote, quote.StoreId, quote.PrescriptionId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to post quotes", "Something went wrong!", err)
		return
	}
	quote.Id, _ = insert.LastInsertId()

	utils.SuccessResponse(c, http.StatusOK, "Quote added successfully.", quote)

}
