package model

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

//ChangePin Pinの構造体
type ChangePin struct {
	TweetID string `json:"tweetID,omitempty"`
}

//PostPinHandler Post /pin ピン
func PostPinHandler(c echo.Context) error {
	pin := ChangePin{}
	c.Bind(&pin)

	var userID string
	Db.Get(&userID, "SELECT user_ID FROM tweet WHERE tweet_ID=?", pin.TweetID)
	if userID != c.Get("UserID") {
		return c.String(http.StatusInternalServerError, "あなたのTweetではありません")
	}

	Db.Exec("INSERT INTO pin (pin_ID, user_ID,tweet_ID) VALUES (?, ?,?)", uuid.New(), c.Get("UserID"), pin.TweetID)
	return c.NoContent(http.StatusOK)
}

//DeletePinHandler Delete /pin Pin消去
func DeletePinHandler(c echo.Context) error {
	pin := ChangePin{}
	c.Bind(&pin)

	err := Db.Exec("DELETE FROM pin WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), pin.TweetID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

//GetIsPinHandler Get /isPin/:tweetID ピンを入れたかの確認
func GetIsPinHandler(c echo.Context) error {
	tweetID := c.Param("tweetID")

	var userID string
	err := Db.Get(&userID, "SELECT user_ID FROM pin WHERE tweet_ID=?", tweetID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if userID != "" {
		return c.NoContent(http.StatusOK)
	}

	return c.String(http.StatusOK, "none")
}
