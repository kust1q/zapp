package user

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func (s *service) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar by user id: %w", err)
	}

	user.AvatarUrl = avatarUrl

	return user, nil
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar by user id: %w", err)
	}

	user.AvatarUrl = avatarURL

	return user, nil
}
