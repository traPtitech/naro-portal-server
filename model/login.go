package model

func Counter(userID string) (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", userID)
	return count, err
}

type UserWithHashedPass struct {
	ID         string `db:"id"`
	HashedPass string `db:"hashed_pass"`
}

func InsertUserWithHashedPass(userID string, hashedPass []byte) error {
	_, err := db.Exec("INSERT INTO users (id, hashed_pass) VALUES (?, ?)", userID, hashedPass)
	return err
}

func SelectUser(userID string) (UserWithHashedPass, error) {
	savedUser := UserWithHashedPass{}
	err := db.Get(&savedUser, "SELECT * FROM users WHERE id=?", userID)
	return savedUser, err
}
