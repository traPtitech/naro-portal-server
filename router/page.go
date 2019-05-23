package router

import (
	//"fmt"
	"net/http"
	"github.com/labstack/echo"
)

func Pong(c echo.Context)error{
	return c.String(http.StatusOK, "pong")
}

