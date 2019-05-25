package timeline

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/jmoiron/sqlx"
)

type Tweet struct{					//Tweetの構造体
	TweetID				uuid.UUID
	UserID				uuid.UUID
	Tweet				string
	CreatedTime			time.Time
	FavoNum				int
}

var (
	db *sqlx.DB
)

func getTimeLineHandler(c echo.Context) error{
	userName:=c.Param("userName")

	tweets:=[]Tweet{}
	var userID uuid.UUID
	db.Get(&userID,"SELECT user_ID FROM User WHERE user_name=?",userName)
	db.Select(&tweets,"SELECT * FROM Tweet WHERE user_ID=?",userID)
	return c.JSON(http.StatusOK,tweets)
}