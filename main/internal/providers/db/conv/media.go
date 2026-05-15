package conv

import (
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
)

func FromDomainToTweetMediaModel(media *entity.TweetMedia) *models.TweetMedia {
	if media == nil {
		return nil
	}

	return &models.TweetMedia{
		ID:        media.ID,
		TweetID:   media.TweetID,
		Path:      media.Path,
		MimeType:  media.MimeType,
		SizeBytes: media.SizeBytes,
	}
}

func FromTweetMediaModelToDomain(media *models.TweetMedia) *entity.TweetMedia {
	if media == nil {
		return nil
	}

	return &entity.TweetMedia{
		ID:        media.ID,
		TweetID:   media.TweetID,
		Path:      media.Path,
		MimeType:  media.MimeType,
		SizeBytes: media.SizeBytes,
	}
}

func FromDomainToAvatarModel(avatar *entity.Avatar) *models.Avatar {
	if avatar == nil {
		return nil
	}

	return &models.Avatar{
		ID:        avatar.ID,
		UserID:    avatar.UserID,
		Path:      avatar.Path,
		MimeType:  avatar.MimeType,
		SizeBytes: avatar.SizeBytes,
	}
}

func FromAvatarModelToDomain(avatar *models.Avatar) *entity.Avatar {
	if avatar == nil {
		return nil
	}

	return &entity.Avatar{
		ID:        avatar.ID,
		UserID:    avatar.UserID,
		Path:      avatar.Path,
		MimeType:  avatar.MimeType,
		SizeBytes: avatar.SizeBytes,
	}
}
