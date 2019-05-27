package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *gorm.DB
)

type DateInGame struct {
	// ゲーム内日時

	// 1or2(1年目、1月以降)
	Year  int `json:"year"`
	Month int `json:"month"`
	Date  int `json:"date"`
}

type CoopStatus struct {
}

type UserStatus struct {
	UserName  string `json:"user_name"`
	UserID    string `json:"user_id"`
	Password  string `json:"password"`
	UserImage string `json:"user_image"`
	DateInGame
	CoopStatus
}

func CreateUserStatusTable() {
	db.CreateTable(&UserStatus{})
}
