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

//Pin Pinの構造体
type Pin struct {
	PinID     string    `json:"pinID,omitempty" db:"pin_ID"`
	TweetID   string    `json:"tweetID,omitempty"  db:"tweet_ID"`
	UserID    string    `json:"userID,omitempty"  db:"user_ID"`
	Tweet     string    `json:"tweet,omitempty"  db:"tweet"`
	CreatedAt time.Time `json:"createdAt,omitempty"  db:"created_at"`
	FavoNum   int       `json:"favoNum,omitempty"  db:"favo_num"`
}

//Favo Favoの構造体
type Favo struct {
	FavoID    string    `json:"favoID,omitempty" db:"favo_ID"`
	TweetID   string    `json:"tweetID,omitempty"  db:"tweet_ID"`
	UserID    string    `json:"userID,omitempty"  db:"user_ID"`
	UserName  string    `json:"userName,omitempty"`
	Tweet     string    `json:"tweet,omitempty"  db:"tweet"`
	CreatedAt time.Time `json:"createdAt,omitempty"  db:"created_at"`
	FavoNum   int       `json:"favoNum,omitempty"  db:"favo_num"`
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
	err = Db.Select(&tweets, "SELECT * FROM tweet JOIN user ON tweet.user_ID = user.ID WHERE user.name = ? ORDER BY created_at DESC", userName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Values["LastReloadTime"] = time.Now()

	return c.JSON(http.StatusOK, tweets)
}

//GetPinHandler Get /timelinePin/:userName タイムラインのピン
func GetPinHandler(c echo.Context) error {
	userName := c.Param("userName")
	pins := []Pin{}
	err := Db.Select(&pins, "SELECT pin.pin_ID,tweet.tweet_ID,tweet.user_ID,tweet.tweet,tweet.created_at,tweet.favo_num FROM pin JOIN tweet ON pin.tweet_ID = tweet.tweet_ID JOIN user ON pin.user_ID = user.ID WHERE user.name = ? ORDER BY tweet.created_at", userName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pins)
}

//GetFavoHandler Get /timelineFavo/:userName タイムラインのピン
func GetFavoHandler(c echo.Context) error {
	userName := c.Param("userName")
	favos := []Favo{}
	err := Db.Select(&favos, "SELECT favorite.favo_ID,tweet.tweet_ID,tweet.user_ID,tweet.tweet,tweet.created_at,tweet.favo_num FROM favorite JOIN tweet ON favorite.tweet_ID = tweet.tweet_ID JOIN user ON favorite.user_ID = user.ID WHERE user.name = ? WHERE favorite.user_ID = ? ORDER BY tweet.created_at", userName)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, favos)
}
