package feed

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	db interface {
		GetUserByID(ctx context.Context, userID int) (*entity.User, error)
		GetFollowingsIds(ctx context.Context, username string, limit, offset int) ([]int, error)

		GetFeedByAuthorsIds(ctx context.Context, userIDs []int, limit, offset int) ([]entity.Tweet, error)
		GetAllTweets(ctx context.Context, limit, offset int) ([]entity.Tweet, error)
	}

	tweetService interface {
		BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error)
	}
)
