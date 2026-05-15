package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
)

func (s *service) ForgotPassword(ctx context.Context, req *entity.ForgotPassword) (*entity.Recovery, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to find user: %w", err)
	} else if errors.Is(err, errs.ErrUserNotFound) {
		return nil, err
	}

	if err = s.tokens.CloseAllSessions(ctx, strconv.Itoa(user.ID)); err != nil {
		return nil, fmt.Errorf("failed to close sessions: %w", err)
	}

	recoveryToken := s.generateRecoveryToken()
	if err = s.tokens.StoreRecovery(ctx, recoveryToken, strconv.Itoa(user.ID), s.cfg.RecoveryTTL); err != nil {
		return nil, fmt.Errorf("failed to store recovery token: %w", err)
	}

	return &entity.Recovery{
		Recovery: recoveryToken,
	}, nil
}

func (s *service) generateRecoveryToken() string {
	newUUID := uuid.New()
	token := newUUID.String()
	return token
}
