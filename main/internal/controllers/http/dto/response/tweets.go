package response

import (
	"time"
)

type (
	Tweet struct {
		ID            int        `json:"id"`
		Content       string     `json:"content"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		ParentTweetID *int       `json:"parent_tweet_id,omitempty"`
		MediaUrl      string     `json:"media_url"`
		Author        *SmallUser `json:"author"`
		Counters      *Counters  `json:"counters"`
	}

	Counters struct {
		ReplyCount   int `json:"reply_count"`
		RetweetCount int `json:"retweet_count"`
		LikeCount    int `json:"like_count"`
	}

	// For docs
	TweetList struct {
		Tweets []Tweet `json:"tweets"`
	}
)
