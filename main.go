package main

import (
	"kuragate-server/auths"
	"kuragate-server/dbs"
	"kuragate-server/messages"
	"kuragate-server/profiles"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := dbs.GetDB()
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	messages.DB = db
	auths.DB = db
	profiles.DB = db

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/isvalidid/:reqID", auths.GetIsValidIDHandler)
	e.POST("/signup", auths.PostSignUpHandler)
	e.POST("/login", auths.PostLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(auths.CheckLogin)

	withLogin.GET("/whoami", auths.GetWhoAmIHandler)
	withLogin.POST("/logout", auths.PostLogoutHandler)

	withLogin.POST("/messages", messages.PostMessageHandler)
	withLogin.GET("/messages", messages.GetMassagesHandler)
	withLogin.GET("/messages/:id", messages.GetSingleMassageHandler)
	withLogin.PUT("/messages/:id/fav", messages.PutMessageFavHandler)
	withLogin.DELETE("/messages/:id/fav", messages.DeleteMessageFavHandler)

	withLogin.GET("/profiles/:id/followed", profiles.GetFollowedHandler)
	withLogin.PUT("/profiles/:id/followed", profiles.PutFollowedHandler)
	withLogin.DELETE("/profiles/:id/followed", profiles.DeleteFollowedHandler)
	withLogin.GET("/profiles/:id/following", profiles.GetFollowingHandler)

	e.Start(":13300")
}
