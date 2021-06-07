package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
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

type IsValidPasswordResponseBody bool

type IsValidNameResponseBody bool

type WhoAmIResponseBody struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

//投稿
type UpdatePostRequestBody struct {
	Text string `json:"text,omitempty" from:"text"`
}

//ファボ/ファボを外す
type FavPostRequestBody int

//投稿の取得
type PostResponseBody struct {
	PostData PostBody `json:"post_data"`
	FavUsers []string `json:"fav_users"`
}
type PostsResponseBody []PostResponseBody

type PostBody struct {
	ID       int    `json:"id,omitempty" db:"id"`
	UserID   string `json:"user_id,omitempty" db:"user_id"`
	Text     string `json:"text,omitempty" db:"text"`
	PostTime string `json:"post_time,omitempty" db:"post_time"`
}
type PostsBody []PostBody

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

	e.POST("/isvalidid", postIsValidIDHandler) //(Stringで)投げられたIDが使われているか(Boolで)返す
	e.POST("/isvalidname", postIsValidNameHandler)
	e.POST("/isvalidpassword", postIsValidPasswordHandler)
	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)

	withLogin.GET("/whoami", getWhoAmIHandler)
	withLogin.GET("/logout", getLogoutHandler)
	withLogin.POST("/updatepost", postUpdatePostHandler)
	withLogin.POST("/favpost", postFavPostHandler)
	withLogin.POST("/unfavpost", postUnFavPostHandler)
	withLogin.GET("/posts", getPostsHandler)

	e.Start(":13300")
}

func postIsValidIDHandler(c echo.Context) error {
	var req string
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err.Error()))
	}
	if len(req) == 0 || len(req) > 20 {
		return c.JSON(http.StatusOK, false)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	if count > 0 {
		return c.JSON(http.StatusOK, false)
	}
	return c.JSON(http.StatusOK, true)
}

func postIsValidNameHandler(c echo.Context) error {
	var req string
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err.Error()))
	}
	if len(req) == 0 || len(req) > 30 {
		return c.JSON(http.StatusOK, false)
	}
	return c.JSON(http.StatusOK, true)
}

func postIsValidPasswordHandler(c echo.Context) error {
	var req string
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err.Error()))
	}
	if len(req) == 0 || len(req) > 20 {
		return c.JSON(http.StatusOK, false)
	}
	return c.JSON(http.StatusOK, true)
}

func postSignUpHandler(c echo.Context) error {
	req := SignUpRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %s", err.Error()))
	} else if req.ID == "" || req.Name == "" || req.Password == "" {
		return c.String(http.StatusBadRequest, "Bad Request: 空の要素があります")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", req.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = db.Exec("INSERT INTO users (id, name, hashed_pass) VALUES (?, ?, ?)", req.ID, req.Name, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	user := User{}
	err = db.Get(&user, "SELECT * FROM users WHERE id=?", req.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.String(http.StatusForbidden, "パスワードが違います")
		} else {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("hash error: %v", err))
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("session error: %v", err))
	}
	sess.Values["userID"] = req.ID
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("session error: %v", err))
		}

		if sess.Values["userID"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userID", sess.Values["userID"].(string))

		return next(c)
	}
}

func getWhoAmIHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	res := WhoAmIResponseBody{}

	err := db.Get(&res, "SELECT id,name FROM users WHERE id=?", userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.JSON(http.StatusOK, res)
}

func postUpdatePostHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	req := UpdatePostRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	time := time.Now()
	_, err = db.Exec("INSERT INTO posts (user_id, text, post_time) VALUES (?, ?, ?)", userID, req.Text, time)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func postFavPostHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req int
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM favolates WHERE post_id=? AND user_id=?", req, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	if count > 0 {
		return c.NoContent(http.StatusOK)
	}

	_, err = db.Exec("INSERT INTO favolates (post_id, user_id) VALUES (?, ?)", req, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}
func postUnFavPostHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req int
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	_, err = db.Exec("DELETE favolates WHERE user_id=? AND post_id=?", userID, req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func favUsers(postID int) []string {
	var userIDs []string
	db.Select(&userIDs, "SELECT user_id FROM favolates WHERE post_id=?", postID)

	if len(userIDs) == 0 {
		return []string{}
	}

	return userIDs
}

func getPostsHandler(c echo.Context) error {
	reqIDStr := c.QueryParam("id")
	var posts PostsBody
	var err error

	if reqIDStr == "" {
		err = db.Select(&posts, "SELECT id, user_id, text, post_time FROM posts ORDER BY id DESC")
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
	} else {
		reqID, err := strconv.Atoi(reqIDStr)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
		}
		err = db.Select(&posts, "SELECT id, user_id, text, post_time FROM posts WHERE id=? ORDER BY id DESC", reqID)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
	}

	//favしたユーザーを取得、追加
	var res PostsResponseBody
	for i := 0; i < len(posts); i++ {
		post := posts[i]
		res = append(res, PostResponseBody{PostData: post, FavUsers: favUsers(post.ID)})
	}
	return c.JSON(http.StatusOK, res)
}

func getLogoutHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("session error: %v", err))
	}

	sess.Values["userID"] = nil
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}
