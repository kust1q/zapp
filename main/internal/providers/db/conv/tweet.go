package conv

import (
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func FromDomainToTweetModel(tweet *entity.Tweet) *models.Tweet {
	if tweet == nil {
		return nil
	}

	return &models.Tweet{
		ID:            tweet.ID,
		UserID:        tweet.Author.ID,
		ParentTweetID: tweet.ParentTweetID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
	}
}

func FromTweetModelToDomain(tweet *models.Tweet) *entity.Tweet {
	if tweet == nil {
		return nil
	}

	return &entity.Tweet{
		ID:            tweet.ID,
		ParentTweetID: tweet.ParentTweetID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		Author: &entity.SmallUser{
			ID: tweet.UserID,
		},
	}
}

func FromTweetModelToDomainList(tweetsModels []models.Tweet) []entity.Tweet {
	tweets := make([]entity.Tweet, 0, len(tweetsModels))
	for _, tweet := range tweetsModels {
		tweets = append(tweets, *FromTweetModelToDomain(&tweet))
	}
	return tweets
}

func FromDomainToRetweetModel(retweet *entity.Retweet) *models.Retweet {
	if retweet == nil {
		return nil
	}

	return &models.Retweet{
		ID:        retweet.ID,
		UserID:    retweet.UserID,
		TweetID:   retweet.TweetID,
		CreatedAt: retweet.CreatedAt,
	}
}

func FromRetweetModelToDomain(retweet *models.Retweet) *entity.Retweet {
	if retweet == nil {
		return nil
	}

	return &entity.Retweet{
		ID:        retweet.ID,
		UserID:    retweet.UserID,
		TweetID:   retweet.TweetID,
		CreatedAt: retweet.CreatedAt,
	}
}

func FromDomainToLikeModel(like *entity.Like) *models.Like {
	if like == nil {
		return nil
	}

	return &models.Like{
		UserID:  like.UserID,
		TweetID: like.TweetID,
	}
}

func FromLikeModelToDomain(like *models.Like) *entity.Like {
	if like == nil {
		return nil
	}

	return &entity.Like{
		UserID:  like.UserID,
		TweetID: like.TweetID,
	}
}

func FromLikeModelToDomainList(likes []models.Like) []entity.Like {
	res := make([]entity.Like, 0, len(likes))
	for _, like := range likes {
		res = append(res, *FromLikeModelToDomain(&like))
	}
	return res
}

func FromCountersModelToDomain(counters *models.Counters) *entity.Counters {
	return &entity.Counters{
		ReplyCount:   counters.ReplyCount,
		RetweetCount: counters.RetweetCount,
		LikeCount:    counters.LikeCount,
	}
}
