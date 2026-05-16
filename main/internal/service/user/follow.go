package user

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func (s *service) FollowToUser(ctx context.Context, followerID, followingID int) (*entity.Follow, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if followerID == followingID {
		return nil, fmt.Errorf("impossible to subscribe to yourself")
	}
	return s.db.FollowToUser(ctx, followerID, followingID, time.Now())
}

func (s *service) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.db.UnfollowUser(ctx, followerID, followingID)
}

func (s *service) GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	followersIDs, err := s.db.GetFollowersIds(ctx, username, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers ids: %w", err)
	}

	users := make([]entity.SmallUser, 0, len(followersIDs))
	for _, id := range followersIDs {
		user, err := s.db.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}

		avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get avatar")
		}

		users = append(users, entity.SmallUser{
			ID:        user.ID,
			Username:  user.Username,
			AvatarUrl: avatarURL,
		})
	}

	return users, nil
}

func (s *service) GetFollowings(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	followingsIDs, err := s.db.GetFollowingsIds(ctx, username, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get followers ids: %w", err)
	}

	users := make([]entity.SmallUser, 0, len(followingsIDs))
	for _, id := range followingsIDs {
		user, err := s.db.GetUserByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}

		avatarURL, err := s.media.GetAvatarUrlByUserID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get avatar")
		}

		users = append(users, entity.SmallUser{
			ID:        user.ID,
			Username:  user.Username,
			AvatarUrl: avatarURL,
		})
	}

	return users, nil
}
