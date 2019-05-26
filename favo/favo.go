package favo

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

//Favorite Favoriteの構造体
type Favorite struct {
	UserID  string `json:"userID,omitempty"  db:"user_ID"`
	TweetID string `json:"tweetID,omitempty"  db:"tweet_ID"`
}

var (
	db *sqlx.DB
)

//PostAddFavoHandler Post /FavoAdd Favo追加
func PostAddFavoHandler(c echo.Context) error {
	favo := Favorite{}
	c.Bind(&favo)

	var userID uuid.UUID
	db.Get(&userID, "SELECT user_ID FROM Favorite WHERE user_ID=? AND tweet_ID=?", favo.UserID, favo.TweetID)
	if userID != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var FavoNum int
	db.Get(&FavoNum, "SELECT favo_num FROM Tweet WHERE tweet_ID=?", favo.TweetID)
	db.Exec("UPDATE Tweet SET favo_num=? WHERE tweet_ID=?", FavoNum+1, favo.TweetID)

	db.Exec("INSERT INTO Favorite (user_ID,tweet_ID,created_at) VALUES (?,?)", favo.UserID, favo.TweetID, time.Now())
	return c.NoContent(http.StatusOK)
}

//PostDeleteFavoHandler Post /Favo_Delete Favo消去
func PostDeleteFavoHandler(c echo.Context) error {
	favo := Favorite{}
	c.Bind(&favo)

	var FavoNum int
	db.Get(&FavoNum, "SELECT favo_num FROM Tweet WHERE tweet_ID=?", favo.TweetID)
	db.Exec("UPDATE Tweet SET favo_num=? WHERE tweet_ID=?", FavoNum-1, favo.TweetID)

	db.Exec("DELETE FROM favorite WHERE user_ID=? AND tweet_ID=?", favo.UserID, favo.TweetID)
	return c.NoContent(http.StatusOK)
}
