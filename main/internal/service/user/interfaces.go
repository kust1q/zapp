package user

import (
	"context"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	db interface {
		GetUserByID(ctx context.Context, userID int) (*entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
		UpdateUserBio(ctx context.Context, userID int, bio string) error
		DeleteUser(ctx context.Context, userID int) error
		FollowToUser(ctx context.Context, followerID, followingID int, createdAt time.Time) (*entity.Follow, error)
		UnfollowUser(ctx context.Context, followerID, followingID int) error
		GetFollowersIds(ctx context.Context, username string, limit, offset int) ([]int, error)
		GetFollowingsIds(ctx context.Context, username string, limit, offset int) ([]int, error)
		GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error)
		GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error)
	}

	mediaService interface {
		GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error)
		DeleteAvatar(ctx context.Context, userID int) error

		GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error)

		DeleteMediasByUserID(ctx context.Context, userID int) error
	}

	eventProducer interface {
		Publish(ctx context.Context, topic string, event any) error
	}
)
