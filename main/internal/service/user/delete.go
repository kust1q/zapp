package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

func (s *service) DeleteUser(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.media.DeleteAvatar(ctx, userID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Warnf("failed to delete user avatar")
	}

	if err := s.media.DeleteMediasByUserID(ctx, userID); err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err,
		}).Warnf("failed to delete all user tweet medias")
	}

	if err := s.db.DeleteUser(ctx, userID); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event := events.UserDeleted{
			EventType: events.TweetDeleteEvent,
			ID:        userID,
		}
		if err := s.producer.Publish(ctx, events.TopicTweet, event); err != nil {
			logrus.WithError(err).Error("failed to publish user.deleted")
		}
	}()

	return nil
}
