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

func SignUpHandler(c echo.Context) error {
	userData := model.DataForSignUp{}
	c.Bind(&userData)
	model.AddNewUserStatus(userData)

	return c.String(http.StatusOK, "Succeded")
}
func LoginHandler(c echo.Context) error {

}
