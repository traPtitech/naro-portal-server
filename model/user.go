package model

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	UserName string `json:"userName,omitempty"  db:"userName"`
	Password string `json:"password,omitempty"  db:"password"`
}

func GetUsersHandler(c echo.Context) error {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to Get Users: %v", err))
	}
	return c.JSON(http.StatusOK, users)
}
