package router

import (
	"github.com/labstack/echo"
)

func CreateRoutes(e *echo.Echo) {
	g := e.Group("")
	g.Use(checkLogin)

	g.GET("/whoami", getWhoAmIHandler)

	g.GET("/users/:id", getUserHandler)

	g.GET("/posts", getPostsHandler)
	g.POST("/posts", createPostsHandler)
}
