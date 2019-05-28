package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"

	"github.com/naro-portal-server/model"
)

func main() {
	model.Establish()

	store, err := mysqlstore.NewMySQLStoreFromConnection(model.Db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.POST("/login", model.PostLoginHandler)
	e.POST("/signup", model.PostSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(model.CheckLogin)
	withLogin.POST("/tweet", model.PostTweetHandler)
	withLogin.POST("/pin", model.PostPinHandler)
	withLogin.POST("/pinDelete",model.PostDeletePinHandler)
	withLogin.GET("/timeline/:userName", model.GetTimeLineHandler)
	withLogin.GET("/pin/:userName", model.GetPinHandler)
	withLogin.POST("/favoAdd", model.PostAddFavoHandler)
	withLogin.POST("/favoDelete", model.PostDeleteFavoHandler)
	withLogin.POST("/isFavo",model.GetIsFavoHandler)

	e.Start(":11401")
}
