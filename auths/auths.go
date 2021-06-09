package auths

import (
	"errors"
	"fmt"
	"kuragate-server/dbs"
	"net/http"
	"regexp"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string `json:"id,omitempty" db:"id"`
	Name       string `json:"name,omitempty" db:"name"`
	HashedPass string `json:"hashed_pass,omitempty" db:"hashed_pass"`
}

type SignUpRequestBody struct {
	ID       string `json:"id,omitempty" from:"id"`
	Name     string `json:"name,omitempty" from:"name"`
	Password string `json:"password,omitempty" from:"password"`
}

type LoginRequestBody struct {
	ID       string `json:"id,omitempty" from:"id"`
	Password string `json:"password,omitempty" from:"password"`
}

type IsValidIDResponseBody bool

type WhoAmIResponseBody struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

func check_regexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func isValidId(id string) (bool, error) {
	if len(id) == 0 || len(id) > 20 || !check_regexp(`^[0-9a-zA-Z]+$`, id) {
		return false, nil
	}
	var count int
	err := dbs.Db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", id)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	return true, nil
}

func isValidName(id string) bool {
	return (len(id) != 0 && len(id) <= 30)
}

func isValidPassword(id string) bool {
	return (len(id) != 0 && len(id) <= 30 && check_regexp(`^[0-9a-zA-Z]+$`, id))
}

func GetIsValidIDHandler(c echo.Context) error {
	req := c.QueryParam("reqID")

	isValid, err := isValidId(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, isValid)
}

func PostSignUpHandler(c echo.Context) error {
	req := SignUpRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err.Error()))
	}

	validId, err := isValidId(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	if !validId || !isValidName(req.Name) || !isValidPassword(req.Password) {
		return c.JSON(http.StatusBadRequest, "Bad Request: ID, Name, Passwordのいずれかが不適切です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	_, err = dbs.Db.Exec("INSERT INTO users (id, name, hashed_pass) VALUES (?, ?, ?)", req.ID, req.Name, hashedPass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusCreated)
}

func PostLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	user := User{}
	err = dbs.Db.Get(&user, "SELECT * FROM users WHERE id=?", req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.JSON(http.StatusForbidden, "パスワードが違います")
		} else {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("hash error: %v", err))
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("session error: %v", err))
	}
	sess.Values["userID"] = req.ID
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return c.JSON(http.StatusForbidden, "please login")
		}

		if sess.Values["userID"] == nil {
			return c.JSON(http.StatusForbidden, "please login")
		}

		userID, ok := sess.Values["userID"].(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, "Internal Server Error: something wrong when casting userID")
		}
		c.Set("userID", userID)

		return next(c)
	}
}

func GetWhoAmIHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	res := WhoAmIResponseBody{}

	err := dbs.Db.Get(&res, "SELECT id,name FROM users WHERE id=?", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

func GetLogoutHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.JSON(http.StatusForbidden, "please login")
	}

	sess.Values["userID"] = nil
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}
