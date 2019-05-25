package tweet

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/jmoiron/sqlx"
)

type TweetAdd struct{					//Tweetの構造体
	UserID				uuid.UUID
	Tweet				string
	CreatedTime			time.Time
}

var (
	db *sqlx.DB
)

func postTweetHandler(c echo.Context) error{
	tweet :=TweetAdd{}
	c.Bind(&tweet)

	db.Exec("INSERT INTO ? (TweetID, UserID,Tweet,CreatedTime,FavoNum) VALUES (?, ?,?,?,?)","Tweet",uuid.New() ,tweet.UserID,tweet.Tweet,tweet.CreatedTime,0)
	return c.NoContent(http.StatusOK)
}

