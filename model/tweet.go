package model

import (
	"fmt"

	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
)

//Add Tweetの構造体
type Add struct {
	Tweet  string `json:"tweet,omitempty"`
}

//PostTweetHandler Post /tweet Tweet追加
func PostTweetHandler(c echo.Context) error {
	sess,err:=session.Get("sessions",c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	tweet := Add{}
	c.Bind(&tweet)

	Db.Exec("INSERT INTO tweet (tweet_ID,user_ID,tweet,created_at,favo_num) VALUES (?,?,?,?,?)", uuid.New(), sess.Values["UserID"], tweet.Tweet, time.Now(), 0)
	return c.NoContent(http.StatusOK)
}
