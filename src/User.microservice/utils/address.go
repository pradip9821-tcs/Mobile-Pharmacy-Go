package utils

import (
	"com.tcs.mobile-pharmacy/user.microservice/services"
	"com.tcs.mobile-pharmacy/user.microservice/utils/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsSelect(c *gin.Context, user_id int, method string) error {

	db = services.ConnectDB()

	type Address struct {
		Id       int64 `json:"id"`
		IsSelect int   `json:"is_select"`
	}

	var address Address
	var err error

	if method == "Set0" {
		sqlStatement := `SELECT id, is_select FROM addresses WHERE user_id=? and is_select = 1`
		row := db.QueryRow(sqlStatement, user_id)
		err = row.Scan(&address.Id, &address.IsSelect)
		update := `UPDATE addresses SET is_select=0 where id=?`
		row = db.QueryRow(update, address.Id)
	}
	if method == "Set1" {
		sqlStatement := `SELECT id, is_select FROM addresses WHERE user_id=? and is_select = 0`
		row := db.QueryRow(sqlStatement, user_id)
		err = row.Scan(&address.Id, &address.IsSelect)
		update := `UPDATE addresses SET is_select=1 where id=?`
		row = db.QueryRow(update, address.Id)
	}

	db.Close()

	if err != nil {
		fmt.Println(err)
		if err.Error() == constant.NoDataFound {
			return nil
		}
		RespondWithError(c, constant.DatabaseError, constant.BadRequestError, err.Error(), constant.InternalError, http.StatusInternalServerError)
		return err
	}
	return nil
}
