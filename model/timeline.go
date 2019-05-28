package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

//Tweet Tweetの構造体
type Tweet struct {
	TweetID   string    `json:"tweetID,omitempty"  db:"tweet_ID"`
	UserID    string    `json:"userID,omitempty"  db:"user_ID"`
	Tweet     string    `json:"tweet,omitempty"  db:"tweet"`
	CreatedAt time.Time `json:"createdAt,omitempty"  db:"created_at"`
	FavoNum   int       `json:"favoNum,omitempty"  db:"favo_num"`
}

//GetTimeLineHandler Get /timeline/:userName タイムライン
func GetTimeLineHandler(c echo.Context) error {
	userName := c.Param("userName")

	tweets := []Tweet{}
	var userID string
	Db.Get(&userID, "SELECT ID FROM user WHERE name=?", userName)
	Db.Select(&tweets, "SELECT * FROM tweet WHERE user_ID=?", userID)
	return c.JSON(http.StatusOK, tweets)
}
