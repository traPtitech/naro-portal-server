package profiles

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sqlx.DB
)

type GetFollowingResponseBody []string

func GetFollowingHandler(c echo.Context) error {
	id := c.Param("id")
	var response GetFollowingResponseBody

	err := DB.Select(&response, "SELECT followed_user_id FROM follows WHERE following_user_id=?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "db error")
	}

	if response == nil {
		response = GetFollowingResponseBody{}
	}
	return c.JSON(http.StatusOK, response)
}

type GetFollowedResponseBody []string

func GetFollowedHandler(c echo.Context) error {
	id := c.Param("id")
	var response GetFollowedResponseBody

	err := DB.Select(&response, "SELECT following_user_id FROM follows WHERE followed_user_id=?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "db error")
	}

	if response == nil {
		response = GetFollowedResponseBody{}
	}

	return c.JSON(http.StatusOK, response)
}

func PutFollowedHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")

	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM follows WHERE following_user_id=? AND followed_user_id=?", userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "db error")
	}
	if count > 0 {
		return c.NoContent(http.StatusOK)
	}

	_, err = DB.Exec("INSERT INTO follows (following_user_id, followed_user_id) VALUES (?, ?)", userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "db error")
	}

	return c.NoContent(http.StatusOK)
}

func DeleteFollowedHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")

	_, err := DB.Exec("DELETE FROM follows WHERE following_user_id=? AND followed_user_id=?", userID, id)
	if err != nil {
		println(err)
		return c.JSON(http.StatusInternalServerError, "db error")
	}

	return c.NoContent(http.StatusOK)
}
