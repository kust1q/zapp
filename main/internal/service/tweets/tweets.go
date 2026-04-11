package tweets

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type service struct {
	db       tweetStorage
	media    mediaService
	producer eventProducer
}

func NewTweetService(db tweetStorage, media mediaService, eventProducer eventProducer) *service {
	return &service{
		db:       db,
		media:    media,
		producer: eventProducer,
	}
}

func (s *service) BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	author, err := s.db.GetUserByID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet author: %w", err)
	}

	avatarUrl, err := s.media.GetAvatarUrlByUserID(ctx, tweet.Author.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user avatar: %w", err)
	}

	mediaUrl, err := s.media.GetMediaUrlByTweetID(ctx, tweet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweet media url: %w", err)
	}
	tweet.MediaUrl = mediaUrl

	counts, err := s.db.GetCounts(ctx, tweet.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get counters: %w", err)
	}

	return &entity.Tweet{
		ID:            tweet.ID,
		ParentTweetID: tweet.ParentTweetID,
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt,
		UpdatedAt:     tweet.UpdatedAt,
		MediaUrl:      tweet.MediaUrl,
		Author: &entity.SmallUser{
			ID:        author.ID,
			Username:  author.Username,
			AvatarUrl: avatarUrl,
		},
		Counters: counts,
	}, nil
}
