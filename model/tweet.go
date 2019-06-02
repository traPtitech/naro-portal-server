package model

import (
	"fmt"

	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
)

//AddTweet Tweetの構造体
type AddTweet struct {
	Tweet string `json:"tweet,omitempty"`
}

//PostTweetHandler Post /tweet Tweet追加
func PostTweetHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	tweet := AddTweet{}
	c.Bind(&tweet)

	Db.Exec("INSERT INTO tweet (tweet_ID,user_ID,tweet,created_at,favo_num) VALUES (?,?,?,?,?)", uuid.New(), sess.Values["UserID"], tweet.Tweet, time.Now(), 1)
	return c.NoContent(http.StatusOK)
}
