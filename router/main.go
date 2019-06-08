package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"

	"github.com/naro-portal-server/model"
)

func main() {
	err:=model.Establish()
	if err != nil {
		panic(err)
	}

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
	withLogin.POST("/logout",model.PostLogoutHandler)
	withLogin.POST("/tweet", model.PostTweetHandler)
	withLogin.POST("/pin", model.PostPinHandler)
	withLogin.DELETE("/pin", model.DeletePinHandler)
	withLogin.GET("/timeline/:userName", model.GetTimelineHandler)
	withLogin.GET("/timelinePin/:userName", model.GetPinHandler)
	withLogin.POST("/favo", model.PostFavoHandler)
	withLogin.DELETE("/favo", model.DeleteFavoHandler)
	withLogin.GET("/isFavo/:tweetID", model.GetIsFavoHandler)
	withLogin.GET("/whoAmI", model.GetWhoAmIHandler)
	withLogin.GET("/reloadTimeline/:userName",model.GetIsReloadTimelineHandler)

	e.Start(":11400")
}
