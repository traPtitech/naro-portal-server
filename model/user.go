package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	Password  []byte `json:"password"` // ハッシュ化済み
	UserImage string `json:"user_image"`
	DateInGame
	CoopStatus
}

type DataForSignUpAndSignIn struct {
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	UserImage string `json:"user_image"`
}

func CreateUserStatusTable() {
	db.CreateTable(&UserStatus{})
}

func AddNewUserStatus(userData DataForSignUpAndSignIn) {
	userStatus := UserStatus{}
	userStatus.UserName = userData.UserName
	userStatus.UserImage = userData.UserImage
	userStatus.DateInGame = DateInGame{Year: 1, Month: 4, Date: 9}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password")
	}
	userStatus.Password = hashedPassword

	u, err := uuid.NewRandom()
	if err != nil {
		panic("failed to make uuid")
	}
	userStatus.UserID = u.String()

	db.Create(&userStatus)
}
