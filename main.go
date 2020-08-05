package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Tweet struct {
	ID        string    `json:"id,omitempty"  db:"id"`
	TweetBody string    `json:"name,omitempty"  db:"tweet_body"`
	Author    string    `json:"countryCode,omitempty"  db:"author"`
	CreatedAt time.Time `json:"district,omitempty"  db:"created_at"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")
	tweet := Tweet{}
	err = db.Get(&tweet, "SELECT author FROM tweets WHERE id='Tokyo'")
	if err != nil {
		fmt.Printf("db error")
		os.Exit(1)
	}

	fmt.Printf("%s\n", tweet.Author)
}
