package router

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/purewhite404/naro-server/model"
	"net/http"
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
	err := c.Bind(tweet)
	if err != nil {
		return c.String(http.StatusBadRequest, "Not suitable for JsonTweet format")
	}

	u, err := uuid.NewRandom()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Cannot create tweet uuid")
	}
	tweet.ID = u.String()

	err = model.InsertTweet(tweet)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot insert tweet into database: %v", err))
	}
	return c.String(http.StatusCreated, "Post tweet successfully")
}
