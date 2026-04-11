package kafka

import (
	"context"

	"github.com/kust1q/Zapp/search/internal/domain/entity"
)

type (
	searchService interface {
		IndexTweet(ctx context.Context, tweet *entity.Tweet) error
		IndexUser(ctx context.Context, user *entity.User) error
		DeleteTweet(ctx context.Context, tweetID int) error
		DeleteUserWithTweets(ctx context.Context, userID int) error
	}
)
