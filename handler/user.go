package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/oribe1115/phan-sns-server/model"
)

func CreateUserStatusHandler(c echo.Context) error {
	model.CreateUserStatusTable()
	return c.String(http.StatusOK, "user_status table crated!\n")
}

func SignUpHandler(c echo.Context) error {
	userData := model.DataForSignUp{}
	c.Bind(&userData)
	err := model.AddNewUserStatus(userData)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Faild to Sign Up")
	}

	return c.String(http.StatusOK, "Succeded")
}
func LoginHandler(c echo.Context) error {
	loginData := model.LoginRequestBody{}
	c.Bind(&loginData)
	userID, err := model.Login(loginData)
	if err != nil {
		fmt.Println(err)
		if err == model.ErrWrongPass {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	sess.Values["userID"] = userID
	sess.Save(c.Request(), c.Response())

	return c.String(http.StatusOK, "Login Succeded")
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userID"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userID", sess.Values["userID"].(string))

		return next(c)

	}
}
