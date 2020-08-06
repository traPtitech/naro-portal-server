package model

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func EstablishConnection() (*sqlx.DB, error) {
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "password"
	}

	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "twitterclone"
	}

	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4")
	db = _db
	return db, err
}
