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

func (s *service) UpdatePassword(ctx context.Context, req *entity.UpdatePassword) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	user, err := s.db.GetUserByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Credential.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("invalid password")
	}

	if err = s.tokens.CloseAllSessions(ctx, strconv.Itoa(req.UserID)); err != nil {
		return fmt.Errorf("failed to close sessions: %w", err)
	}
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}
	return s.db.UpdateUserPassword(ctx, req.UserID, string(newHashPassword))
}
