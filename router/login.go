package router

import (
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

func validation(user *UserInfo) string {
	// empty userID
	if user.UserID == "" {
		return "Empty userID"
	}
	// empty password
	if len(user.Password) < 6 {
		return "Weak password"
	}
	// reject multiple user
	count, err := model.Counter(user.UserID)
	if err != nil {
		return "Database is not working well, and cannot verify user infomation"
	}
	if count > 0 {
		return "The same userID already exists"
	}
	return ""
}

func PostRegisterHandler(c echo.Context) error {
	req := new(UserInfo)
	c.Bind(&req)

	// validation
	switch s := validation(req); s {
	case "Empty userID":
		return c.String(http.StatusBadRequest, s)
	case "Weak password":
		return c.String(http.StatusBadRequest, s)
	case "Database is not working well, and cannot verify user infomation":
		return c.String(http.StatusInternalServerError, s)
	case "The same userID already exists":
		return c.String(http.StatusConflict, s)
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
