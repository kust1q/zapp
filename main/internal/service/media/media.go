package media

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

type service struct {
	db     db
	object objectStorage
}

func NewMediaService(db db, object objectStorage) *service {
	return &service{db: db, object: object}
}

func (s *service) UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (string, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mt, err := s.detectMediaType(filename)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file for size calculation: %w", err)
	}

	bufReader := bytes.NewReader(buf.Bytes())

	path, mime, err := s.object.Upload(ctx, bufReader, mt, filename)
	if err != nil {
		return "", err
	}
	tweetMedia, err := s.db.UpsertByTweetIdTx(ctx, tx, &entity.TweetMedia{
		TweetID:   tweetID,
		Path:      path,
		MimeType:  mime,
		SizeBytes: size,
	})
	if err != nil {
		go s.asyncCleanup(path)
		return "", fmt.Errorf("upsert media failed: %w", err)
	}
	return s.object.GetPresignedURL(ctx, tweetMedia.Path)
}

func (s *service) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	mediaPath, err := s.db.GetMediaPathByTweetID(ctx, tweetID)
	if err != nil {
		return "", fmt.Errorf("failed to get media url: %w", err)
	}

	return s.object.GetPresignedURL(ctx, mediaPath)
}

func (s *service) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweetMedia, err := s.db.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get tweet media data: %w", err)
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	mediaURL, err := s.object.GetPresignedURL(ctx, tweetMedia.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet media url: %w", err)
	}

	return &entity.TweetMedia{
		ID:        tweetMedia.ID,
		TweetID:   tweetMedia.TweetID,
		Path:      mediaURL,
		MimeType:  tweetMedia.MimeType,
		SizeBytes: tweetMedia.SizeBytes,
	}, nil
}

func (s *service) DeleteTweetMedia(ctx context.Context, tweetID, userID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	media, err := s.db.GetMediaDataByTweetID(ctx, tweetID)
	if err != nil {
		return err
	}
	if err := s.db.DeleteMediaByTweetID(ctx, tweetID, userID); err != nil {
		return err
	}
	if media.Path != "" {
		go s.asyncCleanup(media.Path)
	}
	return nil
}

func (s *service) UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (res *entity.Avatar, err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var mt entity.MediaType
	mt, err = s.detectMediaType(filename)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for size calculation: %w", err)
	}

	bufReader := bytes.NewReader(buf.Bytes())

	path, mime, err := s.object.Upload(ctx, bufReader, mt, filename)
	if err != nil {
		return nil, err
	}

	avatar, err := s.db.UploadAvatarTx(ctx, tx, &entity.Avatar{
		UserID:    userID,
		Path:      path,
		MimeType:  mime,
		SizeBytes: size,
	})

	if err != nil {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go s.cleanUpMedia(bgCtx, path)
		return nil, fmt.Errorf("upload avatar failed: %w", err)
	}

	avatar.Path, err = s.object.GetPresignedURL(ctx, avatar.Path)
	if err != nil {
		return nil, fmt.Errorf("get avatar url failed: %w", err)
	}

	return avatar, nil
}

func (s *service) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	avatarPath, err := s.db.GetAvatarPathByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get avatar url: %w", err)
	}

	return s.object.GetPresignedURL(ctx, avatarPath)
}

func (s *service) GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	avatar, err := s.db.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	avatarPath, err := s.db.GetAvatarPathByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar path: %w", err)
	}

	avatarUrl, err := s.object.GetPresignedURL(ctx, avatarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar url: %w", err)
	}

	avatar.Path = avatarUrl

	return avatar, nil
}

func (s *service) DeleteAvatar(ctx context.Context, userID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	avatar, err := s.db.GetAvatarDataByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get data of useravatar: %w", err)
	}

	if err := s.db.DeleteAvatarByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete avatar by user id: %w", err)
	}

	if avatar.Path != "" {
		go s.asyncCleanup(avatar.Path)
	}
	return nil
}

func (s *service) DeleteMediasByUserID(ctx context.Context, userID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	urls, err := s.db.GetMediaUrlsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get medias urls: %w", err)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			s.cleanUpMedia(ctx, url)
		}(url)
	}
	wg.Wait()
	return nil
}

func (s *service) cleanUpMedia(ctx context.Context, Path string) {
	if err := s.object.Remove(ctx, Path); err != nil {
		logrus.WithField("media_url", Path).Warnf("failed to remove media: %s", err)
	}
}

func (s *service) asyncCleanup(path string) {
	if path == "" {
		return
	}
	go func(p string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.cleanUpMedia(bgCtx, p)
	}(path)
}

func (s *service) GetPresignedURL(ctx context.Context, path string) (string, error) {
	return s.object.GetPresignedURL(ctx, path)
}

func (s *service) detectMediaType(filename string) (entity.MediaType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return entity.MediaTypeImage, nil
	case ".gif":
		return entity.MediaTypeGIF, nil
	case ".mp4", ".mov", ".m4v":
		return entity.MediaTypeVideo, nil
	case ".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a", ".webm":
		return entity.MediaTypeAudio, nil
	default:
		return "", fmt.Errorf("invalid media type")
	}
}
