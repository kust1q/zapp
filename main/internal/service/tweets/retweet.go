package tweets

import (
	"context"
	"fmt"
	"time"
)

func (s *service) CreateRetweet(ctx context.Context, userID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.db.Retweet(ctx, userID, tweetID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create retweet")
	}
	return nil
}

func (s *service) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := s.db.DeleteRetweet(ctx, userID, retweetID)
	if err != nil {
		return fmt.Errorf("failed to delete retweet")
	}
	return nil
}
