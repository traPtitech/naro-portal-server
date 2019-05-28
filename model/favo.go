package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/labstack/echo-contrib/session"
)

//Favorite Favoriteの構造体
type Favorite struct {
	TweetID string `json:"tweetID,omitempty"  db:"tweet_ID"`
}

//PostAddFavoHandler Post /favoAdd Favo追加
func PostAddFavoHandler(c echo.Context) error {
	sess,err:=session.Get("sessions",c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	favo := Favorite{}
	c.Bind(&favo)

	var userID uuid.UUID
	Db.Get(&userID, "SELECT user_ID FROM Favorite WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], favo.TweetID)
	if userID != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var FavoNum int
	Db.Get(&FavoNum, "SELECT favo_num FROM Tweet WHERE tweet_ID=?", favo.TweetID)
	Db.Exec("UPDATE Tweet SET favo_num=? WHERE tweet_ID=?", FavoNum+1, favo.TweetID)

	Db.Exec("INSERT INTO (Favorite user_ID,tweet_ID,created_at) VALUES (?,?)", sess.Values["UserID"], favo.TweetID, time.Now())
	return c.NoContent(http.StatusOK)
}

//PostDeleteFavoHandler Post /favoDelete Favo消去
func PostDeleteFavoHandler(c echo.Context) error {
	sess,err:=session.Get("sessions",c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	favo := Favorite{}
	c.Bind(&favo)

	var FavoNum int
	Db.Get(&FavoNum, "SELECT favo_num FROM Tweet WHERE tweet_ID=?", favo.TweetID)
	Db.Exec("UPDATE Tweet SET favo_num=? WHERE tweet_ID=?", FavoNum-1, favo.TweetID)

	Db.Exec("DELETE FROM favorite WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], favo.TweetID)
	return c.NoContent(http.StatusOK)
}
