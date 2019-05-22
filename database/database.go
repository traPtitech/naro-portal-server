package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/auths"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/posts"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/sessions"
	"github.com/sapphi-red/webengineer_naro-_7_server/database/users"
	"github.com/srinathgs/mysqlstore"
	"log"
	"os"
)

var (
	db *sqlx.DB
	SessionStore *mysqlstore.MySQLStore
	Sessions *sessions.SessionDB
	Auths *auths.AuthDB
	Users *users.UserDB
	Posts *posts.PostDB
)

func ConnectDB() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	initialize()
}

func initialize() {
	SessionStore = sessions.CreateStore(db)
	Sessions = sessions.CreateSessionDB()
	Auths = auths.CreateAuthDB(db)
	Users = users.CreateUserDB(db)
	Posts = posts.CreatePostDB(db)
}
