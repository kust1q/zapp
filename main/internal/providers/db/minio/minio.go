package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/main/internal/config"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/minio/minio-go/v7"
)

type minioDB struct {
	client        *minio.Client
	config        *config.MinioConfig
	mediaPolicies map[entity.MediaType]entity.MediaPolicy
}

func NewMinioDB(client *minio.Client, config *config.MinioConfig, mediaPolicies map[entity.MediaType]entity.MediaPolicy) *minioDB {
	return &minioDB{
		client:        client,
		config:        config,
		mediaPolicies: mediaPolicies,
	}
}

func (s *minioDB) Upload(ctx context.Context, file io.Reader, mediaType entity.MediaType, filename string) (path string, mimeType string, err error) {
	config, ok := s.mediaPolicies[mediaType]
	if !ok {
		return "", "", fmt.Errorf("unsupported media type: %s", mediaType)
	}
	ext := strings.ToLower(filepath.Ext(filename))
	data, err := s.readAndValidate(file, ext, config)
	if err != nil {
		return "", "", err
	}

	if config.ForceMimeType != "" {
		mimeType = config.ForceMimeType
	} else {
		mimeType = http.DetectContentType(data)
	}

	path = filepath.Join(string(mediaType), uuid.New().String(), filename)

	_, err = s.client.PutObject(
		ctx,
		s.config.BucketName,
		path,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: mimeType,
		})

	if err != nil {
		return "", "", err
	}

	return path, mimeType, nil
}

func (s *minioDB) Remove(ctx context.Context, objectPath string) error {
	return s.client.RemoveObject(
		ctx,
		s.config.BucketName,
		objectPath,
		minio.RemoveObjectOptions{},
	)
}

func (s *minioDB) GetPresignedURL(ctx context.Context, objectPath string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.config.BucketName, objectPath, s.config.TTL, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *minioDB) readAndValidate(reader io.Reader, ext string, policy entity.MediaPolicy) ([]byte, error) {
	limitedReader := io.LimitReader(reader, policy.MaxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > policy.MaxSize {
		return nil, errs.ErrFileTooLarge
	}

	detectedType := http.DetectContentType(data)
	if !s.isAllowedMimeType(detectedType, policy) {
		return nil, fmt.Errorf("invalid media type: %s", detectedType)
	}
	if !s.isValidExtension(ext, policy) {
		return nil, fmt.Errorf("invalid extension: %s", ext)
	}

	return data, nil
}

func (s *minioDB) isValidExtension(ext string, policy entity.MediaPolicy) bool {
	if len(policy.AllowedExt) == 0 {
		return true
	}

	for _, allowed := range policy.AllowedExt {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

func (s *minioDB) isAllowedMimeType(mimeType string, policy entity.MediaPolicy) bool {
	for _, allowed := range policy.AllowedMime {
		if mimeType == allowed {
			return true
		}
	}
	return false
}
