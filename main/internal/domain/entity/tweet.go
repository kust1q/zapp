package entity

import (
	"time"
)

type (
	Tweet struct {
		ID            int
		ParentTweetID *int
		Content       string
		CreatedAt     time.Time
		UpdatedAt     time.Time
		MediaUrl      string
		Author        *SmallUser
		File          *File
		Counters      *Counters
	}

	Counters struct {
		ReplyCount   int
		RetweetCount int
		LikeCount    int
	}

	Retweet struct {
		ID        int
		UserID    int
		TweetID   int
		CreatedAt time.Time
	}

	Like struct {
		UserID  int
		TweetID int
	}
)
