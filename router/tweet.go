package router

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/purewhite404/naro-server/model"
	"net/http"
	"time"
)

func GetTweetHandler(c echo.Context) error {
	tweets, err := model.SelectTweet()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot get timeline: %v", err))
	}

	return c.JSON(http.StatusOK, tweets)
}

func PostTweetHandler(c echo.Context) error {
	tweet := new(model.JsonTweet)
	err := c.Bind(tweet) // ここでidとcreated_atはダミー値なので以下の処理においてサーバ側で生成する
	if err != nil {
		return c.String(http.StatusBadRequest, "Not suitable for JsonTweet format")
	}

	// uuidを生成し、tweetのidとする
	u, err := uuid.NewRandom()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot create tweet uuid")
	}
	tweet.ID = u.String()

	// created_atを生成しタイムスタンプとする
	tweet.CreatedAt = time.Now()

	err = model.InsertTweet(tweet)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot insert tweet into database: %v", err))
	}
	return c.JSON(http.StatusCreated, tweet)
}
