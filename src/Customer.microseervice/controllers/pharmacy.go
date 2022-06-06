package controllers

import (
	"com.tcs.mobile-pharmacy/customer.microservice/services"
	"com.tcs.mobile-pharmacy/customer.microservice/utils"
	"com.tcs.mobile-pharmacy/customer.microservice/utils/constant"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sort"
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
