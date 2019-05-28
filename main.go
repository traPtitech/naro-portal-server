package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
}

type Review struct {
	ID		 int	`json:"id" form:"id" db:"id"`
	Title	 string `json:"title,omitempty" form:"title" db:"title"`
	Contents string `json:"contents,omitempty" form:"Contents" db:"contents"`
	Username string `json:"username,omitempty" form:"username" db:"username"`
	FavCount int    `json:"fav_count" form:"fav_count" db:"fav_count"`
}

type Fav struct {
	FavID	 int	`json:"fav_id" form:"fav_id" db:"fav_id"`
	ReviewID string `json:"review_id,omitempty" form:"review_id" db:"review_id"`
	FavUser  string `json:"fav_user,omitempty" form:"fav_user" db:"fav_user"`
}

type Follow struct {
	FollowID      int	`json:"follow_id" form:"follow_id" db:"follow_id"`
	FollowUser    string `json:"follow_user,omitempty" form:"follow_user" db:"follow_user"`
	FollowedUser  string `json:"followed_user,omitempty" form:"followed_user" db:"followed_user"`
}

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
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
	e.POST("/api/login", postLoginHandler)
	e.POST("/api/signup", postSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/whoami", getWhoAmIHandler)
	withLogin.POST("/api/review", postReviewHandler)
	withLogin.GET("/api/show", getAllReviewHandler)
	withLogin.GET("/api/myreviews", getMyReviewHandler)
	withLogin.POST("/api/givefav", giveFavHandler)
	withLogin.POST("/api/follow", giveFollowHandler)
	withLogin.GET("/api/myfav", getFavInfoHandler)
	withLogin.GET("/api/mytimeline", getTimeLineHandler)
	withLogin.GET("/api/titles/:titleName", getTitleInfoHandler)
	withLogin.GET("/api/users/:userName", getUserInfoHandler)

	e.Start(":11240")
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

	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
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

type Me struct {
	Username string `json:"username,omitempty"  db:"username"`
}

func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Me{
		Username: c.Get("userName").(string),
	})
}

func postReviewHandler(c echo.Context) error {
	req := Review{}
	c.Bind(&req)

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]


	if req.Contents == "" {
		return c.String(http.StatusBadRequest, "レビューが空です")
	}

	db.Exec("INSERT INTO reviews (title, contents, username, fav_count) VALUES (?,?,?,?);", req.Title, req.Contents, username, 0)
	return c.NoContent(http.StatusCreated)
}

func getAllReviewHandler(c echo.Context) error {

	reviews := []Review{}
	db.Select(&reviews, "SELECT id,username,title,contents,fav_count FROM reviews;")
	fmt.Println(reviews)
	
	return c.JSON(http.StatusOK, reviews)
}

func getMyReviewHandler(c echo.Context) error {

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]
	fmt.Println(username)

	reviews := []Review{}
	db.Select(&reviews, "SELECT id,username,title,contents,fav_count FROM reviews WHERE username=?;", username)
	fmt.Println(reviews)
	
	return c.JSON(http.StatusOK, reviews)
}

func getTitleInfoHandler(c echo.Context) error {
	titleName := c.Param("titleName")
	fmt.Println(titleName)
	strings.Replace(titleName,"%20","Replaced",' ')
	reviews := []Review{}
	db.Select(&reviews, "SELECT id,username, contents,fav_count FROM reviews WHERE title=?;", titleName)
	fmt.Println(reviews)

	return c.JSON(http.StatusOK, reviews)
}

func getUserInfoHandler(c echo.Context) error {
	userName := c.Param("userName")
	fmt.Println(userName)
	strings.Replace(userName,"%20","Replaced",' ')
	reviews := []Review{}
	db.Select(&reviews, "SELECT id,title, contents,fav_count FROM reviews WHERE username=?;", userName)
	fmt.Println(reviews)

	return c.JSON(http.StatusOK, reviews)
}

func giveFavHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]//操作者のユーザー名の取得

	req := Review{}
	c.Bind(&req)
	fmt.Println(req)

	db.Exec("INSERT INTO Fav (review_id, fav_user) VALUES (?,?);",req.ID,username)

	return c.NoContent(http.StatusCreated)
}

func giveFollowHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]//操作者のユーザー名の取得

	req := Review{}
	c.Bind(&req)
	fmt.Println(req)

	db.Exec("INSERT INTO follow (follow_user, followed_user) VALUES (?,?);",username,req.Username)

	return c.NoContent(http.StatusCreated)
}

func getFavInfoHandler(c echo.Context) error {

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]//操作者のユーザー名の取得

	reviews := []Review{}
	db.Select(&reviews, "SELECT id, title, contents, username FROM reviews JOIN Fav ON id = review_id WHERE Fav.fav_user=?;", username)
	fmt.Println(reviews)

	return c.JSON(http.StatusOK, reviews)
}

func getTimeLineHandler(c echo.Context) error {

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	if sess.Values["userName"] == nil {
		return c.String(http.StatusForbidden, "please login")
	}
	c.Set("userName", sess.Values["userName"].(string))

	var username = sess.Values["userName"]//操作者のユーザー名の取得

	reviews := []Review{}
	db.Select(&reviews, "SELECT id, title, contents, username FROM reviews JOIN follow ON username = followed_user WHERE follow_user=?;", username)
	fmt.Println(reviews)

	return c.JSON(http.StatusOK, reviews)
}

