package messages

import (
	"fmt"
	"kuragate-server/dbs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
)

//投稿
type PostMassageRequestBody struct {
	Text string `json:"text,omitempty" from:"text"`
}

type UpdatePostRequestBody struct {
	Text string `json:"text,omitempty" from:"text"`
}

//ファボ/ファボを外す
type FavPostRequestBody int

//投稿の取得
//一つの投稿
type GetMessageBody struct {
	ID       int      `json:"id,omitempty" db:"id"`
	UserID   string   `json:"user_id,omitempty" db:"user_id"`
	Text     string   `json:"text,omitempty" db:"text"`
	PostTime string   `json:"post_time,omitempty" db:"post_time"`
	FavUsers []string `json:"fav_users"`
}
type GetMessagesBody []GetMessageBody

func PostUpdatePostHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	req := UpdatePostRequestBody{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	time := time.Now()
	_, err = dbs.Db.Exec("INSERT INTO messages (user_id, text, post_time) VALUES (?, ?, ?)", userID, req.Text, time)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func PostMessageHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	req := PostMassageRequestBody{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	time := time.Now()
	_, err = dbs.Db.Exec("INSERT INTO messages (user_id, text, post_time) VALUES (?, ?, ?)", userID, req.Text, time)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)

}

func PutMessageFavHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req int
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	var count int
	err = dbs.Db.Get(&count, "SELECT COUNT(*) FROM favolates WHERE message_id=? AND user_id=?", req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	if count > 0 {
		return c.NoContent(http.StatusOK)
	}

	_, err = dbs.Db.Exec("INSERT INTO favolates (message_id, user_id) VALUES (?, ?)", req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func DeleteMessageFavHandler(c echo.Context) error {
	userID := c.Get("userID").(string)
	var req int
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
	}

	_, err = dbs.Db.Exec("DELETE favolates WHERE user_id=? AND message_id=?", userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func favUsers(messageID int) ([]string, error) {
	var userIDs []string
	err := dbs.Db.Select(&userIDs, "SELECT user_id FROM favolates WHERE message_id=?", messageID)

	if len(userIDs) == 0 {
		return []string{}, err
	}

	return userIDs, nil
}

func GetMassagesHandler(c echo.Context) error {
	var messages GetMessagesBody

	err := dbs.Db.Select(&messages, "SELECT id, user_id, text, post_time FROM messages ORDER BY id DESC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	//favしたユーザーを取得
	for i := 0; i < len(messages); i++ {
		message := messages[i]
		users, err := favUsers(message.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		message.FavUsers = users
	}
	return c.JSON(http.StatusOK, messages)
}

func GetSingleMassageHandler(c echo.Context) error {
	var message GetMessageBody
	id := c.QueryParam("id")

	err := dbs.Db.Get(&message, "SELECT id, user_id, text, post_time FROM messages WHERE id=?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	users, err := favUsers(message.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	message.FavUsers = users

	return c.JSON(http.StatusOK, message)
}
