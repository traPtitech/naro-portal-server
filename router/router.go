package router

import (
	"net/http"

	"github.com/Ras96/naro-portal-server/model"
	"github.com/labstack/echo/v4"
)

func SetRouting(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	users := e.Group("/users")
	{
		users.GET("", model.GetUsersHandler)
	}
}
