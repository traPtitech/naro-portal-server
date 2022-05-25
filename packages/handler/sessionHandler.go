package handler

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Me struct {
	UserName string `json:"username,omitempty"  db:"UserName"`
	UserID   string `json:"userID,omitempty"  db:"UserID"`
}
type LoginRequestBody struct {
	UserName string `json:"username,omitempty" form:"UserName"`
	Password string `json:"password,omitempty" form:"Password"`
}

type User struct {
	UserName   string `json:"username,omitempty"  db:"UserName"`
	UserID     string `json:"userID,omitempty"  db:"UserID"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

func GetUserNameHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userName := c.Get("UserName").(string)
		return c.String(http.StatusOK, userName)
	}
}

func PostSignUpHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		// もう少し真面目にバリデーションするべき
		if req.Password == "" || req.UserName == "" {
			// エラーは真面目に返すべき
			return c.String(http.StatusBadRequest, "項目が空です")
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
		}

		// ユーザーの存在チェック
		var count int

		err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE UserName=?", req.UserName)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}

		if count > 0 {
			return c.String(http.StatusConflict, "ユーザーが既に存在しています")
		}

		_, err = db.Exec("INSERT INTO users (UserName, HashedPass) VALUES (?, ?)", req.UserName, hashedPass)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}

		user := User{}
		err = db.Get(&user, "SELECT * FROM users WHERE UserName=?", req.UserName)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}

		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		sess.Values["UserName"] = user.UserName
		sess.Values["UserID"] = user.UserID
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusCreated)
	}
}

func PostLoginHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		user := User{}
		err := db.Get(&user, "SELECT * FROM users WHERE UserName=?", req.UserName)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db errror: %v", err))
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
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
		sess.Values["UserName"] = req.UserName
		sess.Values["UserID"] = user.UserID
		sess.Save(c.Request(), c.Response())

		return c.String(http.StatusOK, "login")
	}
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["UserName"] == nil || sess.Values["UserID"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("UserName", sess.Values["UserName"].(string))
		c.Set("UserID", sess.Values["UserID"].(string))

		return next(c)
	}
}
func GetWhoAmIHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("sessions", c)

		return c.JSON(http.StatusOK, Me{
			UserName: sess.Values["UserName"].(string),
			UserID:   sess.Values["UserID"].(string),
		})
	}
}
func LogoutHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("sessions", c)
		sess.Values["UserName"] = nil
		sess.Values["UserID"] = nil
		sess.Save(c.Request(), c.Response())
		return c.String(http.StatusOK, "logout")
	}
}
