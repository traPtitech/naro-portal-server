package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Post struct {
	PostTime    string `json:"postTime,omitempty" db:"PostTime"`
	UserName    string `json:"userName,omitempty" db:"UserName"`
	UserID      string `json:"userID,omitempty" db:"UserID"`
	Body        string `json:"body,omitempty" db:"Body"`
	LikeCount   string `json:"likeCount,omitempty" db:"LikeCount"`
	RepostCount string `json:"repostCount,omitempty" db:"RepostCount"`
	ID          int    `json:"id,omitempty" db:"ID"`
}
type PostRequest struct {
	Body string `json:"body,omitempty" db:"Body"`
}

func AddPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		userName := sess.Values["UserName"]
		userID := sess.Values["UserID"]
		data := &PostRequest{}
		err = c.Bind(data)
		if err != nil { // エラーが発生した際
			// fmt.Sprintf("%+v", data): dataをstringに変換
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}
		State := `INSERT INTO posts(UserName, UserID, Body, LikeCount, RepostCount) VALUES (?,?,?,?,?)`
		_, err2 := db.Exec(State, userName, userID, data.Body, 0, 0)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err2.Error()))
		}

		return c.String(http.StatusOK, "success")
	}
}

func GetPostListHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		posts := []Post{}
		if err := db.Select(&posts, "SELECT * FROM posts ORDER BY ID DESC LIMIT 20"); errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("No Such Post List"))
		} else if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}

		return c.JSON(http.StatusOK, posts)
	}
}

func GetPostHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID := c.Param("postID")
		post := Post{}
		id, atoiErr := strconv.Atoi(postID)
		if atoiErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("IDをURLの末尾に入れてください"))
		}
		if err := db.Get(&post, "SELECT * FROM posts WHERE ID=?", id); errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("No Such City Name=%s", postID))
		} else if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf(err.Error()))
		}

		return c.JSON(http.StatusOK, post)
	}
}
