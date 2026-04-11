package search

import (
	"context"

	"github.com/kust1q/Zapp/search/internal/domain/entity"
)

type (
	searchRepository interface {
		SearchTweets(ctx context.Context, query string) ([]int, error)
		SearchUsers(ctx context.Context, query string) ([]int, error)
		IndexTweet(ctx context.Context, tweet *entity.Tweet) error
		IndexUser(ctx context.Context, user *entity.User) error
		DeleteTweet(ctx context.Context, tweetID int) error
		DeleteUser(ctx context.Context, userID int) error
		DeleteTweetsByUserID(ctx context.Context, userID int) error
	}
)
