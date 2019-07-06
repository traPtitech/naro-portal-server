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

//PostFavoHandler Post /favo Favo追加,消去
func PostFavoHandler(c echo.Context) error {
	favo := Favorite{}
	c.Bind(&favo)

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM favorite WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), favo.TweetID)
	if userID != "" {
		_, err := Db.Exec("DELETE FROM favorite WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), favo.TweetID)
		_, err = Db.Exec("UPDATE tweet SET favo_num=favo_num-1 WHERE tweet_ID=?", favo.TweetID)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
	_, err := Db.Exec("INSERT INTO favorite (favo_ID,user_ID,tweet_ID,created_at) VALUES (?,?,?,?)", uuid.New(), c.Get("UserID"), favo.TweetID, time.Now())
	_, err = Db.Exec("UPDATE tweet SET favo_num=favo_num+1 WHERE tweet_ID=?", favo.TweetID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
