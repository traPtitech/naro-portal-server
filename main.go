package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	e.POST("/login", postLoginHandler)
	e.POST("/logout", postLogoutHandler)
	e.POST("/signup", postSignUpHandler)
	e.POST("/tweet", postTweetHandler)
	e.POST("/deleteTweet", postDeleteTweetHandler)
	e.POST("/favorite", postFavoriteHandler)
	e.POST("/unfavorite", postUnfavoriteHandler)
	e.POST("/updateProfile", postUpdateProfileHandler)
	e.POST("/follow", postFollowHandler)
	e.POST("/unfollow", postUnfollowHandler)
	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/list", getTweetListHandler)
	withLogin.GET("/favoriteList/:userId", getFavoriteListHandler)
	withLogin.GET("/followList/:userId", getFollowListHandler)
	withLogin.GET("/followerList/:userId", getFollowerListHandler)
	withLogin.GET("/whoami",getWhoAmIHandler)
	withLogin.GET("/userInfo/:userId", getUserInfoHandler)
	e.Start(":10400")
}

type LoginRequestBody struct {
	Userid string `json:"userid,omitempty" form:"userid"`
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type TweetRequestBody struct {
	Id string `json:"id,omitempty" form:"id"`
	Userid string `json:"userid,omitempty" form:"userid"`
	Time string `json:"time,omitempty" form:"time"`
	Text string `json:"text,omitempty" form:"text"`
}

type User struct {
	Userid   string `json:"userid,omitempty"  db:"Userid"`
	Username   string `json:"username,omitempty"  db:"Username"`
	HashedPass string `json:"-"  db:"HashedPass"`
	Biography string `json:"biography,omitempty"  db:"Biography"`
	Website string `json:"website,omitempty"  db:"Website"`
}

type Tweet struct {
	Id   string `json:"id,omitempty"  db:"Id"`
	Userid   string `json:"userid,omitempty"  db:"Userid"`
	Username   string `json:"username,omitempty"  db:"Username"`
	Time   string `json:"time,omitempty"  db:"Time"`
	Text   string `json:"text,omitempty"  db:"Text"`
}

type Favorite struct {
	Userid   string `json:"userid,omitempty"  db:"Userid"`
	Tweetid   string `json:"tweetid,omitempty"  db:"Tweetid"`
}
type Follow struct {
	Userid   string `json:"userid,omitempty"  db:"Userid"`
	Followid   string `json:"followid,omitempty"  db:"Followid"`
}
func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	// もう少し真面目にバリデーションするべき
	if req.Password == "" || req.Userid == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	// ユーザーの存在チェック
	var count int

	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE Userid=?", req.Userid)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = db.Exec("INSERT INTO users (Userid,Username,HashedPass,Biography,Website) VALUES (?, ?, ?,'','')", req.Userid, req.Username, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE Userid=?", req.Userid)
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
	sess.Values["userId"] = req.Userid
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}
func postLogoutHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userId"] = ""
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
func postTweetHandler(c echo.Context) error {
	req := TweetRequestBody{}
	c.Bind(&req)
	_, err := db.Exec("INSERT INTO tweets (Id,Userid,Time,Text) VALUES (?, ?, ?, ?)", req.Id, req.Userid, req.Time, req.Text)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)

}
func postDeleteTweetHandler(c echo.Context) error {
	req := TweetRequestBody{}
	c.Bind(&req)
	_, err := db.Exec("DELETE FROM tweets WHERE Id=?", req.Id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)

}
func postFavoriteHandler(c echo.Context) error {
	req := Favorite{}
	c.Bind(&req)
	_, err := db.Exec("INSERT INTO favorites (Userid,Tweetid) VALUES (?, ?)", req.Userid, req.Tweetid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
func postUnfavoriteHandler(c echo.Context) error {
	req := Favorite{}
	c.Bind(&req)
	_, err := db.Exec("DELETE FROM favorites WHERE Userid=? AND Tweetid=?", req.Userid, req.Tweetid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
func postFollowHandler(c echo.Context) error {
	req := Follow{}
	c.Bind(&req)
	_, err := db.Exec("INSERT INTO follows (Userid,Followid) VALUES (?, ?)", req.Userid, req.Followid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
func postUnfollowHandler(c echo.Context) error {
	req := Follow{}
	c.Bind(&req)
	_, err := db.Exec("DELETE FROM follows WHERE Userid=? AND Followid=?", req.Userid, req.Followid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
func postUpdateProfileHandler(c echo.Context) error {
	req := User{}
	c.Bind(&req)
	_, err := db.Exec("UPDATE users set Username=?,Biography=?,Website=? WHERE Userid=?", req.Username, req.Biography, req.Website, req.Userid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)

}
func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userId"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userId", sess.Values["userId"].(string))
		//c.String(http.StatusOK, sess.Values["userName"].(string)+"\n")
		return next(c)
	}
}
func getTweetListHandler(c echo.Context) error {
	rows, _ := db.Query("SELECT Id,T.Userid,Username,Time,Text FROM tweets AS T INNER JOIN users AS U ON T.Userid = U.Userid ORDER BY Time DESC")
	defer rows.Close()
	var result []Tweet
	for rows.Next() {
		tweet := Tweet{}
		rows.Scan(&tweet.Id,&tweet.Userid,&tweet.Username,&tweet.Time,&tweet.Text)
		result = append(result, tweet)
	}
	return c.JSON(http.StatusOK, result)
}
func getFavoriteListHandler(c echo.Context) error {
	userId := c.Param("userId")
	rows, _ := db.Query("SELECT Id,T.Userid,Username,Time,Text FROM tweets AS T INNER JOIN users AS U ON T.Userid = U.Userid INNER JOIN favorites AS F ON T.Id=F.Tweetid WHERE F.Userid=? ORDER BY Time DESC",userId)
	defer rows.Close()
	var result []Tweet
	for rows.Next() {
		tweet := Tweet{}
		rows.Scan(&tweet.Id,&tweet.Userid,&tweet.Username,&tweet.Time,&tweet.Text)
		result = append(result, tweet)
	}
	return c.JSON(http.StatusOK, result)
}
func getFollowListHandler(c echo.Context) error {
	userId := c.Param("userId")
	rows, _ := db.Query("SELECT U.Userid,U.Username,U.Hashedpass,U.Biography,U.Website FROM follows AS F INNER JOIN users AS U ON F.Followid=U.Userid WHERE F.Userid=?",userId)
	defer rows.Close()
	var result []User
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Userid,&user.Username,&user.HashedPass,&user.Biography,&user.Website)
		result = append(result, user)
	}
	return c.JSON(http.StatusOK, result)
}
func getFollowerListHandler(c echo.Context) error {
	userId := c.Param("userId")
	rows, _ := db.Query("SELECT U.Userid,U.Username,U.Hashedpass,U.Biography,U.Website FROM follows AS F INNER JOIN users AS U ON F.Userid=U.Userid WHERE F.Followid=?",userId)
	defer rows.Close()
	var result []User
	for rows.Next() {
		user := User{}
		rows.Scan(&user.Userid,&user.Username,&user.HashedPass,&user.Biography,&user.Website)
		result = append(result, user)
	}
	return c.JSON(http.StatusOK, result)
}
func getWhoAmIHandler(c echo.Context) error {
	if c.Get("userId") == "" {
        return c.String(http.StatusForbidden, "please login")
    }
	user := User{}
	db.Get(&user, "SELECT * FROM users WHERE Userid=?", c.Get("userId").(string))
	return c.JSON(http.StatusOK, user)
}
func getUserInfoHandler(c echo.Context) error {
	userId := c.Param("userId")
	user := User{}
	db.Get(&user, "SELECT * FROM users WHERE Userid=?", userId)
	if user.Userid == "" {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}