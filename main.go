package main

import (
	"net/http"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})

	e.Start(":12100")
}
