package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"github.com/WistreHosshii/naro-portal-server/model"
	"github.com/WistreHosshii/naro-portal-server/router"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	err := model.EstablishConnection()
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	if os.Getenv("NODE_ENV") = "development" {
		port:="12500"
	} else {
		port := os.Getenv("PORT")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ping", router.Pong)
	e.POST("/signup", router.PostSignUpHandler)
	e.POST("/login", PostLoginHandler)

	e.Start(":" + port)
}
