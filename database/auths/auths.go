package auths

import (
	"github.com/jmoiron/sqlx"
)

type AuthDB struct {
	db        *sqlx.DB
	tableName string
}

type AuthUser struct {
	ID         string `json:"id"  db:"id"`
	HashedPass string `json:"hashed_pass"  db:"hashed_pass"`
}

func CreateAuthDB(db *sqlx.DB) *AuthDB {
	return &AuthDB{
		db:        db,
		tableName: "auth",
	}
}

func (a *AuthDB) GetUser(id string, user *AuthUser) (err error) {
	err = a.db.Get(
		&user,
		`SELECT * FROM `+a.tableName+` WHERE Id = ?`,
		id,
	)
	return
}

func (a *AuthDB) GetUserExistance(id string) (res bool, err error) {
	var count int
	err = a.db.Get(&count, `SELECT COUNT(*) FROM `+a.tableName+` WHERE Id = ?`, id)
	res = count > 0
	return
}

func (a *AuthDB) AddUser(id string, hashedPass []byte) (err error) {
	_, err = a.db.Exec(
		`INSERT INTO `+a.tableName+` (Id, HashedPass) VALUES (?, ?)`,
		id,
		hashedPass,
	)
	return
}

func (a *AuthDB) DeleteUser(id string) (err error) {
	_, err = a.db.Exec(
		`DELETE FROM `+a.tableName+` WHERE Id = ?`,
		id,
	)
	return
}
