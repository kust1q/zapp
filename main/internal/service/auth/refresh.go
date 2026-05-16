package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/sirupsen/logrus"
)

func (s *service) Refresh(ctx context.Context, req *entity.Refresh) (*entity.Tokens, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req.Refresh = strings.TrimSpace(req.Refresh)

	userID, err := s.tokens.GetUserIdByRefreshToken(ctx, req.Refresh)
	if err != nil {
		if errors.Is(err, errs.ErrTokenNotFound) {
			return nil, errs.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if err = s.tokens.RemoveRefresh(ctx, req.Refresh); err != nil {
		logrus.Warnf("failed to delete refresh token: %v", err)
	}

	userIntID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errs.ErrInvalidRefreshToken
	}

	user, err := s.db.GetUserByID(ctx, userIntID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	role := "user"
	if user.IsSuperuser {
		role = "admin"
	}
	accessToken, err := s.generateAccessToken(user.ID, user.Credential.Email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, strconv.Itoa(user.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &entity.Tokens{
		Access: &entity.Access{
			Access: accessToken,
		},
		Refresh: &entity.Refresh{
			Refresh: refreshToken,
		},
	}, nil
}
