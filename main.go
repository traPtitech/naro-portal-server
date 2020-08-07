package main

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/purewhite404/naro-server/model"
	"github.com/purewhite404/naro-server/router"
	"github.com/srinathgs/mysqlstore"
)

func main() {
	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))
	e.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:8000"},
		}),
	)

	e.POST("/register", router.PostRegisterHandler)
	e.POST("/login", router.PostLoginHandler)
	e.GET("/tweet", router.GetTweetHandler)

	withLogin := e.Group("")
	withLogin.Use(router.HasLoggedin)
	withLogin.GET("/whoami", router.GetMeHandler)
	withLogin.GET("/timeline", router.GetTweetHandler)
	withLogin.POST("/tweet", router.PostTweetHandler)

	e.Start(":11900")
}
