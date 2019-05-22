package router

import (
	"github.com/labstack/echo"
	"github.com/sapphi-red/webengineer_naro-_7_server/database"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/posts"
	"net/http"
)

func getPostsHandler(c echo.Context) error {
	var posts []posts.Post
	err := database.Posts.GetPosts(posts)
	if err != nil {
		return return500(c, "getPostsError", err)
	}
	return c.JSON(http.StatusAccepted, posts)
}
