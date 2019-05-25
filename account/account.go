package account

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"github.com/pborman/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type LoginRequestBody struct {
	UserName string
	UserPassword string
}

type User struct{
	UserName			string	``
	UserPassword		string
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := db.Get(&user, "SELECT (user_name,user_password) FROM user WHERE user_name=?", req.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
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
	sess.Values["userName"] = req.UserName
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	var userID uuid.UUID
	db.Get(&userID,"SELECT user_ID FROM user WHERE user_Name=?",req.UserName)
	if userID!=nil {
		return c.String(http.StatusBadRequest, "ユーザーが既に存在しています")
	}

	if utf8.RuneCountInString(req.UserPassword)<8{
		return c.String(http.StatusBadRequest, "パスワードは8文字以上です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	_, err = db.Exec("INSERT INTO user (Username, HashedPass) VALUES (?, ?)", req.UserName, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}