package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Me struct {
	Username string `json:"username,omitempty"  db:"username"`
}

type MyPost struct {
	ID      int64  `json:"id,omitempty" db:"ID"`
	Name    string `json:"name,omitempty"  db:"Name"`
	Content string `json:"content,omitempty"  db:"Content"`
	Time    int64  `json:"time,omitempty" db:"Time"`
	// TODO
	// Time int
	// Good int `json:"good", db:"GOOD"`
}

type MyPostID struct {
	ID int64 `json:"id,omitempty" db:"MAX(ID)"`
}

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

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/login", postLoginHandler)
	e.POST("/signup", postSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/whoami", getWhoAmIHandler)
	withLogin.GET("/twttt/get", func(c echo.Context) error {
		posts := []MyPost{}
		err := db.Select(&posts, "SELECT * FROM twttt")
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		return c.JSON(http.StatusOK, posts)
	})
	withLogin.POST("/twttt/post", postHandler)

	e.Start(":10600")
}

type MyContent struct {
	Content string `json:"content,omitempty"  db:"Content"`
}

func postHandler(c echo.Context) error {
	req := MyContent{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("server error 0: %v", err))
	}

	fmt.Printf("C: %#v", req)

	IDs := []MyPostID{}
	err = db.Select(&IDs, "SELECT MAX(ID) FROM twttt")
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error 1: %v", err))
	}

	sess, _ := session.Get("sessions", c)
	Name := sess.Values["userName"].(string)

	_, err = db.Exec("INSERT INTO twttt (ID, Name, Content, Time) VALUES (?, ?, ?, ?)", IDs[0].ID+1, Name, req.Content, time.Now().Unix())

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error 2: %v", err))
	}

	return c.NoContent(http.StatusCreated)
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
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)
	fmt.Printf("C: %#v", req)

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

func getWhoAmIHandler(c echo.Context) error {
	sess, _ := session.Get("sessions", c)

	return c.JSON(http.StatusOK, Me{
		Username: sess.Values["userName"].(string),
	})
}
