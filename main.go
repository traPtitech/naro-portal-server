package main

import (
	"github.com/labstack/echo/v4"
	"github.com/Ras96/naro-portal-server/router"
)

func main() {
	e := echo.New()
	router.SetRouting(e)

	e.Start(":4000")
}
