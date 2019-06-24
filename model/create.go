package model

//Create DBのテーブル作成
func Create() error {
	_, err := Db.Exec(`CREATE TABLE "favorite" (favo_ID TEXT,user_ID VARCHAR(36),tweet_ID VARCHAR(36),created_at DATETIME)`)
	if err != nil {
		panic(err)
	}
	_, err = Db.Exec(`CREATE TABLE "pin" (pin_ID VARCHAR(36),user_ID VARCHAR(36),tweet_ID VARCHAR(36))`)
	if err != nil {
		panic(err)
	}
	_, err = Db.Exec(`CREATE TABLE "tweet" (tweet_ID VARCHAR(36),user_ID VARCHAR(36),tweet TEXT,created_at DATETIME,favo_num INT(11))`)
	if err != nil {
		panic(err)
	}
	_, err = Db.Exec(`CREATE TABLE "user" (name VARCHAR(30),ID VARCHAR(36),password VARCHAR(200))`)
	if err != nil {
		panic(err)
	}
	_, err = Db.Exec(`CREATE TABLE "session" (id INT(11),session_data LONGBLOB,created_on TIMESTAMP,modified_on TIMESTAMP,expires_on TIMESTAMP)`)
	if err != nil {
		panic(err)
	}
	return err
}
