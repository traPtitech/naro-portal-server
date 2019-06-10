package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
)

//Favorite Favoriteの構造体
type Favorite struct {
	TweetID string `json:"tweetID,omitempty"  db:"tweet_ID"`
}

//PostFavoHandler Post /favo Favo追加
func PostFavoHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	favo := Favorite{}
	c.Bind(&favo)

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM favorite WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], favo.TweetID)
	if userID != "" {
		return c.NoContent(http.StatusBadRequest)
	}

	var FavoNum int
	Db.Get(&FavoNum, "SELECT favo_num FROM tweet WHERE tweet_ID=?", favo.TweetID)
	Db.Exec("UPDATE tweet SET favo_num=? WHERE tweet_ID=?", FavoNum+1, favo.TweetID)

	Db.Exec("INSERT INTO favorite (favo_ID,user_ID,tweet_ID,created_at) VALUES (?,?,?,?)", uuid.New(), sess.Values["UserID"], favo.TweetID, time.Now())
	return c.NoContent(http.StatusOK)
}

//DeleteFavoHandler Delete /favo Favo消去
func DeleteFavoHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	favo := Favorite{}
	c.Bind(&favo)

	var FavoNum int
	Db.Get(&FavoNum, "SELECT favo_num FROM tweet WHERE tweet_ID=?", favo.TweetID)
	Db.Exec("UPDATE tweet SET favo_num=? WHERE tweet_ID=?", FavoNum-1, favo.TweetID)

	Db.Exec("DELETE FROM favorite WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], favo.TweetID)
	return c.NoContent(http.StatusOK)
}

//GetIsFavoHandler Get /isFavo/:tweetID ファボを入れたかの確認
func GetIsFavoHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	tweetID := c.Param("tweetID")

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM favorite WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], tweetID)
	if userID != "" {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, "none")
}
