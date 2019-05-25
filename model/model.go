package model

import (
	"github.com/pborman/uuid"
	"time"
)

type User struct{					//Userの構造体
	UserID				uuid.UUID
	UserName			string
	UserPassword		string
	UserInfo			string
	UserIconURL			string
}

type Favorite struct{				//Favoriteの構造体
	FavoriteID			uuid.UUID
	MessageID			uuid.UUID
}

type Tweet struct{					//Tweetの構造体
	TweetID				uuid.UUID
	UserID				uuid.UUID
	Tweet				string
	CreatedTime			time.Time
	FavoNum				int
}

type Pin struct{					//Pinの構造体
	PinID				uuid.UUID
	UserID				uuid.UUID
	MessageID			uuid.UUID
}