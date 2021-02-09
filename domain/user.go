package domain

type User struct {
	UserName   string `json:"userName"  db:"userName"`
	HashedPass string `json:"hashedPass,omitempty"  db:"hashedPass"`
}
