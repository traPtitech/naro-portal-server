package main

import (
	"fmt"

	"github.com/Ras96/naro-portal-server/model"
	"github.com/Ras96/naro-portal-server/router"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	store, err := model.NewSession(db)
	if err != nil {
		panic(fmt.Errorf("failed in session constructor:%v", err))
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	router.SetRouting(e)

	_ = e.Start(":4000")
}
