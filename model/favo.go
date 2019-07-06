package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

//Favorite Favoriteの構造体
type Favorite struct {
	TweetID string `json:"tweetID,omitempty"  db:"tweet_ID"`
}

//PostFavoHandler Post /favo Favo追加
func PostFavoHandler(c echo.Context) error {
	favo := Favorite{}
	c.Bind(&favo)

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM favorite WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), favo.TweetID)
	if userID != "" {
		return c.NoContent(http.StatusBadRequest)
	}

	_, err = Db.Exec("UPDATE tweet SET favo_num=favo_num+1 WHERE tweet_ID=?", favo.TweetID)

	_, err = Db.Exec("INSERT INTO favorite (favo_ID,user_ID,tweet_ID,created_at) VALUES (?,?,?,?)", uuid.New(), c.Get("UserID"), favo.TweetID, time.Now())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

//DeleteFavoHandler Delete /favo Favo消去
func DeleteFavoHandler(c echo.Context) error {
	favo := Favorite{}
	c.Bind(&favo)

	_, err := Db.Exec("UPDATE tweet SET favo_num=favo_num-1 WHERE tweet_ID=?", favo.TweetID)
	_, err = Db.Exec("DELETE FROM favorite WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), favo.TweetID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

//GetIsFavoHandler Get /isFavo/:tweetID ファボを入れたかの確認
func GetIsFavoHandler(c echo.Context) error {
	tweetID := c.Param("tweetID")

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM favorite WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), tweetID)
	if userID != "" {
		return c.NoContent(http.StatusOK)
	}

	return c.String(http.StatusOK, "none")
}
