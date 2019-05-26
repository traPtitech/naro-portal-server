package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"

	"github.com/jmoiron/sqlx"

	"github.com/naro-portal-server/account"
	"github.com/naro-portal-server/favo"
	"github.com/naro-portal-server/pin"
	"github.com/naro-portal-server/timeline"
	"github.com/naro-portal-server/tweet"
)

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.POST("/login", account.PostLoginHandler)
	e.POST("/signup", account.PostSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(account.CheckLogin)
	withLogin.POST("/tweet",tweet.PostTweetHandler)
	withLogin.POST("/pin",pin.PostPinHandler)
	withLogin.GET("/timeline/:userName",timeline.GetTimeLineHandler)
	withLogin.POST("/favoAdd",favo.PostAddFavoHandler)
	withLogin.POST("/favoDelete",favo.PostDeleteFavoHandler)

	e.Start(":11401")
}