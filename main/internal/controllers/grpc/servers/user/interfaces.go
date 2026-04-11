package usergrpc

import (
	"context"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	userService interface {
		GetUserByID(ctx context.Context, userID int) (*entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
		GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
		GetFollowings(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
		GetUserProfile(ctx context.Context, username string, limit, offset int) (*entity.UserProfile, error)
		DeleteUser(ctx context.Context, userID int) error
	}
)
