package model

func Counter(username string) (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE username=?", username)
	return count, err
}

type UserWithHashedPass struct {
	Username   string `json:"username,omitempty" db:"username"`
	HashedPass string `json:"-" db:"hashed_pass"`
}

func InsertUserWithHashedPass(username string, hashedPass []byte) error {
	_, err := db.Exec("INSERT INTO users (username, hashed_pass) VALUES (?, ?)", username, hashedPass)
	return err
}

func SelectUser(username string) (*UserWithHashedPass, error) {
	savedUser := new(UserWithHashedPass)
	err := db.Get(&savedUser, "SELECT * FROM users WHERE username=?", username)
	return savedUser, err
}
