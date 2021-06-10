package profiles

import (
	"fmt"
	"kuragate-server/dbs"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
)

type GetFollowingResponseBody []string

func GetFollowingHandler(c echo.Context) error {
	id := c.Param("id")
	var response GetFollowingResponseBody

	err := dbs.Db.Select(&response, "SELECT followed_user_id FROM follows WHERE following_user_id=?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if response == nil {
		response = GetFollowingResponseBody{}
	}
	return c.JSON(http.StatusOK, response)
}

type GetFollowdResponseBody []string

func GetFollowdHandler(c echo.Context) error {
	id := c.Param("id")
	var response GetFollowdResponseBody

	err := dbs.Db.Select(&response, "SELECT following_user_id FROM follows WHERE followed_user_id=?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if response == nil {
		response = GetFollowdResponseBody{}
	}

	return c.JSON(http.StatusOK, response)
}

func PutFollowdHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")

	var count int
	err := dbs.Db.Get(&count, "SELECT COUNT(*) FROM follows WHERE following_user_id=? AND followed_user_id=?", userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	if count > 0 {
		return c.NoContent(http.StatusOK)
	}

	_, err = dbs.Db.Exec("INSERT INTO follows (following_user_id, followed_user_id) VALUES (?, ?)", userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func DeleteFollowdHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	id := c.Param("id")

	_, err := dbs.Db.Exec("DELETE FROM follows WHERE following_user_id=? AND followed_user_id=?", userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}
