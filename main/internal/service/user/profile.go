package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func (s *service) GetMe(ctx context.Context, userID, limit, offset int) (*entity.UserProfile, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return s.GetUserProfile(ctx, user.Username, limit, offset)
}

func (s *service) GetUserProfile(ctx context.Context, username string, limit, offset int) (*entity.UserProfile, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by user id: %w", err)
	}

	avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar by user id: %w", err)
	}

	user.AvatarUrl = avatarURL

	tweets, err := s.db.GetTweetsAndRetweetsByUsername(ctx, user.Username, limit, offset)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get tweets by username: %w", err)
		}
	}

	for i := range tweets {
		author, err := s.db.GetUserByID(ctx, tweets[i].Author.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tweet author: %w", err)
		}

		avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, tweets[i].Author.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user avatar: %w", err)
		}

		tweets[i].Author = &entity.SmallUser{
			ID:        author.ID,
			Username:  author.Username,
			AvatarUrl: avatarUrl,
		}

		mediaUrl, err := s.media.GetMediaUrlByTweetID(ctx, tweets[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tweet media url")
		}
		tweets[i].MediaUrl = mediaUrl

		counts, err := s.db.GetCounts(ctx, tweets[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get counters: %w", err)
		}
		tweets[i].Counters = counts
	}

	return &entity.UserProfile{
		User:   user,
		Tweets: tweets,
	}, nil
}
