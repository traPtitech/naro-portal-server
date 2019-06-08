package model

import (
	"fmt"
	"time"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

//GetIsReloadTimelineHandler Get /reloadTimeline/:userName Timelineの再読み込みするかの判定
func GetIsReloadTimelineHandler(c echo.Context) error{
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	userName:=c.Param("userName")
	lastReloadTime:=sess.Values["LastReloadTime"].(time.Time)

	var newestMessage time.Time
	var userID string
	Db.Get(&userID,"SELECT ID FROM user WHERE user_name=?",userName)
	Db.Get(&newestMessage,"SELECT created_at FROM tweet ORDER BY created_at DESC LIMIT 1 WHERE user_ID=?",userID)
	sess.Values["LastReloadTime"]=time.Now()
	
	if lastReloadTime.Before(newestMessage){
		return c.String(http.StatusOK,"new message exist")
	}
	return c.NoContent(http.StatusOK)
}