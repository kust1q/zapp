package media

import (
	"context"
	"database/sql"
	"io"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	db interface {
		UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *entity.TweetMedia) (*entity.TweetMedia, error)
		GetMediaPathByTweetID(ctx context.Context, tweetID int) (string, error)
		GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
		DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error

		UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error)
		GetAvatarPathByUserID(ctx context.Context, userID int) (string, error)
		GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
		DeleteAvatarByUserID(ctx context.Context, userID int) error

		GetMediaUrlsByUserID(ctx context.Context, userID int) ([]string, error)
	}

	objectStorage interface {
		Upload(ctx context.Context, file io.Reader, mediaType entity.MediaType, filename string) (path string, mimeType string, err error)
		Remove(ctx context.Context, objectPath string) error
		GetPresignedURL(ctx context.Context, objectPath string) (string, error)
	}
)
