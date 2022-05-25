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

	e.POST("/add", addHandler)

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
	countmae := c.QueryParam("count")
	var count int
	var err error
	if countmae == "" {
		count = 30
	} else {
		count, err = strconv.Atoi(countmae)
	}

	if err != nil { // エラーが発生した際
		// fmt.Sprintf("%+v", data): dataをstringに変換
		return c.String(400, "Bad Request")
	}
	ans := ""
	for i := 1; i <= count; i++ {
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

func addHandler(c echo.Context) error {
	data := new(jsonData)
	err := c.Bind(data)

	if err != nil {
		return c.JSON(400, {"error":"Bad Request"})
	}
	
	count1, err1 = strconv.Atoi(data["right"])
	if err1 != nil {
		return c.JSON(400, {"error":"Bad Request"})
	}

	count2, err2 = strconv.Atoi(data["left"])
	if err2 != nil {
		return c.JSON(400, {"error":"Bad Request"})
	}
	
	jsonAnswer = {count1 + count2}

	return c.JSON(http.StatusOK, data)
}
