package router

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type ResponseData struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func return400(c echo.Context, err error) error {
	return c.String(http.StatusBadRequest, err.Error())
}

func return500(c echo.Context, name string, err error) error {
	fmt.Printf("%s: %v\n", name, err)
	return c.NoContent(http.StatusInternalServerError)
}

func returnErrorJSON(c echo.Context, content string) error {
	return c.JSON(http.StatusBadRequest, ResponseData{
		Type:    "Error",
		Content: content,
	})
}

func returnSuccessJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseData{
		Type:    "Success",
		Content: "Success",
	})
}
