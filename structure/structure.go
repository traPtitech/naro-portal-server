package structure

type DateInGame struct {
	// ゲーム内日時

	Year  int // 1or2(1年目、1月以降)
	Month int
	Date  int
}

type CoopStatus struct {
}

type UserStatus struct {
	UserName  string
	UserId    string
	Password  string
	UserImage string
	DateInGame
	CoopStatus
}
