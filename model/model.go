package model

import (
	"os"

	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func EstablishConecction() {
	databaseURL := os.Getenv("DATABASE_URL")
	_db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		panic("failed to connect database")
	}
	db = _db
}
