package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type tweet_content struct {
	Tweet_id int    `json:"tweet_id,omitempty" db:"tweet_id"`
	Text     string `json:"text,omitempty" db:"text"`
	Fav      bool   `json:"fav,omitempty" db:"fav"`
}

var (
	db *sqlx.DB
)

func main() {
	e := echo.New()

	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	} //DataBaseに接続してる（エラーならその内容を表示）
	db = _db //func mainの外でも使えるように外で定義したdbに_dbを代入

	e.POST("/PostTweet", tweet_handler)
}

func tweet_handler(c echo.Context) error {
	tweet_text_handler(c)

	return c.NoContent(http.StatusOK)
}
func tweet_text_handler(c echo.Context) error {
	tweet_content := tweet_content{}
	c.Bind(&tweet_content) //対応するリクエストのKeyの値を構造体にうまくあてはめてくれる

	if tweet_content.Text == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "Void content")
	}
	return c.NoContent(http.StatusOK)
}
