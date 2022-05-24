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
	Username string `json:"username,omitempty"  db:"username"`
}
type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

func GetUserNameHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userName := c.Get("userName").(string)
		return c.String(http.StatusOK, userName)
	}
}

func PostSignUpHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		// もう少し真面目にバリデーションするべき
		if req.Password == "" || req.Username == "" {
			// エラーは真面目に返すべき
			return c.String(http.StatusBadRequest, "項目が空です")
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
		}

		// ユーザーの存在チェック
		var count int

		err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", req.Username)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}

		if count > 0 {
			return c.String(http.StatusConflict, "ユーザーが既に存在しています")
		}

		_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", req.Username, hashedPass)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}

		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}
		sess.Values["userName"] = req.Username
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusCreated)
	}
}

func PostLoginHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := LoginRequestBody{}
		c.Bind(&req)

		user := User{}
		err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
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
		sess.Values["userName"] = req.Username
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	}
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
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
func GetWhoAmIHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("sessions", c)

		return c.JSON(http.StatusOK, Me{
			Username: sess.Values["userName"].(string),
		})
	}
}
func LogoutHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("sessions", c)
		sess.Values["userName"] = nil
		sess.Save(c.Request(), c.Response())
		return c.String(http.StatusOK, "logout")
	}
}
