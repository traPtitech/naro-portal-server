package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
)

type UserStatus struct {
	// date形式はどうすればいい？
	username string
	// 聞いたまま書いただけ　あとで確認
	now_date time.Date
}

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/db", func(c echo.Context) error {
		databaseURL := os.Getenv("DATABASE_URL")
		db, err := gorm.Open("postgres", databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		data := UserStatus{}
		db.First(&data)
		// data, _ := db.Query("SELECT * FROM user_status;")
		defer db.Close()
		return c.JSON(http.StatusOK, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	e.Start(":" + port) // ここを前述の通り自分のポートにすること
}
