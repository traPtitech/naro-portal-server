package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type SignupRequestBody struct {
	ID       string `json:"id,omitempty" from:"id"`
	Name     string `json:"name,omitempty" from:"name"`
	Password string `json:"password,omitempty" from:"password"`
}

type LoginRequestBody struct {
	ID         string `json:"id,omitempty" from:"id"`
	HashedPass string `json:"hashed_pass,omitempty" from:"hashed_pass"`
}

type User struct {
	ID         string `json:"id,omitempty" db:"id"`
	Name       string `json:"name,omitempty" db:"name"`
	HashedPass string `json:"hashed_pass,omitempty" db:"hashed_pass"`
}

func main() {
	_db, err := sqlx.Connect(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOSTNAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Start(":13300")
}
