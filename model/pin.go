package model

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
)

//Pin Pinの構造体
type Pin struct {
	UserID  string `json:"userID,omitempty"`
	TweetID string `json:"tweetID,omitempty"`
}

//PostPinHandler Post /pin ピン
func PostPinHandler(c echo.Context) error {
	pin := Pin{}
	c.Bind(&pin)

	sess,err:=session.Get("sessions",c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	
	Db.Exec("INSERT INTO pin (pin_ID, user_ID,tweet_ID) VALUES (?, ?,?)", uuid.New(),sess.Values["UserID"], pin.UserID, pin.TweetID)
	return c.NoContent(http.StatusOK)
}
