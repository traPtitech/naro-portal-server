package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/srinathgs/mysqlstore"
)

const storeName = "session"

var store *mysqlstore.MySQLStore

type SessionDB struct {
	storeName string
}

func CreateStore(db *sqlx.DB) *mysqlstore.MySQLStore {
	_store, err := mysqlstore.NewMySQLStoreFromConnection(
		db.DB,
		storeName,
		"/",
		60*60*24*14,
		[]byte("secret-token"),
	)
	if err != nil {
		panic(err)
	}

	store = _store
	return store
}

func CreateSessionDB() *SessionDB {
	return &SessionDB{
		storeName: storeName,
	}
}

func (s *SessionDB) Get(c echo.Context) (*sessions.Session, error) {
	return session.Get(s.storeName, c)
}

func (s *SessionDB) SetID(c echo.Context, id string) error {
	sess, err := s.Get(c)
	if err != nil {
		return err
	}

	sess.Values["id"] = id
	sess.Save(c.Request(), c.Response())
	return nil
}

func (s *SessionDB) Destroy(c echo.Context) error {
	sess, err := s.Get(c)
	if err != nil {
		return err
	}

	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	return nil
}
