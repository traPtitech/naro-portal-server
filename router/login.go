package router

import (
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/purewhite404/naro-server/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserInfo struct {
	UserID   string `json:"user_id,omitempty"`
	Password string `json:"password,omitempty"`
}

func validation(user *UserInfo) (int, error) {
	// empty userID
	if user.UserID == "" {
		return http.StatusBadRequest, errors.New("Empty userID")
	}
	// empty password
	if len(user.Password) < 6 {
		return http.StatusBadRequest, errors.New("Weak password")
	}
	// reject multiple user
	count, err := model.Counter(user.UserID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Database is not working well, and cannot verify user infomation")
	}
	if count > 0 {
		return http.StatusConflict, errors.New("The same userID already exists")
	}
	return http.StatusAccepted, nil
}

func PostRegisterHandler(c echo.Context) error {
	req := new(UserInfo)
	c.Bind(&req)

	// validation
	if statusCode, err := validation(req); err != nil {
		return c.String(statusCode, fmt.Sprintf("%v", err))
	}

	// passwordをhash化してuserを作成
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Hash algorithm is not working well: %v", err))
	}
	err = model.InsertUserWithHashedPass(req.UserID, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Database is not working well: %v", err))
	}
	return c.String(http.StatusCreated, "User created successfully")
}

func PostLoginHandler(c echo.Context) error {
	req := new(UserInfo)
	c.Bind(&req)

	savedUser, err := model.SelectUser(req.UserID)
	if err != nil {
		// TODO: userIDがない場合とDBが壊れた場合でstatusを分ける
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Database is not working well: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedUser.HashedPass), []byte(req.Password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return c.String(http.StatusForbidden, "Wrong password")
		default:
			return c.String(http.StatusInternalServerError, "Database is not working well, and we cannot check your password")
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Something wrong in getting session: %v", err))
	}
	sess.Values["userID"] = req.UserID
	sess.Save(c.Request(), c.Response())

	return c.String(http.StatusOK, "Login successfully")
}

func GetMeHandler(c echo.Context) error {
	return c.String(http.StatusOK, c.Get("userID").(string))
}

func HasLoggedin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Something wrong in getting session: %v", err))
		}

		if sess.Values["userID"] == nil {
			return c.String(http.StatusForbidden, "Please login")
		}
		c.Set("userID", sess.Values["userID"].(string))

		return next(c)
	}
}
