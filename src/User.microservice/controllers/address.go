package controllers

import (
	"com.tcs.mobile-pharmacy/user.microservice/utils"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

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

func AddAddress(c *gin.Context) {

	type Body struct {
		Id                  int64   `json:"id"`
		PrimaryAddress      string  `json:"primary_address" binding:"required"`
		AdditionAddressInfo string  `json:"addition_address_info" binding:"required"`
		AddressType         int     `json:"address_type"`
		Latitude            float64 `json:"latitude" binding:"required"`
		Longitude           float64 `json:"longitude" binding:"required"`
		IsSelect            int     `json:"is_select"`
		UserId              int     `json:"user_id"`
	}

	var body Body

	body.UserId, _, _ = utils.GetUserId(c)

	err := utils.ParseBody(c, &body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body)

	if body.IsSelect == 1 {
		err := utils.IsSelect(c, body.UserId, constant.SetAsNotSelected)
		if err != nil {
			// need to central error handling
			fmt.Println(err)
			return
		}
	}

	insert, err := db.Exec(`INSERT INTO addresses (primary_address, addition_address_info, address_type, latitude, longitude, is_select, user_id)
		values (?,?,?,?,?,?,?)`, body.PrimaryAddress, body.AdditionAddressInfo, body.AddressType, body.Latitude, body.Longitude, body.IsSelect, body.UserId)

	if err != nil {
		err := utils.IsSelect(c, body.UserId, constant.SetAsSelected)
		if err != nil {
			// need to central error handling
			fmt.Println(err)
			return
		}
		// need to central error handling middleware
		fmt.Println("Data update failed!")
		return
	}

	body.Id, _ = insert.LastInsertId()

	utils.SuccessResponse(c, http.StatusOK, "Address added successfully.", body)

}

func GetAddress(c *gin.Context) {

	type Address struct {
		Id                  int64   `json:"id"`
		PrimaryAddress      string  `json:"primary_address" binding:"required"`
		AdditionAddressInfo string  `json:"addition_address_info" binding:"required"`
		AddressType         int     `json:"address_type"`
		Latitude            float64 `json:"latitude" binding:"required"`
		Longitude           float64 `json:"longitude" binding:"required"`
		IsSelect            int     `json:"is_select"`
		UserId              int     `json:"user_id"`
		CreatedAt           string  `json:"createdAt"`
		UpdatedAt           string  `json:"updatedAt"`
	}
	var addresses []Address
	userId, _, _ := utils.GetUserId(c)

	sqlStatement := `SELECT * FROM addresses WHERE user_id=?`
	results, _ := db.Query(sqlStatement, userId)

	for results.Next() {
		var address Address
		err := results.Scan(&address.Id, &address.PrimaryAddress, &address.AdditionAddressInfo, &address.AddressType, &address.Latitude, &address.Longitude, &address.IsSelect, &address.CreatedAt, &address.UpdatedAt, &address.UserId)
		if err != nil {
			fmt.Println(err)
		}
		addresses = append(addresses, address)
	}

	if addresses == nil {
		utils.SuccessResponse(c, http.StatusNotFound, "Address Not Found!", addresses)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Get addresses successfully", addresses)

}

func DeleteAddress(c *gin.Context) {

	type Address struct {
		Id                  int64   `json:"id"`
		PrimaryAddress      string  `json:"primary_address" binding:"required"`
		AdditionAddressInfo string  `json:"addition_address_info" binding:"required"`
		AddressType         int     `json:"address_type"`
		Latitude            float64 `json:"latitude" binding:"required"`
		Longitude           float64 `json:"longitude" binding:"required"`
		IsSelect            int     `json:"is_select"`
		UserId              int     `json:"user_id"`
		CreatedAt           string  `json:"createdAt"`
		UpdatedAt           string  `json:"updatedAt"`
	}
	var address Address

	QueryParams := c.Request.URL.Query()

	addressId := QueryParams.Get("address_id")

	sqlStatement := `SELECT * FROM addresses WHERE id=?`
	row := db.QueryRow(sqlStatement, addressId)
	err := row.Scan(&address.Id, &address.PrimaryAddress, &address.AdditionAddressInfo, &address.AddressType, &address.Latitude, &address.Longitude, &address.IsSelect, &address.CreatedAt, &address.UpdatedAt, &address.UserId)
	if err != nil {
		if err.Error() == constant.NoDataFound {
			utils.SuccessResponse(c, http.StatusOK, "Address Not Found!", nil)
			return
		}
		utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
		return
	}

	sqlStatement = `DELETE FROM addresses WHERE id=?;`
	_ = db.QueryRow(sqlStatement, addressId)

	utils.SuccessResponse(c, http.StatusOK, "Address deleted successfully.", address)

}

func UpdateAddress(c *gin.Context) {

	type Address struct {
		Id                  int64   `json:"id"`
		PrimaryAddress      string  `json:"primary_address" `
		AdditionAddressInfo string  `json:"addition_address_info" `
		AddressType         int     `json:"address_type"`
		Latitude            float64 `json:"latitude" `
		Longitude           float64 `json:"longitude" `
		IsSelect            int     `json:"is_select"`
		UserId              int     `json:"user_id"`
	}

	var address Address

	QueryParams := c.Request.URL.Query()

	addressId := QueryParams.Get("address_id")

	sqlStatement := `SELECT id, primary_address, addition_address_info, address_type, latitude, longitude, is_select, user_id FROM addresses WHERE id=?`
	row := db.QueryRow(sqlStatement, addressId)
	err := row.Scan(&address.Id, &address.PrimaryAddress, &address.AdditionAddressInfo, &address.AddressType, &address.Latitude, &address.Longitude, &address.IsSelect, &address.UserId)
	if err != nil {
		if err.Error() == constant.NoDataFound {
			utils.SuccessResponse(c, http.StatusOK, "Address Not Found!", nil)
			return
		}
		utils.RespondWithError(c, constant.UserNotFound, constant.EmptyData, "404", constant.EmptyData, http.StatusNotFound)
		return
	}

	userId, _, _ := utils.GetUserId(c)

	if userId != address.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": constant.Unauthorized, "status": constant.FailedStatus})
		c.Abort()
		return
	}

	utils.ParseBody(c, &address)

	if address.IsSelect == 1 {
		err := utils.IsSelect(c, address.UserId, constant.SetAsNotSelected)
		if err != nil {
			// need to central error handling
			fmt.Println(err)
			return
		}
	}

	_, err = db.Exec(`UPDATE addresses SET primary_address=?,addition_address_info=?,address_type=?,latitude=?,longitude=?,is_select=? where id=?`, address.PrimaryAddress, address.AdditionAddressInfo, address.AddressType, address.Latitude, address.Longitude, address.IsSelect, addressId)
	if err != nil {
		err := utils.IsSelect(c, address.UserId, constant.SetAsSelected)
		if err != nil {
			// need to central error handling
			fmt.Println(err)
			return
		}
		// need to central error handling middleware
		fmt.Println("Data update failed!")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Address Get successfully.", address)
}
