package model

import (
	"errors"
	"os"

	"github.com/antonlindstrom/pgstore"
	"github.com/jinzhu/gorm"
)

var (
	db          *gorm.DB
	databaseURL string
)

func EstablishConecction() (*gorm.DB, error) {
	databaseURL = os.Getenv("DATABASE_URL")
	_db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.New("faild to connect to DB")
	}
	db = _db

	return db, nil
}

func StoreForSession() (*pgstore.PGStore, error) {
	store, err := pgstore.NewPGStore(databaseURL)

	return store, err
}
