package tweets

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/sirupsen/logrus"
)

func (s *service) CreateTweet(ctx context.Context, tweet *entity.Tweet) (res *entity.Tweet, err error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var createdTweet *entity.Tweet
	createdTweet, err = s.db.CreateTweetTx(ctx, tx, tweet)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	var mediaUrl string
	if tweet.File != nil {
		mediaUrl, err = s.media.UploadAndAttachTweetMediaTx(ctx, createdTweet.ID, tweet.File.File, tweet.File.Header.Filename, tx)
		if err != nil {
			return nil, err
		}
	}

	createdTweet.MediaUrl = mediaUrl

	res, err = s.BuildEntityTweetToResponse(ctx, createdTweet)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event := events.TweetEvent{
			EventType: events.TweetCreateEvent,
			ID:        createdTweet.ID,
			Content:   createdTweet.Content,
			UserID:    createdTweet.Author.ID,
			Username:  createdTweet.Author.Username,
		}

		if pErr := s.producer.Publish(ctx, events.TopicTweet, event); pErr != nil {
			logrus.WithError(pErr).Error("failed to publish tweet.created")
		}
	}()

	return res, nil
}
