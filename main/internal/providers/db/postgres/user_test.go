package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
)

type mockCache struct {
	mock.Mock
}

func (m *mockCache) SetUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockCache) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockCache) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockCache) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockCache) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *mockCache) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *mockCache) InvalidateUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockCache) SetTweet(ctx context.Context, tweet *models.Tweet) error { return nil }
func (m *mockCache) GetTweet(ctx context.Context, tweetID int) (*models.Tweet, error) { return nil, nil }
func (m *mockCache) MGetTweets(ctx context.Context, tweetIDs []int) (map[int]*models.Tweet, error) { return nil, nil }
func (m *mockCache) InvalidateTweet(ctx context.Context, tweetID int) error { return nil }
func (m *mockCache) SetUserTweetIDs(ctx context.Context, username string, ids []int) error { return nil }
func (m *mockCache) GetUserTweetIDs(ctx context.Context, username string) ([]int, error) { return nil, nil }
func (m *mockCache) InvalidateUserTweets(ctx context.Context, username string) error { return nil }
func (m *mockCache) SetReplyIDs(ctx context.Context, parentTweetID int, ids []int) error { return nil }
func (m *mockCache) GetReplyIDs(ctx context.Context, parentTweetID int) ([]int, error) { return nil, nil }
func (m *mockCache) InvalidateReplies(ctx context.Context, parentTweetID int) error { return nil }
func (m *mockCache) SetTweetLikerIDs(ctx context.Context, tweetID int, userIDs []int) error { return nil }
func (m *mockCache) GetTweetLikerIDs(ctx context.Context, tweetID int) ([]int, error) { return nil, nil }
func (m *mockCache) InvalidateTweetLikers(ctx context.Context, tweetID int) error { return nil }
func (m *mockCache) SetTweetCounters(ctx context.Context, tweetID int, counters *models.Counters) error { return nil }
func (m *mockCache) GetTweetCounters(ctx context.Context, tweetID int) (*models.Counters, error) { return nil, nil }
func (m *mockCache) InvalidateTweetCounters(ctx context.Context, tweetID int) error { return nil }

func TestPostgres_GetUserByEmail_Success(t *testing.T) {
	db, smock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	
	cacheMock := new(mockCache)
	pg := NewPostgresDB(sqlxDB, cacheMock)
	ctx := context.Background()

	email := "test@test.com"

	cacheMock.On("GetUserByEmail", mock.Anything, email).Return(nil, errs.ErrCacheKeyNotFound)
	
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "gen", "created_at", "is_superuser"}).
		AddRow(1, "test", email, "pass", "bio", "male", time.Now(), false)
	smock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1").WithArgs(email).WillReturnRows(rows)

	cacheMock.On("SetUser", mock.Anything, mock.Anything).Return(nil)

	res, err := pg.GetUserByEmail(ctx, email)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 1, res.ID)
		assert.Equal(t, email, res.Credential.Email)
	}
}
