package main

import (
	"github.com/labstack/echo"
	//"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	//"github.com/WistreHosshii/naro-portal-server/model"
	"github.com/WistreHosshii/naro-portal-server/router"

	//"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ping",router.pong)

	e.Start(":12500")
}