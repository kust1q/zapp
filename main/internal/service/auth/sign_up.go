package auth

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/o1egl/govatar"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(ctx context.Context, req *entity.User) (user *entity.User, err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req.Credential.Email = strings.ToLower(strings.TrimSpace(req.Credential.Email))
	req.Username = strings.TrimSpace(req.Username)

	if err = s.checkUserExists(ctx, req.Credential.Email, req.Username); err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Credential.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	userEntity := entity.User{
		Username:    req.Username,
		Bio:         req.Bio,
		Gen:         req.Gen,
		CreatedAt:   req.CreatedAt,
		IsSuperuser: req.IsSuperuser,
		Credential: &entity.Credential{
			Email:    req.Credential.Email,
			Password: string(hashedPassword),
		},
	}

	var createdUser *entity.User
	createdUser, err = s.db.CreateUserTx(ctx, tx, &userEntity)
	if err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	var avatar *entity.Avatar
	avatar, err = s.generateAndUploadAvatar(ctx, createdUser.ID, createdUser.Username, createdUser.Gen, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate or upload avatar: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction failed: %w", err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event := events.UserEvent{
			EventType: events.UserCreateEvent,
			ID:        createdUser.ID,
			Username:  createdUser.Username,
			Bio:       createdUser.Bio,
		}

		if pErr := s.producer.Publish(bgCtx, events.TopicUser, event); pErr != nil {
			logrus.WithError(pErr).Error("failed to publish user.created")
		}
	}()

	return &entity.User{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Bio:       createdUser.Bio,
		Gen:       createdUser.Gen,
		AvatarUrl: avatar.Path,
		CreatedAt: createdUser.CreatedAt,
		Credential: &entity.Credential{
			Email: createdUser.Credential.Email,
		},
	}, nil
}

func (s *service) generateAndUploadAvatar(ctx context.Context, userID int, username, gender string, tx *sql.Tx) (*entity.Avatar, error) {
	genMap := map[string]govatar.Gender{
		"male":   govatar.MALE,
		"female": govatar.FEMALE,
	}
	gen, ok := genMap[strings.ToLower(gender)]
	if !ok {
		return nil, fmt.Errorf("invalid gender")
	}

	var img image.Image
	var err error
	img, err = govatar.GenerateForUsername(gen, username)
	if err != nil {
		return nil, fmt.Errorf("avatar generation failed: %w", err)
	}

	var buf bytes.Buffer
	if err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("JPEG encoding failed: %w", err)
	}
	avatarSaveName := uuid.New().String() + ".jpg"

	return s.media.UploadAvatarTx(ctx, userID, &buf, avatarSaveName, tx)
}

func (s *service) checkUserExists(ctx context.Context, email, username string) error {
	var err error
	_, err = s.db.GetUserByEmail(ctx, email)
	if err == nil {
		return errs.ErrEmailAlreadyUsed
	}
	if !errors.Is(err, errs.ErrUserNotFound) {
		return fmt.Errorf("email check failed: %w", err)
	}

	_, err = s.db.GetUserByUsername(ctx, username)
	if err == nil {
		return errs.ErrUsernameAlreadyUsed
	}
	if !errors.Is(err, errs.ErrUserNotFound) {
		return fmt.Errorf("username check failed: %w", err)
	}
	return nil
}
