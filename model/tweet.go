package model

import (
	"time"
)

type JsonTweet struct {
	ID        string    `json:id,omitempty db:id`
	TweetBody string    `json:tweet_body,omitempty db:tweet_body`
	Author    string    `json:author,omitempty db:author`
	CreatedAt time.Time `json:created_at,omitempty db:created_at`
}

func SelectTweet() ([]JsonTweet, error) {
	tweets := []JsonTweet{}
	err := db.Select(&tweets, "SELECT * FROM tweets ORDER BY created_at DESC LIMIT 30")
	return tweets, err
}

func InsertTweet(tweet *JsonTweet) error {
	_, err := db.Exec("INSERT INTO tweets (id, tweet_body, author, created_at) VALUES (?, ?, ?, ?)",
		tweet.ID, tweet.TweetBody, tweet.Author, tweet.CreatedAt)

	return err
}
