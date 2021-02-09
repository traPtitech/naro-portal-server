package model

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/srinathgs/mysqlstore"
)

var (
	db *sqlx.DB
)

// EstablishConnection DBと接続
func EstablishConnection() (*sqlx.DB, error) {
	user := os.Getenv("DB_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "pass"
	}

	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "custom_theme"
	}

	_db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4")
	db = _db

	return db, err
}

func NewSession(db *sqlx.DB) (store *mysqlstore.MySQLStore, err error) {
	store, err = mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	return
}
