package auth

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	db interface {
		BeginTx(ctx context.Context) (*sql.Tx, error)
		CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
		GetUserByID(ctx context.Context, userID int) (*entity.User, error)
		UpdateUserPassword(ctx context.Context, userID int, password string) error
		DeleteUser(ctx context.Context, userID int) error
		UserExistsByUsername(ctx context.Context, username string) (bool, error)
		UserExistsByEmail(ctx context.Context, email string) (bool, error)
	}

	tokenStorage interface {
		StoreRefresh(ctx context.Context, refreshToken, userID string, ttl time.Duration) error
		GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (string, error)
		CloseAllSessions(ctx context.Context, userID string) error
		RemoveRefresh(ctx context.Context, refreshToken string) error

		StoreRecovery(ctx context.Context, token, userID string, ttl time.Duration) error
		GetUserIdByRecoveryToken(ctx context.Context, recoveryToken string) (string, error)
	}

	mediaService interface {
		UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*entity.Avatar, error)
	}

	eventProducer interface {
		Publish(ctx context.Context, topic string, event any) error
	}
)
