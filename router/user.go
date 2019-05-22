package router

import (
	"github.com/labstack/echo"
	"github.com/sapphi-red/webengineer_naro-_7_server/database"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/users"
	"net/http"
)

func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	var user users.User
	err := database.Users.GetUser(id, &user)
	if err != nil {
		return return500(c, "getPostsError", err)
	}
	return c.JSON(http.StatusAccepted, user)
}
