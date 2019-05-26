package pin

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

var (
	db *sqlx.DB
)

//Pin Pinの構造体
type Pin struct {
	PinID   string `json:"pinID,omitempty"`
	UserID  string `json:"userID,omitempty"`
	TweetID string `json:"tweetID,omitempty"`
}

//PostPinHandler Post /pin ピン
func PostPinHandler(c echo.Context) error {
	pin := Pin{}
	c.Bind(&pin)

	db.Exec("INSERT INTO ? (pin_ID, user_ID,tweet_ID) VALUES (?, ?,?)", "Pin", uuid.New(), pin.UserID, pin.TweetID)
	return c.NoContent(http.StatusOK)
}
