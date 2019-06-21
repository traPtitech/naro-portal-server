package model

import (
	"os"
	//"github.com/labstack/gommon/log"
	"github.com/WistreHosshii/naro-portal-server/model/mystruct"
	"github.com/jmoiron/sqlx"

	"fmt"
)

var (
	db *sqlx.DB
)

//dbとの通信
func EstablishConnection() error {
	_db, err := sqlx.Connect("postgres", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	db = _db
	return err
}

func GetUserCount(userName string) (int, error) {
	var count int
	var err error
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE user_name=?", userName)
	return count, err
}

func GetUserName(userName string) (mystruct.User, error) {
	user := mystruct.User{}
	err := db.Get(&user, "SELECT FROM users WHERE user_name=?", userName)
	return user, err
}

func ExecUserInfo(userName string, hashedPass []byte, id string) error {
	_, err := db.Exec("INSERT INTO users (user_name, hashed_pass, id) VALUES (?, ?, ?)", userName, hashedPass, id)
	return err
}
