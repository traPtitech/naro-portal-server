package users

import (
	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	db        *sqlx.DB
	tableName string
}

type User struct {
	ID   string `json:"id"  db:"id"`
	Name string `json:"name"  db:"name"`
}

func CreateUserDB(db *sqlx.DB) *UserDB {
	return &UserDB{
		db:        db,
		tableName: "user",
	}
}

func (u *UserDB) GetUser(id string, user *User) (err error) {
	err = u.db.Get(
		&user,
		`SELECT * FROM `+u.tableName+` WHERE Id = ?`,
		id,
	)
	return
}

func (u *UserDB) AddUser(user *User) (err error) {
	_, err = u.db.NamedExec(
		`INSERT INTO `+u.tableName+` (Id, Name) VALUES (:Id, :Name)`,
		user,
	)
	return
}
