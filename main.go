package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/oribe1115/phan-sns-server/handler"
	"github.com/oribe1115/phan-sns-server/model"
)

func main() {
	_, err := model.EstablishConecction()
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	store, err := model.StoreForSession()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/create/tabele/userstatus", handler.CreateUserStatusHandler)
	e.POST("/signup", handler.SignUpHandler)
	e.POST("/login", handler.LoginHandler)

	// withLogin := e.Group("")
	// withLogin.Use(handler.CheckLogin)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	e.Start(":" + port) // ここを前述の通り自分のポートにすること
}
