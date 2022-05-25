package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type jsonData struct {
	Number int    `json:"number,omitempty"`
	String string `json:"string,omitempty"`
	Bool   bool   `json:"bool,omitempty"`
}

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/pikachu", func(c echo.Context) error {
		return c.String(http.StatusOK, "はじめましてピカチュウです！\n#ケモナーは一般性癖.\n")

	})

	e.GET("/hello/:username", helloHandler)

	e.GET("/json", jsonHandler)

	e.POST("/post", postHandler)

	e.GET("/ping", pingHandler)

	e.GET("/fizzbuzz", fizzbuzzHandler)

	e.Logger.Fatal(e.Start(":11000"))
}

func jsonHandler(c echo.Context) error {
	res := jsonData{
		Number: 10,
		String: "hoge",
		Bool:   false,
	}

	return c.JSON(http.StatusOK, &res)
}

func postHandler(c echo.Context) error {
	data := new(jsonData)
	err := c.Bind(data)

	if err != nil { // エラーが発生した際
		// fmt.Sprintf("%+v", data): dataをstringに変換
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%+v", data))
	}
	return c.JSON(http.StatusOK, data)
}

func helloHandler(c echo.Context) error {
	userID := c.Param("username")
	return c.String(http.StatusOK, "Hello, "+userID+".\n")
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong\n")
}

func fizzbuzzHandler(c echo.Context) error {
	count, _ := strconv.Atoi(c.QueryParam("count"))
	ans := ""
	for i := 0; i < count; i++ {
		if i%15 == 0 {
			ans += "FizzBuzz"
		} else {
			if i%3 == 0 {
				ans += "Fizz"
			} else {
				if i%5 == 0 {
					ans += "Buzz"
				} else {
					ans += strconv.Itoa(i)
				}
			}
		}
		ans += "\n"
	}
	return c.String(http.StatusOK, ans)
}
