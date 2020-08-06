package model

import (
	"time"
)

type JsonTweet struct {
	ID        string    `json:"id,omitempty" db:"id"`
	UserID    string    `json:"user_id,omitempty" db:"user_id"`
	TweetBody string    `json:"tweet_body,omitempty" db:"tweet_body"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func SelectTweet() ([]JsonTweet, error) {
	tweets := []JsonTweet{}
	err := db.Select(&tweets, "SELECT * FROM tweets ORDER BY created_at DESC LIMIT 30")
	return tweets, err
}

func InsertTweet(tweet *JsonTweet) error {
	_, err := db.Exec("INSERT INTO tweets (id, user_id, tweet_body, created_at) VALUES (?, ?, ?, ?)",
		tweet.ID, tweet.UserID, tweet.TweetBody, tweet.CreatedAt)

	return err
}
