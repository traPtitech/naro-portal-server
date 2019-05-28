package main

import (
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"

	"github.com/oribe1115/phan-sns-server/handler"
	"github.com/oribe1115/phan-sns-server/model"
)

func main() {
	model.EstablishConecction()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/create/tabele/userstatus", handler.CreateUserStatusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	e.Start(":" + port) // ここを前述の通り自分のポートにすること
}
