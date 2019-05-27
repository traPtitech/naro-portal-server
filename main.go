package main

import (
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"

	"github.com/oribe1115/phan-sns-server/handler"
)

var (
	db *gorm.DB
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	_db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		panic("failed to connect database")
	}
	db = _db
	defer db.Close()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/create", handler.CreateUserStatusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	e.Start(":" + port) // ここを前述の通り自分のポートにすること
}
