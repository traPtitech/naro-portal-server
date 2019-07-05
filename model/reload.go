package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

//GetIsReloadTimelineHandler Get /reloadTimeline/:userName Timelineの再読み込みするかの判定
func GetIsReloadTimelineHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	userName := c.Param("userName")
	lastReloadTime := sess.Values["LastReloadTime"].(time.Time)

	var newestMessage time.Time
	var userID string
	err = Db.Get(&userID, "SELECT ID FROM user WHERE name=?", userName)
	err = Db.Get(&newestMessage, "SELECT created_at FROM tweet WHERE user_ID=? ORDER BY created_at DESC LIMIT 1", userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if lastReloadTime.Before(newestMessage) {
		return c.String(http.StatusOK, "new message exist")
	}
	return c.NoContent(http.StatusOK)
}
