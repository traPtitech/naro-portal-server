package router

import (
	"net/http"

	"github.com/Ras96/naro-portal-server/model"
	"github.com/labstack/echo/v4"
)

func SetRouting(e *echo.Echo) {
	e.GET("/ping", getpingHandler)


	e.POST("/signup", model.PostSignUpHandler)
	e.POST("/login", model.PostLoginHandler)

	withlogin := e.Group("")
	{
		withlogin.Use(model.CheckLogin)
		withlogin.GET("/pingping", getpingHandler)

		users := withlogin.Group("/users")
		{
			users.GET("", model.GetUsersHandler)
		}
	}
}


func getpingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}