package pin

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type Pin struct{					//Pinの構造体
	PinID				uuid.UUID
	UserID				uuid.UUID
	MessageID			uuid.UUID
}

func postPinHandler(c echo.Context) error{
	pin:=Pin{}
	c.Bind(&pin)

	db.Exec("INSERT INTO ? (PinID, UserID,MessageID) VALUES (?, ?,?)","Pin",uuid.New() ,pin.UserID,pin.MessageID)
	return c.NoContent(http.StatusOK)
}