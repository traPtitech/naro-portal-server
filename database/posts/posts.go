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
	ID          uint64    `json:"id,omitempty"  db:"id"`
	Content     uint64    `json:"content"  db:"content"`
	Desc        string    `json:"desc"  db:"desc"`
	CreatedUser string    `json:"created_user"  db:"created_user"`
	CreatedDate time.Time `json:"created_date,omitempty"  db:"created_date"`
}

func CreatePost(desc string, createdUser string) *Post {
	return &Post{
		Desc:        desc,
		CreatedUser: createdUser,
		CreatedDate: time.Now(),
	}
}

func CreatePostWithContent(desc string, createdUser string, content uint64) *Post {
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
		`SELECT * FROM `+p.tableName+` WHERE id = ?`,
		id,
	)
	return
}

func (p *PostDB) GetPosts(posts []Post) (err error) {
	err = p.db.Select(
		&posts,
		`SELECT * FROM `+p.tableName,
	)
	return
}

func (p *PostDB) AddPost(post *Post) (err error) {
	_, err = p.db.NamedExec(
		`INSERT INTO `+p.tableName+` (content, desc, created_user, created_date) VALUES (:content, :desc, :created_user, created_date)`,
		post,
	)
	return
}
