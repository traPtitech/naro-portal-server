package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/srinathgs/mysqlstore"
)

// LoadDatabase データベースを読み込みます
func LoadDatabase() *sqlx.DB {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	fmt.Println("Connected!")

	return db
}

// SetupSessionDatabase DBからセッションストアを作成します
func SetupSessionDatabase(db *sqlx.DB) *mysqlstore.MySQLStore {
	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}
	return store
}
