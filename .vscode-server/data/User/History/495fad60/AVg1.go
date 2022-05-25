package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
	Code        string `json:"code,omitempty"  db:"Code"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")
	var city City
	var kuni City
	cityhosii := os.Args[1]
	if err := db.Get(&city, "SELECT * FROM city WHERE Name='"+cityhosii+"'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = %s", "Tokyo")
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}
	kunihosii := city.CountryCode
	if err := db.Get(&kuni, "SELECT Name,Code,Population FROM country WHERE Code='"+kunihosii+"'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = %s", "Tokyo")
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	fmt.Printf(city.Name+"の人口は%d人です\n", city.Population)
	fmt.Printf(city.Name+"は"+kuni.Name+"の街で、国の人口に対しての割合は"+"%d %である。\n", int(math.Round((float64(city.Population)/float64(kuni.Population))*100)))
}
