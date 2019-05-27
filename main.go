package main

import (
	"net/http"

	"github.com/motoki317/naro-portal-server/database"
	"github.com/motoki317/naro-portal-server/router"

	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := database.LoadDatabase()
	store := database.SetupSessionDatabase(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	router.SetupLoginRoutes(e, db)
	router.SetupWithLoginRoutes(e, db)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server successfully started!")
	})

	e.Start(":10900")
}
