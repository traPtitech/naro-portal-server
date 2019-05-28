package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/sapphi-red/webengineer_naro-_7_server/database"
	"github.com/sapphi-red/webengineer_naro-_7_server/router"
	"net/http"
)

func main() {
	database.ConnectDB()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowCredentials: true,
	  }))
	e.Use(middleware.Logger())
	e.Use(session.Middleware(database.SessionStore))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	router.CreateRoutes(e)
	router.CreateLoginRoutes(e)

	e.Start(":12100")
}
