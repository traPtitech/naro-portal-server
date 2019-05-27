package router

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// SetupTweetRotues Tweet関連のRouteを置きます
func SetupTweetRoutes(e *echo.Echo, db *sqlx.DB) {
	withLogin := e.Group("/tweets")
	withLogin.POST("/create", makePostCreateTweet(db))
	withLogin.POST("/delete/:id", makePostDeleteTweet(db))
}

func makePostCreateTweet(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		panic("Implement me")
	}
}

func makePostDeleteTweet(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		panic("Implement me")
	}
}
