package search

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	mediaService interface {
		GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error)
		GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error)
	}

	tweetService interface {
		BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
	}

	searchStorage interface {
		GetTweetsByIDs(ctx context.Context, ids []int) ([]entity.Tweet, error)
		GetUsersByIDs(ctx context.Context, ids []int) ([]entity.User, error)
	}

	searchProvider interface {
		SearchTweets(ctx context.Context, query string) ([]int, error)
		SearchUsers(ctx context.Context, query string) ([]int, error)
	}
)
