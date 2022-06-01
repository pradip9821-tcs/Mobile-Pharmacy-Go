package services

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func ConnectDB() *sql.DB {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("Connected.")
	return db
}
