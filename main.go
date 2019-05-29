package main

import (
	"github.com/labstack/gommon/log"
	"github.com/labstack/echo"
	//"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	//"github.com/WistreHosshii/naro-portal-server/model"
	"github.com/WistreHosshii/naro-portal-server/router"
	"github.com/WistreHosshii/naro-portal-server/model"
	

	//"net/http"
)

func main() {
	err := model.EstablishConnection()
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ping",router.Pong)
	e.POST("/signup",router.PostSignUpHandler)

	e.Start(":12500")
}