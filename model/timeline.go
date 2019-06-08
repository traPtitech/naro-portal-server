package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

//Tweet Tweetの構造体
type Tweet struct {
	TweetID   string    `json:"tweetID,omitempty"  db:"tweet_ID"`
	UserID    string    `json:"userID,omitempty"  db:"user_ID"`
	Tweet     string    `json:"tweet,omitempty"  db:"tweet"`
	CreatedAt time.Time `json:"createdAt,omitempty"  db:"created_at"`
	FavoNum   int       `json:"favoNum,omitempty"  db:"favo_num"`
}

//TweetIDOfPin Pin止めされたTweetの構造体
type TweetIDOfPin struct{
	TweetID string `json:"tweetID,omitempty" db:"tweet_ID"`
}

//GetTimelineHandler Get /timeline/:userName タイムライン
func GetTimelineHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	userName := c.Param("userName")

	tweets := []Tweet{}
	var userID string
	Db.Get(&userID, "SELECT ID FROM user WHERE name=?", userName)
	Db.Select(&tweets, "SELECT * FROM tweet WHERE user_ID=?", userID)
	sess.Values["LastReloadTime"]=time.Now()
	return c.JSON(http.StatusOK, tweets)
}

//GetPinHandler Get /pin/:userName タイムラインのピン
func GetPinHandler(c echo.Context) error {
	userName := c.Param("userName")

	dbPins := []TweetIDOfPin{}
	var userID string
	Db.Get(&userID, "SELECT ID FROM user WHERE name=?", userName)
	Db.Select(&dbPins, "SELECT tweet_ID FROM pin WHERE user_ID=?", userID)

	pins:=[]Tweet{}
	for i,v:=range dbPins{
		Db.Get(&pins[i],"SELECT * FROM tweet WHERE tweet_ID=?",v.TweetID)
	} 

	return c.JSON(http.StatusOK, pins)
}