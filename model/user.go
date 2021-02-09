package model

import (
	"fmt"
	"net/http"

	"github.com/Ras96/naro-portal-server/domain"
	"github.com/labstack/echo/v4"
)

func GetUsersHandler(c echo.Context) error {
	users := []domain.User{}
	err := db.Select(&users, "SELECT userName FROM users")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to Get Users: %v", err))
	}
	return c.JSON(http.StatusOK, users)
}
