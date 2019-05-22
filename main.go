package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/db", func(c echo.Context) error {
		connStr := "postgres://ssadsyncjczxby:8f647d3f6a031c4cb2d6fd97106053e259982e97c1205a6a2deff50e989e85e1@ec2-54-221-212-126.compute-1.amazonaws.com:5432/dcbbm7iv8usrv4"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		data, _ := db.Query("SELECT * FROM user_status;")
		return c.JSON(http.StatusOK, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	e.Start(":" + port) // ここを前述の通り自分のポートにすること
}
