package controllers

import (
	"com.tcs.mobile-pharmacy/customer.microservice/services"
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var db *sql.DB = services.ConnectDB()

func GetNearByPharmacy(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type User struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	type StoreAddress struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	type Store struct {
		Id                  int     `json:"id"`
		StoreName           string  `json:"store_name"`
		StoreImage          string  `json:"store_image"`
		PrimaryAddress      string  `json:"primary_address"`
		AdditionAddressInfo string  `json:"addition_address_info"`
		Distance            float64 `json:"distance"`
		UserId              int     `json:"user_id"`
	}

	var storeAddress StoreAddress
	var stores []Store
	var user User

	results, err := db.Query(`Select latitude,longitude from addresses where is_select = 1 and user_id=?`, userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	for results.Next() {
		err := results.Scan(&user.Latitude, &user.Longitude)
		if err != nil {
			fmt.Println(err)
		}
	}

	results, err = db.Query(`Select s.id,s.store_name,s.store_image,s.user_id,a.primary_address,a.addition_address_info,a.latitude,a.longitude from addresses a right join stores s on a.user_id = s.user_id where a.is_select = 1`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for results.Next() {
		var store Store
		err := results.Scan(&store.Id, &store.StoreName, &store.StoreImage, &store.UserId, &store.PrimaryAddress, &store.AdditionAddressInfo, &storeAddress.Latitude, &storeAddress.Longitude)
		if err != nil {
			fmt.Println(err)
		}

		var R float64 = 6371

		φ1 := (user.Latitude * math.Pi) / 180
		φ2 := (storeAddress.Latitude * math.Pi) / 180

		Δφ := ((storeAddress.Latitude - user.Latitude) * math.Pi) / 180
		Δλ := ((storeAddress.Longitude - user.Longitude) * math.Pi) / 180

		a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
			math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

		d := R * c

		store.Distance = math.Round(d*100) / 100

		if store.Distance < 100 {
			stores = append(stores, store)
		}
	}

	sort.Slice(stores, func(i, j int) bool {
		return stores[i].Distance < stores[j].Distance
	})

	utils.SuccessResponseWithCount(c, http.StatusOK, "Get NearBy Pharmacy Successfully!", stores, len(stores))
}

func GetNearByPharmacyV2(c *gin.Context) {

	_, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type User struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	type StoreAddress struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	type Store struct {
		Id                  int     `json:"id"`
		StoreName           string  `json:"store_name"`
		StoreImage          string  `json:"store_image"`
		PrimaryAddress      string  `json:"primary_address"`
		AdditionAddressInfo string  `json:"addition_address_info"`
		Distance            float64 `json:"distance"`
		UserId              int     `json:"user_id"`
	}

	var storeAddress StoreAddress
	var stores []Store
	var user User

	QueryParams := c.Request.URL.Query()

	user.Latitude, _ = strconv.ParseFloat(QueryParams.Get("latitude"), 64)
	user.Longitude, _ = strconv.ParseFloat(QueryParams.Get("longitude"), 64)

	results, err := db.Query(`Select s.id,s.store_name,s.store_image,s.user_id,a.primary_address,a.addition_address_info,a.latitude,a.longitude from addresses a right join stores s on a.user_id = s.user_id where a.is_select = 1`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for results.Next() {
		var store Store
		err := results.Scan(&store.Id, &store.StoreName, &store.StoreImage, &store.UserId, &store.PrimaryAddress, &store.AdditionAddressInfo, &storeAddress.Latitude, &storeAddress.Longitude)
		if err != nil {
			fmt.Println(err)
		}

		var R float64 = 6371

		φ1 := (user.Latitude * math.Pi) / 180
		φ2 := (storeAddress.Latitude * math.Pi) / 180

		Δφ := ((storeAddress.Latitude - user.Latitude) * math.Pi) / 180
		Δλ := ((storeAddress.Longitude - user.Longitude) * math.Pi) / 180

		a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
			math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

		d := R * c

		store.Distance = math.Round(d*100) / 100

		if store.Distance < 100 {
			stores = append(stores, store)
		}
	}

	sort.Slice(stores, func(i, j int) bool {
		return stores[i].Distance < stores[j].Distance
	})

	utils.SuccessResponseWithCount(c, http.StatusOK, "Get NearBy Pharmacy Successfully!", stores, len(stores))
}

func CreatePrescription(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Data struct {
		Id                 int64       `json:"id"`
		Name               string      `json:"name"`
		TextNote           string      `json:"text_note"`
		PrescriptionImages interface{} `json:"prescription_images"`
		Medicines          interface{} `json:"medicines"`
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

	var medicines []Medicine
	var images []Image
	var data Data

	textNote := c.PostForm("text_note")

	insert, err := db.Exec(`INSERT INTO prescriptions (text_note, status, user_id) values (?,0,?)`, textNote, userId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription data insertion failed!", err)
		return
	}
	data.Id, _ = insert.LastInsertId()

	form, _ := c.MultipartForm()
	files := form.File["image"]

	for _, file := range files {
		var image Image

		filename := uuid.NewString() + "-" + strings.ReplaceAll(file.Filename, " ", "_")

		c.SaveUploadedFile(file, "../../images/"+filename)

		image.Url = "images/" + filename
		image.MimeType = file.Header.Get("Content-Type")

		images = append(images, image)

	}

	for i := 0; i < len(images); i++ {
		insert, err := db.Exec(`INSERT INTO prescription_images (url, type, prescription_id) VALUES (?,?,?)`, images[i].Url, images[i].MimeType, data.Id)
		if err != nil {
			fmt.Println(err)
			_, err := db.Exec(`DELETE FROM prescriptions WHERE id=?`, data.Id)
			if err != nil {
				fmt.Println(err)
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription data deletion failed!", err)
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription image data insertion failed!", err)
			return
		}
		images[i].Id, _ = insert.LastInsertId()
	}

	MedicineList := []byte(c.PostForm("medicine"))

	err = json.Unmarshal(MedicineList, &medicines)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Failed to parse json data!", err)
		return
	}

	for i := 0; i < len(medicines); i++ {
		insert, err := db.Exec(`INSERT INTO medicines (name,prescription_id) VALUES (?,?)`, medicines[i].Name, data.Id)
		if err != nil {
			_, err := db.Exec(`DELETE FROM prescriptions WHERE id=?`, data.Id)
			if err != nil {
				fmt.Println(err)
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription data deletion failed!", err)
				return
			}
			for j := 0; j < len(images); j++ {
				err = os.Remove("../../" + images[j].Url)
				if err != nil {
					fmt.Println(err)
					utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription images deletion failed on server!", err)
					return
				}
			}
			fmt.Println(err)
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Medicine data insertion failed!", err)
			return
		}
		medicines[i].Id, _ = insert.LastInsertId()
	}

	_, err = db.Exec(`UPDATE prescriptions SET name = ?`, medicines[0].Name)
	if err != nil {
		fmt.Println(err)
		_, err := db.Exec(`DELETE FROM prescriptions WHERE id=?`, data.Id)
		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription data deletion failed!", err)
			return
		}
		for j := 0; j < len(images); j++ {
			err = os.Remove("../../" + images[j].Url)
			if err != nil {
				fmt.Println(err)
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription images deletion failed on server!", err)
				return
			}
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOCreatePrescription, "Prescription name updation failed!", err)
		return
	}
	data.Name = medicines[0].Name
	data.TextNote = textNote
	data.PrescriptionImages = images
	data.Medicines = medicines

	utils.SuccessResponse(c, http.StatusOK, "Prescription created successfully..", data)
}

func DeletePrescription(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	type Image struct {
		Id  int64  `json:"id"`
		Url string `json:"url"`
	}

	var images []Image

	QueryParams := c.Request.URL.Query()

	prescriptionId := QueryParams.Get("prescription_id")

	if prescriptionId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "prescription_id not provided!", "status": constant.FailedStatus})
		return
	}

	result, _ := db.Query(`SELECT id,url FROM prescription_images WHERE prescription_id=?`, prescriptionId)
	if result.Next() == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "No prescription found for loggedIn user!", "status": constant.FailedStatus})
		return
	}
	for result.Next() {
		var image Image
		result.Scan(&image.Id, &image.Url)
		images = append(images, image)
	}

	delete, err := db.Exec(`DELETE FROM prescriptions WHERE user_id=? and id=?`, userId, prescriptionId)
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedToDeletePrescription, "Prescription data deletion failed!", err)
		return
	}

	AffectedRow, _ := delete.RowsAffected()
	if AffectedRow == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "No prescription found for loggedIn user!", "status": constant.FailedStatus})
		return
	}

	for j := 0; j < len(images); j++ {
		os.Remove("../../" + images[j].Url)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prescription deleted successfully.", "status": constant.SuccessStatus})
}

func GetPrescription(c *gin.Context) {

	userId, role, _ := utils.GetUserId(c)
	IsAuth := utils.Authorization(role)
	if !IsAuth {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only customer can access this routes!", "status": constant.FailedStatus})
		return
	}

	QueryParams := c.Request.URL.Query()

	state := QueryParams.Get("state")
	page, _ := strconv.ParseInt(QueryParams.Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(QueryParams.Get("limit"), 10, 64)

	if state == "" || page == 0 || limit == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "State or page or limit not provided!", "status": constant.FailedStatus})
		return
	}

	offset := (page - 1) * limit

	type Prescriptions struct {
		Id                 int64       `json:"id"`
		Name               string      `json:"name"`
		TextNote           string      `json:"text_note"`
		PrescriptionImages interface{} `json:"prescription_images"`
		Medicines          interface{} `json:"medicines"`
		Status             int64       `json:"status"`
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

	if state == "current" {

		result, err := db.Query(`SELECT id, name, text_note, status,createdAt, updatedAt, user_id FROM prescriptions where user_id=? and status=0 LIMIT ? OFFSET ?`, userId, limit, offset)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
			return
		}
		for result.Next() {
			var prescription Prescriptions

			result.Scan(&prescription.Id, &prescription.Name, &prescription.TextNote, &prescription.Status, &prescription.CreatedAt, &prescription.UpdatedAt, &prescription.UserId)

			var images []Image
			var medicines []Medicine

			img, err := db.Query(`SELECT id, url, type FROM prescription_images WHERE prescription_id=?`, prescription.Id)
			if err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
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
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
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

	} else if state == "past" {

		result, err := db.Query(`SELECT id, name, text_note, status,createdAt, updatedAt, user_id FROM prescriptions where user_id=? and status!=0 LIMIT ? OFFSET ?`, userId, limit, offset)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
			return
		}
		for result.Next() {
			var prescription Prescriptions

			result.Scan(&prescription.Id, &prescription.Name, &prescription.TextNote, &prescription.Status, &prescription.CreatedAt, &prescription.UpdatedAt, &prescription.UserId)

			var images []Image
			var medicines []Medicine

			img, err := db.Query(`SELECT id, url, type FROM prescription_images WHERE prescription_id=?`, prescription.Id)
			if err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
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
				utils.ErrorResponse(c, http.StatusInternalServerError, constant.FailedTOFetchPrescription, "Something went wrong!", err)
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

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error_message": "Please, Enter valid state!", "status": constant.FailedStatus})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Prescription fetched successfully..", prescriptions)
}
