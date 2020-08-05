package model

func Counter(username string) (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", username)
	return count, err
}

type UserWithHashedPass struct {
	Username   string `json:"username,omitempty" db:"id"`
	HashedPass string `json:"-" db:"hashed_pass"`
}

func InsertUserWithHashedPass(username string, hashedPass []byte) error {
	_, err := db.Exec("INSERT INTO users (id, hashed_pass) VALUES (?, ?)", username, hashedPass)
	return err
}

func SelectUser(username string) (UserWithHashedPass, error) {
	savedUser := UserWithHashedPass{}
	err := db.Get(&savedUser, "SELECT * FROM users WHERE id=?", username)
	return savedUser, err
}
