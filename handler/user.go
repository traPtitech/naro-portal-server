package handler

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/oribe1115/phan-sns-server/model"
)

func CreateUserStatusHandler(c echo.Context) error {
	model.CreateUserStatusTable()
	return c.String(http.StatusOK, "user_status table crated!\n")
}
