package model

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

//User Userの構造体
type User struct {
	UserName     string `json:"userName,omitempty"  db:"name"`
	UserPassword string `json:"userPassword,omitempty"  db:"password"`
}

//PostLoginHandler POST /login ログイン
func PostLoginHandler(c echo.Context) error {
	req := User{}
	c.Bind(&req)

	user := User{}
	err := Db.Get(&user, "SELECT name,password FROM user WHERE name=?", req.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(req.UserPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["UserName"] = req.UserName
	var userID string
	Db.Get(&userID,"SELECT ID FROM user WHERE name=?",req.UserName)
	sess.Values["UserID"]=userID
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

//CheckLogin ログイン確認
func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["UserName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("UserName", sess.Values["UserName"].(string))

		return next(c)
	}
}

//PostSignUpHandler Post /signup サインアップ
func PostSignUpHandler(c echo.Context) error {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	Db = _db

	req := User{}
	c.Bind(&req)

	var userID string
	Db.Get(&userID, "SELECT ID FROM user WHERE name=?", req.UserName)
	if userID != "" {
		return c.String(http.StatusBadRequest, "ユーザーが既に存在しています")
	}

	if utf8.RuneCountInString(req.UserPassword) < 8 {
		return c.String(http.StatusBadRequest, "パスワードは8文字以上です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	_, err = Db.Exec("INSERT INTO (user name,ID,password) VALUES (?, ?,?)", req.UserName, uuid.New(), hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
