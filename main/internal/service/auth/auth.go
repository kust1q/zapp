package auth

import (
	"time"

	"github.com/kust1q/Zapp/main/internal/config"
)

type service struct {
	cfg      *config.AuthServiceConfig
	db       db
	tokens   tokenStorage
	media    mediaService
	producer eventProducer
}

func NewAuthService(cfg *config.AuthServiceConfig, db db, media mediaService, tokens tokenStorage, producer eventProducer) *service {
	return &service{
		cfg:      cfg,
		db:       db,
		media:    media,
		tokens:   tokens,
		producer: producer,
	}
}

func (s *service) GetRefreshTTL() time.Duration {
	return s.cfg.RefreshTTL
}
