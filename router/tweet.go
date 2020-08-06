package router

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/purewhite404/naro-server/model"
	"net/http"
)

func GetTweetHandler(c echo.Context) error {
	tweets, err := model.GetTweets()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot get timeline: %v", err))
	}

	return c.JSON(http.StatusOK, tweets)
}

func PostTweetHandler(c echo.Context) error {
	req := new(model.RequestTweet)
	err := c.Bind(req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Not suitable for JsonTweet format: %v", err))
	}

	// uuidを生成し、tweetのidとする
	u, err := uuid.NewRandom()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot create tweet uuid: %v", err))
	}
	uuid := u.String()

	// uuidとreq
	err = model.InsertTweet(uuid, req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot insert tweet into database: %v", err))
	}

	// 今INSERTしたものをもう一度取り出す
	tweetFromDB, err := model.GetPostedTweet(uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Tweet was created, but cannot get from DB: %v", err))
	}
	return c.JSON(http.StatusCreated, tweetFromDB)
}
