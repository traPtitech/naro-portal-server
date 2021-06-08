package main

import (
	"fmt"
	"errors"
	"log"
	"net/http"
	"os"
	//"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect(
		"mysql", 
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOSTNAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/home", getTweetHandler)
	withLogin.POST("/home", postTweetHandler)
	withLogin.GET("/:userid", getAccountHome)

	e.Start(":11800")
}

type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("binding error: %v", err))
	}

	if req.Password == "" || req.Username == "" {
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

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
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("binding error: %v", err))
	}

	user := User{}
	err = db.Get(&user, "SELECT * FROM users WHERE username = ?", req.Username)
	// users { Username, HashedPass }

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
	// sessions { id, session_data, etc... }

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
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

type Tweet struct {
	DateTime time.Time `json:"datetime,omitempty" db:"DateTime"`
	UserID   string `json:"userid,omitempty" db:"UserID"`
	Text     string `json:"text,omitempty" db:"Tweet"`
}

type TweetRequestBody struct {
	Text string `json:"text,omitempty" form:"text"`
}

func getTweetHandler(c echo.Context) error {
	tweets := []Tweet{}
	err := db.Select(&tweets, "SELECT * FROM tweets ORDER BY DateTime DESC LIMIT 20")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.JSON(http.StatusOK, tweets)
}

func postTweetHandler(c echo.Context) error {
	req := TweetRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("binding error: %v", err))
	}

	if req.Text == "" {
		return c.String(http.StatusBadRequest, "empty tweet");
	}

	username := c.Get("userName");
	_, err = db.Exec("INSERT INTO tweets (DateTime, UserID, Tweet) VALUES (?, ?, ?)", time.Now(), username, req.Text)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func getAccountHome(c echo.Context) error {
	UserID := c.Param("userid")
	tweets := []Tweet{}
	err := db.Select(&tweets, "SELECT * FROM tweets WHERE UserID = ? ORDER BY DateTime DESC LIMIT 20", UserID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.JSON(http.StatusOK, tweets)
}