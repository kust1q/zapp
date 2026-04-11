package conv

import (
	"time"

	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/response"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

// Requests
func FromTweetRequestToDomain(userID int, parent_tweet_id *int, file *entity.File, req *request.Tweet) *entity.Tweet {
	if req == nil {
		return nil
	}

	return &entity.Tweet{
		ParentTweetID: parent_tweet_id,
		Content:       req.Content,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		File:          file,
		Author: &entity.SmallUser{
			ID: userID,
		},
		Counters: nil,
	}
}

func FromTweetUpdateRequestToDomain(userID, tweetID int, file *entity.File, req *request.Tweet) *entity.Tweet {
	if req == nil {
		return nil
	}

	return &entity.Tweet{
		ID:        tweetID,
		Content:   req.Content,
		UpdatedAt: time.Now(),
		File:      file,
		Author: &entity.SmallUser{
			ID: userID,
		},
	}
}

// Responses
func FromDomainToTweetResponse(tweet *entity.Tweet) *response.Tweet {
	if tweet == nil {
		return nil
	}

	responseTweet := &response.Tweet{
		ID:            tweet.ID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		ParentTweetID: tweet.ParentTweetID,
		MediaUrl:      tweet.MediaUrl,
		Author:        FromDomainToSmallUserResponse(tweet.Author),
	}

	if tweet.Counters != nil {
		responseTweet.Counters = &response.Counters{
			ReplyCount:   tweet.Counters.ReplyCount,
			RetweetCount: tweet.Counters.RetweetCount,
			LikeCount:    tweet.Counters.LikeCount,
		}
	}

	return responseTweet
}

func FromDomainToTweetListResponse(tweets []entity.Tweet) []response.Tweet {
	if tweets == nil {
		return nil
	}

	res := make([]response.Tweet, 0, len(tweets))
	for _, t := range tweets {
		tweetResponse := FromDomainToTweetResponse(&t)
		if tweetResponse != nil {
			res = append(res, *tweetResponse)
		}
	}
	return res
}
