package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kust1q/Zapp/main/internal/config"
	"github.com/kust1q/Zapp/main/internal/service/auth"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	tx, _ := args.Get(0).(*sql.Tx)
	return tx, args.Error(1)
}

func (m *mockDB) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockDB) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockDB) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *mockDB) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockDB) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *mockDB) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type mockTokenStorage struct {
	mock.Mock
}

func (m *mockTokenStorage) StoreRefresh(ctx context.Context, refreshToken, userID string, ttl time.Duration) error {
	args := m.Called(ctx, refreshToken, userID, ttl)
	return args.Error(0)
}

func (m *mockTokenStorage) GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *mockTokenStorage) CloseAllSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockTokenStorage) RemoveRefresh(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *mockTokenStorage) StoreRecovery(ctx context.Context, token, userID string, ttl time.Duration) error {
	args := m.Called(ctx, token, userID, ttl)
	return args.Error(0)
}

func (m *mockTokenStorage) GetUserIdByRecoveryToken(ctx context.Context, recoveryToken string) (string, error) {
	args := m.Called(ctx, recoveryToken)
	return args.String(0), args.Error(1)
}

type mockMediaService struct {
	mock.Mock
}

func (m *mockMediaService) UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*entity.Avatar, error) {
	args := m.Called(ctx, userID, file, filename, tx)
	avatar := args.Get(0)
	if avatar == nil {
		return nil, args.Error(1)
	}
	return avatar.(*entity.Avatar), args.Error(1)
}

type mockEventProducer struct {
	mock.Mock
}

func (m *mockEventProducer) Publish(ctx context.Context, topic string, event any) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}

func generateTestRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

func TestService_GetRefreshTTL(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		RefreshTTL: 24 * time.Hour,
	}

	service := auth.NewAuthService(
		cfg,
		&mockDB{},
		&mockMediaService{},
		&mockTokenStorage{},
		&mockEventProducer{},
	)

	ttl := service.GetRefreshTTL()
	assert.Equal(t, 24*time.Hour, ttl)
}

func TestService_SignIn_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  time.Hour,
		RefreshTTL: 24 * time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Credential: &entity.Credential{
			Email:    "test@example.com",
			Password: string(hashedPassword),
		},
		IsSuperuser: false,
	}

	mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil).Once()
	mockTokens.On("StoreRefresh", mock.Anything, mock.Anything, "1", cfg.RefreshTTL).Return(nil).Once()

	req := &entity.Credential{
		Email:    "test@example.com",
		Password: password,
	}

	tokens, err := service.SignIn(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.Access.Access)
	assert.NotEmpty(t, tokens.Refresh.Refresh)

	mockDB.AssertExpectations(t)
	mockTokens.AssertExpectations(t)
}

func TestService_SignIn_InvalidCredentials(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("differentpassword"), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Credential: &entity.Credential{
			Email:    "test@example.com",
			Password: string(hashedPassword),
		},
	}

	mockDB.On("GetUserByEmail", mock.Anything, "test@example.com").Return(user, nil).Once()

	req := &entity.Credential{
		Email:    "test@example.com",
		Password: password,
	}

	tokens, err := service.SignIn(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Equal(t, errs.ErrInvalidCredentials, err)
}

func TestService_SignUp_EmailAlreadyExists(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()

	req := &entity.User{
		Username: "newuser",
		Credential: &entity.Credential{
			Email:    "existing@example.com",
			Password: "password123",
		},
	}

	existingUser := &entity.User{
		ID:       1,
		Username: "existinguser",
		Credential: &entity.Credential{
			Email: "existing@example.com",
		},
	}

	mockDB.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil).Once()
	mockDB.On("GetUserByUsername", mock.Anything, "newuser").Return(nil, errs.ErrUserNotFound).Once()

	result, err := service.SignUp(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, errs.ErrEmailAlreadyUsed) ||
		(err != nil && err.Error() == "failed to begin transaction: email already used"))
}

func TestService_Refresh_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  time.Hour,
		RefreshTTL: 24 * time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	refreshToken := "valid-refresh-token"

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Credential: &entity.Credential{
			Email: "test@example.com",
		},
		IsSuperuser: false,
	}

	mockTokens.On("GetUserIdByRefreshToken", mock.Anything, refreshToken).Return("1", nil).Once()
	mockTokens.On("RemoveRefresh", mock.Anything, refreshToken).Return(nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockTokens.On("StoreRefresh", mock.Anything, mock.Anything, "1", cfg.RefreshTTL).Return(nil).Once()

	req := &entity.Refresh{Refresh: refreshToken}
	tokens, err := service.Refresh(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.Access.Access)
	assert.NotEmpty(t, tokens.Refresh.Refresh)

	userID, err := service.VerifyAccessToken(tokens.Access.Access)
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

func TestService_Refresh_InvalidToken(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	refreshToken := "invalid-refresh-token"

	mockTokens.On("GetUserIdByRefreshToken", mock.Anything, refreshToken).Return("", errs.ErrTokenNotFound).Once()

	req := &entity.Refresh{Refresh: refreshToken}
	tokens, err := service.Refresh(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Equal(t, errs.ErrInvalidRefreshToken, err)
}

func TestService_SignOut_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	refreshToken := "refresh-token-to-delete"

	mockTokens.On("RemoveRefresh", mock.Anything, refreshToken).Return(nil).Once()

	req := &entity.Refresh{Refresh: refreshToken}
	err := service.SignOut(ctx, req)

	assert.NoError(t, err)
	mockTokens.AssertExpectations(t)
}

func TestService_UpdatePassword_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	oldPassword := "oldpassword123"
	newPassword := "newpassword123"

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	user := &entity.User{
		ID: 1,
		Credential: &entity.Credential{
			Password: string(hashedOldPassword),
		},
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockTokens.On("CloseAllSessions", mock.Anything, "1").Return(nil).Once()
	mockDB.On("UpdateUserPassword", mock.Anything, 1, mock.AnythingOfType("string")).Return(nil).Once()

	req := &entity.UpdatePassword{
		UserID:      1,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	err := service.UpdatePassword(ctx, req)
	assert.NoError(t, err)
}

func TestService_UpdatePassword_InvalidOldPassword(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	oldPassword := "oldpassword123"
	wrongPassword := "wrongpassword"
	newPassword := "newpassword123"

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	user := &entity.User{
		ID: 1,
		Credential: &entity.Credential{
			Password: string(hashedOldPassword),
		},
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()

	req := &entity.UpdatePassword{
		UserID:      1,
		OldPassword: wrongPassword,
		NewPassword: newPassword,
	}

	err := service.UpdatePassword(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestService_ForgotPassword_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		RecoveryTTL: time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	email := "test@example.com"

	user := &entity.User{
		ID: 1,
		Credential: &entity.Credential{
			Email: email,
		},
	}

	mockDB.On("GetUserByEmail", mock.Anything, email).Return(user, nil).Once()
	mockTokens.On("CloseAllSessions", mock.Anything, "1").Return(nil).Once()
	mockTokens.On("StoreRecovery", mock.Anything, mock.AnythingOfType("string"), "1", cfg.RecoveryTTL).Return(nil).Once()

	req := &entity.ForgotPassword{Email: email}
	recovery, err := service.ForgotPassword(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, recovery)
	assert.NotEmpty(t, recovery.Recovery)

	_, err = uuid.Parse(recovery.Recovery)
	assert.NoError(t, err)
}

func TestService_ForgotPassword_UserNotFound(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		RecoveryTTL: time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	email := "nonexistent@example.com"

	mockDB.On("GetUserByEmail", mock.Anything, email).Return(nil, errs.ErrUserNotFound).Once()

	req := &entity.ForgotPassword{Email: email}

	recovery, err := service.ForgotPassword(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, recovery)
	assert.Equal(t, errs.ErrUserNotFound, err)
}

func TestService_RecoveryPassword_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	recoveryToken := "valid-recovery-token"
	newPassword := "newpassword123"

	mockTokens.On("GetUserIdByRecoveryToken", mock.Anything, recoveryToken).Return("1", nil).Once()
	mockDB.On("UpdateUserPassword", mock.Anything, 1, mock.AnythingOfType("string")).Return(nil).Once()

	req := &entity.RecoveryPassword{
		RecoveryToken: recoveryToken,
		NewPassword:   newPassword,
	}

	err := service.RecoveryPassword(ctx, req)
	assert.NoError(t, err)
}

func TestService_RecoveryPassword_InvalidToken(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	ctx := context.Background()
	recoveryToken := "invalid-recovery-token"
	newPassword := "newpassword123"

	mockTokens.On("GetUserIdByRecoveryToken", mock.Anything, recoveryToken).Return("", errors.New("token not found")).Once()

	req := &entity.RecoveryPassword{
		RecoveryToken: recoveryToken,
		NewPassword:   newPassword,
	}

	err := service.RecoveryPassword(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get userID")
}

func TestService_VerifyAccessToken_Success(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	claims := auth.AccessClaims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	userID, err := service.VerifyAccessToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

func TestService_VerifyAccessToken_Expired(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	claims := auth.AccessClaims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	userID, err := service.VerifyAccessToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}

func TestService_VerifyAccessToken_InvalidSignature(t *testing.T) {
	privateKey, publicKey := generateTestRSAKeys(t)
	cfg := &config.AuthServiceConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		AccessTTL:  time.Hour,
	}

	mockDB := &mockDB{}
	mockTokens := &mockTokenStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := auth.NewAuthService(cfg, mockDB, mockMedia, mockTokens, mockProducer)

	wrongPrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	claims := auth.AccessClaims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(wrongPrivateKey)
	require.NoError(t, err)

	userID, err := service.VerifyAccessToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "token validation failed")
}
