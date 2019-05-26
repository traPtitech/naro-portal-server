package tweet

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

//Add Tweetの構造体
type Add struct {
	UserID string `json:"ID,omitempty"`
	Tweet  string `json:"tweet,omitempty"`
}

var (
	db *sqlx.DB
)

//PostTweetHandler Post /tweet Tweet追加
func PostTweetHandler(c echo.Context) error {
	tweet := Add{}
	c.Bind(&tweet)

	db.Exec("INSERT INTO ? (tweet_ID, user_ID,tweet,created_at,favo_num) VALUES (?, ?,?,?,?)", "Tweet", uuid.New(), tweet.UserID, tweet.Tweet, time.Now(), 0)
	return c.NoContent(http.StatusOK)
}
