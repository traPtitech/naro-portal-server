package posts

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type PostDB struct {
	db        *sqlx.DB
	tableName string
}

type Post struct {
	ID          int64     `json:"id,omitempty"  db:"Id"`
	Content     string    `json:"content"  db:"Content"`
	Desc        string    `json:"desc"  db:"Desc"`
	CreatedUser string    `json:"created_user,omitempty"  db:"CreatedUser"`
	CreatedDate time.Time `json:"created_date,omitempty"  db:"CreatedDate"`
}

func CreatePost(content string, desc string, createdUser string) *Post {
	return &Post{
		Content:     content,
		Desc:        desc,
		CreatedUser: createdUser,
		CreatedDate: time.Now(),
	}
}

func CreatePostDB(db *sqlx.DB) *PostDB {
	return &PostDB{
		db:        db,
		tableName: "post",
	}
}

func (p *PostDB) GetPost(id string, post *Post) (err error) {
	err = p.db.Get(
		&post,
		`SELECT * FROM `+p.tableName+` WHERE Id = ?`,
		id,
	)
	return
}

func (p *PostDB) AddPost(post *Post) (err error) {
	_, err = p.db.NamedExec(
		`INSERT INTO `+p.tableName+` (Content, Desc, CreatedUser, CreatedDate) VALUES (:Content, :Desc, :CreatedUser, CreatedDate)`,
		post,
	)
	return
}
