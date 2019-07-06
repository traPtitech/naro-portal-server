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
	Db.Get(&userID, "SELECT user_ID FROM pin WHERE tweet_ID=?", tweetID)
	if userID != "" {
		_, err := Db.Exec("DELETE FROM pin WHERE user_ID=? AND tweet_ID=?", c.Get("UserID"), pin.TweetID)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
	_, err := Db.Exec("INSERT INTO pin (pin_ID, user_ID,tweet_ID) VALUES (?, ?,?)", uuid.New(), c.Get("UserID"), pin.TweetID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
