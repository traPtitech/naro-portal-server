package model

//Create DBのテーブル作成
func Create() {
	Db.Exec(`CREATE TABLE "favorite" (favo_ID TEXT,user_ID VARCHAR(36),tweet_ID VARCHAR(36),created_at DATETIME)`)
	Db.Exec(`CREATE TABLE "pin" (pin_ID VARCHAR(36),user_ID VARCHAR(36),tweet_ID VARCHAR(36))`)
	Db.Exec(`CREATE TABLE "tweet" (tweet_ID VARCHAR(36),user_ID VARCHAR(36),tweet TEXT,created_at DATETIME,favo_num INT(11))`)
	Db.Exec(`CREATE TABLE "user" (name VARCHAR(30),ID VARCHAR(36),password VARCHAR(200))`)
	Db.Exec(`CREATE TABLE "session" (id INT(11),session_data LONGBLOB,created_on TIMESTAMP,modified_on TIMESTAMP,expires_on TIMESTAMP)`)
}
