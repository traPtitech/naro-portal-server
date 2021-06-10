package main

import (
	"kuragate-server/auths"
	"kuragate-server/dbs"
	"kuragate-server/messages"
	"kuragate-server/profiles"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbs.InitDb()

	store, err := mysqlstore.NewMySQLStoreFromConnection(dbs.Db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/isvalidid/:reqID", auths.GetIsValidIDHandler)
	e.POST("/signup", auths.PostSignUpHandler)
	e.POST("/login", auths.PostLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(auths.CheckLogin)

	withLogin.GET("/whoami", auths.GetWhoAmIHandler)
	withLogin.GET("/logout", auths.GetLogoutHandler)

	withLogin.POST("/messages", messages.PostMessageHandler)
	withLogin.GET("/messages", messages.GetMassagesHandler)
	withLogin.GET("/messages/:id", messages.GetSingleMassageHandler)
	withLogin.PUT("/messages/:id/fav", messages.PutMessageFavHandler)
	withLogin.DELETE("/messages/:id/fav", messages.DeleteMessageFavHandler)

	withLogin.GET("/profiles/:id/followed", profiles.GetFollowdHandler)
	withLogin.PUT("/profiles/:id/followed", profiles.PutFollowdHandler)
	withLogin.DELETE("/profiles/:id/followed", profiles.DeleteFollowdHandler)
	withLogin.GET("/profiles/:id/following", profiles.GetFollowingHandler)

	e.Start(":13300")
}
