package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) RecoveryPassword(ctx context.Context, req *entity.Recovery) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.RecoveryToken = strings.TrimSpace(req.RecoveryToken)

	userIDstr, err := s.tokens.GetUserIdByRecoveryToken(ctx, req.RecoveryToken)
	if err != nil {
		return fmt.Errorf("failed to get userID: %w", err)
	}

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		return fmt.Errorf("failed to get userID: %w", err)
	}

	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	return s.db.UpdateUserPassword(ctx, userID, string(newHashPassword))
}
