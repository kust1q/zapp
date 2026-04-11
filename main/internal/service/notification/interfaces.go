package notification

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	hubProvider interface {
		SendNotification(notification *entity.Notification)
	}

	db interface {
		GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error)
		GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	}
)
