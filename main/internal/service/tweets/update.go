package tweets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

func (s *service) UpdateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	exTweet, err := s.db.GetTweetById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTweetNotFound
		}
		return nil, fmt.Errorf("failed to get tweet by id: %w", err)
	}

	if exTweet.Author.ID != req.Author.ID {
		return nil, errs.ErrUnauthorizedUpdate
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	exTweet.Content = req.Content
	exTweet.UpdatedAt = time.Now()

	updatedTweet, err := s.db.UpdateTweet(ctx, exTweet)
	if err != nil {
		return nil, fmt.Errorf("failed to update tweet: %w", err)
	}

	var mediaUrl string
	if req.File != nil {
		if err := s.media.DeleteTweetMedia(ctx, req.ID, req.Author.ID); err != nil {
			logrus.WithFields(logrus.Fields{
				"tweet_id": req.ID,
				"error":    err,
			}).Warn("failed to delete old media record")
		}
		mediaUrl, err = s.media.UploadAndAttachTweetMediaTx(ctx, updatedTweet.ID, req.File.File, req.File.Header.Filename, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to upload new media: %w", err)
		}
	}

	updatedTweet.MediaUrl = mediaUrl

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event := events.TweetEvent{
			EventType: events.TweetUpdateEvent,
			ID:        updatedTweet.ID,
			Content:   updatedTweet.Content,
			UserID:    updatedTweet.Author.ID,
			Username:  updatedTweet.Author.Username,
		}
		if err := s.producer.Publish(cntx, events.TopicTweet, event); err != nil {
			logrus.WithError(err).Error("failed to publish tweet.updated")
		}
	}()

	return s.BuildEntityTweetToResponse(ctx, updatedTweet)
}
