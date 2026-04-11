package tweets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

func (s *service) LikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.db.LikeTweet(ctx, userID, tweetID)
	if err != nil && !errors.Is(err, errs.ErrTweetNotFound) {
		return fmt.Errorf("failed to like tweet: %w", err)
	} else if errors.Is(err, errs.ErrTweetNotFound) {
		return err
	}
	return nil
}

func (s *service) UnlikeTweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.db.UnLikeTweet(ctx, userID, tweetID)
	if err != nil {
		return fmt.Errorf("failed to unlike tweet: %w", err)
	}
	return nil
}

func (s *service) GetLikes(ctx context.Context, tweetID, limit, offset int) ([]entity.SmallUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	users, err := s.db.GetLikes(ctx, tweetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes: %w", err)
	}

	for i := range users {
		users[i].AvatarUrl, err = s.media.GetAvatarUrlByUserID(ctx, users[i].ID)
		if err != nil {
			logrus.WithError(err).Warn("failed to get avatar url")
		}
	}
	return users, nil
}
