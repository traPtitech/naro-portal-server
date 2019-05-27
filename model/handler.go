package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// Handler 使い方
type Handler struct {
	E  echo.Echo
	DB *sqlx.DB
}
