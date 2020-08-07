package model

import (
	"time"
)

// requestではツイートした人と本文のみが渡される
type RequestTweet struct {
	UserID    string `json:"user_id,omitempty" db:"user_id"`
	TweetBody string `json:"tweet_body,omitempty" db:"tweet_body"`
}

// routerでuuid, DBでcreated_atを付けてクライアントに返す
type ResponseTweet struct {
	ID        string    `json:"id,omitempty" db:"id"`
	UserID    string    `json:"user_id,omitempty" db:"user_id"`
	TweetBody string    `json:"tweet_body,omitempty" db:"tweet_body"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func GetTweets() ([]ResponseTweet, error) {
	tweets := []ResponseTweet{}
	err := db.Select(&tweets, "SELECT * FROM tweets ORDER BY created_at DESC LIMIT 30")
	return tweets, err
}

func InsertTweet(id string, tweet *RequestTweet) error {
	_, err := db.Exec("INSERT INTO tweets (id, user_id, tweet_body) VALUES (?, ?, ?)", id, tweet.UserID, tweet.TweetBody)

	return err
}

// 送信したツイートをDBに格納して、created_atを付けてからもう一度クライアントに送り返す
func GetPostedTweet(id string) (*ResponseTweet, error) {
	tweet := new(ResponseTweet)
	err := db.Get(tweet, "SELECT * FROM tweets WHERE id = ?", id)
	return tweet, err
}
