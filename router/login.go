package router

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/sapphi-red/webengineer_naro-_7_server/database"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/auths"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	ID_MAX_LENGTH       = 30
	PASSWORD_MAX_LENGTH = 72
)

type LoginRequestBody struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req *LoginRequestBody) ValidateLoginInputs() error {
	if req.ID == "" {
		return errors.New("IDが空です")
	}
	if len(req.ID) > ID_MAX_LENGTH {
		return errors.New("IDが長すぎます")
	}

	if req.Password == "" {
		return errors.New("パスワードが空です")
	}
	if len(req.Password) > PASSWORD_MAX_LENGTH {
		return errors.New("パスワードが長すぎます")
	}
	return nil
}

func generatePass(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

func comparePass(hashedPass string, pass string) (isMismatch bool, err error) {
	isMismatch = false
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	if err != nil {
		isMismatch = (bcrypt.ErrMismatchedHashAndPassword == err)
		err = nil
	}
	return
}

func CreateLoginRoutes(e *echo.Echo) {
	e.POST("/signup", signUpHandler)
	e.POST("/login", loginHandler)
	e.POST("/logout", logoutHandler)
}

func signUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	err := req.ValidateLoginInputs()
	if err != nil {
		return return400(c, err)
	}

	hashedPass, err := generatePass(req.Password)
	if err != nil {
		return return500(c, "bcryptGenerateError", err)
	}

	exsists, err := database.Auths.GetUserExistance(req.ID)
	if err != nil {
		return return500(c, "UserGettingError", err)
	}
	if exsists {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	err = database.Auths.AddUser(req.ID, hashedPass)
	if err != nil {
		return return500(c, "UserAddingError", err)
	}

	return c.NoContent(http.StatusCreated)
}

func loginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := auths.AuthUser{}
	err := database.Auths.GetUser(req.ID, &user)
	if err != nil {
		return return500(c, "UserGettingError", err)
	}

	isMismatch, err := comparePass(user.HashedPass, req.Password)
	if isMismatch {
		return c.NoContent(http.StatusForbidden)
	}
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	err = database.Sessions.SetID(c, req.ID)
	if err != nil {
		return return500(c, "loginSessionDBError", err)
	}

	return c.NoContent(http.StatusOK)
}

func logoutHandler(c echo.Context) error {
	err := database.Sessions.Destroy(c)
	if err != nil {
		return return500(c, "logoutSessionDBError", err)
	}
	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := database.Sessions.Get(c)
		if err != nil {
			return return500(c, "checkSessionDBError", err)
		}

		if sess.Values["id"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("id", sess.Values["id"].(string))

		return next(c)
	}
}

type Me struct {
	ID string `json:"id,omitempty"`
}

func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Me{
		ID: c.Get("id").(string),
	})
}
