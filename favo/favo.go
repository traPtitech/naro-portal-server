package favo

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/jmoiron/sqlx"
)

type Favorite struct{				//Favoriteの構造体
	UserID				uuid.UUID
	MessageID			uuid.UUID
}

var (
	db *sqlx.DB
)

func postFavoHandler(c echo.Context) error{
	favo:=Favorite{}
	c.Bind(&favo)

	var userID uuid.UUID
	db.Get(&userID,"SELECT user_ID FROM favorite WHERE user_ID=? AND message_ID=?",uuid.String(favo.UserID),uuid.String(favo.MessageID))
	if userID!=nil{
		return c.NoContent(http.StatusBadRequest)
	}

	db.Exec("INSERT INTO favorite (user_ID,message_ID) VALUES (?,?)",uuid.String(favo.UserID),uuid.String(favo.MessageID))
	return c.NoContent(http.StatusOK)
}

func deleteFavoHandler(c echo.Context) error{
	favo:=Favorite{}
	c.Bind(&favo)

	db.Exec("DELETE FROM favorite WHERE user_ID=? AND message_ID=?",uuid.String(favo.UserID),uuid.String(favo.MessageID))
	return c.NoContent(http.StatusOK)
}