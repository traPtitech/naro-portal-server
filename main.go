package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"

	"github.com/motoki317/github-webhook/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := database.LoadDatabase()
	store := database.SetupSessionDatabase(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server successfully started!")
	})

	e.Start(":10900")
}
