package tweetgrpc

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	tweetService interface {
		GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error)
		GetRepliesToTweet(ctx context.Context, tweetID, limit, offset int) ([]entity.Tweet, error)
		GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error)
		GetLikes(ctx context.Context, tweetID, limit, offset int) ([]entity.SmallUser, error)
	}
)
