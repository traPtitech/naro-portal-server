package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"

	"git.trapti.tech/mazrean/twitter_clone_server/model"
)

func main() {
	err := model.Establish()
	if err != nil {
		panic(err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(model.Db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	err=model.Create()
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
	withLogin.POST("/logout", model.PostLogoutHandler)
	withLogin.POST("/tweet", model.PostTweetHandler)
	withLogin.POST("/pin", model.PostPinHandler)
	withLogin.DELETE("/pin", model.DeletePinHandler)
	withLogin.GET("/isPin/:tweetID", model.GetIsPinHandler)
	withLogin.GET("/timeline/:userName", model.GetTimelineHandler)
	withLogin.GET("/timelinePin/:userName", model.GetPinHandler)
	withLogin.GET("/timelineFavo/:userName", model.GetFavoHandler)
	withLogin.POST("/favo", model.PostFavoHandler)
	withLogin.DELETE("/favo", model.DeleteFavoHandler)
	withLogin.GET("/isFavo/:tweetID", model.GetIsFavoHandler)
	withLogin.GET("/whoAmI", model.GetWhoAmIHandler)
	withLogin.GET("/userName", model.GetUserListHandler)
	withLogin.GET("/reloadTimeline/:userName", model.GetIsReloadTimelineHandler)

	e.Start(":11400")
}
