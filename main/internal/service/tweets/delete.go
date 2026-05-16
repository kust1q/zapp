package tweets

import (
	"context"
	"errors"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

func (s *service) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := s.media.DeleteTweetMedia(ctx, tweetID, userID); err != nil && !errors.Is(err, errs.ErrTweetMediaNotFound) {
		logrus.WithFields(logrus.Fields{
			"tweet_id": tweetID,
			"err":      err,
		}).Warnf("failed to delete tweet media")
	}

	if err := s.db.DeleteTweet(ctx, userID, tweetID); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event := events.TweetDeleted{
			EventType: events.TweetDeleteEvent,
			ID:        tweetID,
		}
		if err := s.producer.Publish(bgCtx, events.TopicTweet, event); err != nil {
			logrus.WithError(err).Error("failed to publish tweet.deleted")
		}
	}()

	return nil
}
