package model

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	//Db Establishの構造体
	Db *sqlx.DB
)

//Establish データベースに接続
func Establish() error {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("database"), os.Getenv("username"), os.Getenv("password"), os.Getenv("hostname")))
	if err != nil {
		return err
	}
	Db = _db
	return nil
}
