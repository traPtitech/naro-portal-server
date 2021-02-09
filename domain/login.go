package domain

type LoginRequestBody struct {
	UserName string `json:"userName"  db:"userName"`
	Password string `json:"password"  db:"password"`
}