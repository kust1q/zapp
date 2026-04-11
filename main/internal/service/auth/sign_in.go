package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type AccessClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *service) SignIn(ctx context.Context, req *entity.Credential) (*entity.Tokens, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Credential.Password), []byte(req.Password)); err != nil {
		return nil, errs.ErrInvalidCredentials
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

func (s *service) generateAccessToken(userID int, email, role string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "zapp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.cfg.PrivateKey)
}

func (s *service) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(tokenBytes)
	if err := s.tokens.StoreRefresh(ctx, refreshToken, userID, s.cfg.RefreshTTL); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s *service) VerifyAccessToken(tokenString string) (int, error) {
	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.cfg.PublicKey, nil
		},
		jwt.WithLeeway(10*time.Second),
	)

	if err != nil {
		return 0, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
