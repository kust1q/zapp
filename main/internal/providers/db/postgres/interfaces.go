package postgres

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/providers/db/models"
)

type (
	cache interface {
		SetUser(ctx context.Context, user *models.User) error
		GetUserByID(ctx context.Context, id int) (*models.User, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		ExistsByUsername(ctx context.Context, username string) (bool, error)
		ExistsByEmail(ctx context.Context, email string) (bool, error)
		InvalidateUser(ctx context.Context, userID int) error

		SetTweet(ctx context.Context, tweet *models.Tweet) error
		GetTweet(ctx context.Context, tweetID int) (*models.Tweet, error)
		MGetTweets(ctx context.Context, tweetIDs []int) (map[int]*models.Tweet, error)
		InvalidateTweet(ctx context.Context, tweetID int) error

		SetUserTweetIDs(ctx context.Context, username string, ids []int) error
		GetUserTweetIDs(ctx context.Context, username string) ([]int, error)
		InvalidateUserTweets(ctx context.Context, username string) error

		SetReplyIDs(ctx context.Context, parentTweetID int, ids []int) error
		GetReplyIDs(ctx context.Context, parentTweetID int) ([]int, error)
		InvalidateReplies(ctx context.Context, parentTweetID int) error

		SetTweetLikerIDs(ctx context.Context, tweetID int, userIDs []int) error
		GetTweetLikerIDs(ctx context.Context, tweetID int) ([]int, error)
		InvalidateTweetLikers(ctx context.Context, tweetID int) error

		SetTweetCounters(ctx context.Context, tweetID int, counters *models.Counters) error
		GetTweetCounters(ctx context.Context, tweetID int) (*models.Counters, error)
		InvalidateTweetCounters(ctx context.Context, tweetID int) error
	}
)
