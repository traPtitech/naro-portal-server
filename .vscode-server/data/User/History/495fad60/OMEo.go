package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type MyError struct {
	Error string `json:"error"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	e := echo.New()

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/add", addHandler)

	e.Start(":11000")

}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	var city City
	if err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
		log.Printf("No Such City Name=%s", cityName)
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	return c.JSON(http.StatusOK, city)
}


func addHandler(c echo.Context) error {
	data := new(City)
	err := c.Bind(data)

	fmt.Print(data)

	if err != nil {
		msg := MyError{Error: "Bad Request"}
		return c.JSON(400, msg)
	}

	if err := db.Get(&city, "SELECT * FROM city WHERE Name='"+cityhosii+"'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = %s", "Tokyo")
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	INSERT INTO city (name, countryCode, district, ...) VALUES (値1, 値2, 値3, ...)

	jsonAnswer := data.Right + data.Left
	Ans := Answer{jsonAnswer}
	return c.JSON(http.StatusOK, Ans)
}
