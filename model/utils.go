package model

import (
	"os"
	"github.com/labstack/gommon/log"
	"github.com/jmoiron/sqlx"
	"github.com/WistreHosshii/naro-portal-server/router"

	"fmt"

)

var (
	db *sqlx.DB
)

//dbとの通信
func EstablishConnection() error {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	router.Db = _db
	
	return nil 
}