package tweets_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/kust1q/Zapp/main/internal/service/tweets"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTweetStorage struct {
	mock.Mock
}

func (m *mockTweetStorage) BeginTx(ctx context.Context) (*sql.Tx, error) {
	args := m.Called(ctx)
	tx, _ := args.Get(0).(*sql.Tx)
	return tx, args.Error(1)
}

func (m *mockTweetStorage) CreateTweetTx(ctx context.Context, tx *sql.Tx, tweet *entity.Tweet) (*entity.Tweet, error) {
	args := m.Called(ctx, tx, tweet)
	return args.Get(0).(*entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	args := m.Called(ctx, tweet)
	return args.Get(0).(*entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	args := m.Called(ctx, tweetID)
	tweet := args.Get(0)
	if tweet == nil {
		return nil, args.Error(1)
	}
	return tweet.(*entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) UpdateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	args := m.Called(ctx, tweet)
	return args.Get(0).(*entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	args := m.Called(ctx, userID, tweetID)
	return args.Error(0)
}

func (m *mockTweetStorage) LikeTweet(ctx context.Context, userID, tweetID int) error {
	args := m.Called(ctx, userID, tweetID)
	return args.Error(0)
}

func (m *mockTweetStorage) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	args := m.Called(ctx, userID, tweetID)
	return args.Error(0)
}

func (m *mockTweetStorage) Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error {
	args := m.Called(ctx, userID, tweetID, createdAt)
	return args.Error(0)
}

func (m *mockTweetStorage) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	args := m.Called(ctx, userID, retweetID)
	return args.Error(0)
}

func (m *mockTweetStorage) GetRepliesToTweet(ctx context.Context, parentTweetID, limit, offset int) ([]entity.Tweet, error) {
	args := m.Called(ctx, parentTweetID, limit, offset)
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error) {
	args := m.Called(ctx, username, limit, offset)
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

func (m *mockTweetStorage) GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error) {
	args := m.Called(ctx, tweetID)
	counters := args.Get(0)
	if counters == nil {
		return nil, args.Error(1)
	}
	return counters.(*entity.Counters), args.Error(1)
}

func (m *mockTweetStorage) GetLikes(ctx context.Context, tweetID int, limit, offset int) ([]entity.SmallUser, error) {
	args := m.Called(ctx, tweetID, limit, offset)
	return args.Get(0).([]entity.SmallUser), args.Error(1)
}

func (m *mockTweetStorage) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

type mockMediaService struct {
	mock.Mock
}

func (m *mockMediaService) UploadAndAttachTweetMediaTx(ctx context.Context, tweetID int, file io.Reader, filename string, tx *sql.Tx) (string, error) {
	args := m.Called(ctx, tweetID, file, filename, tx)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	args := m.Called(ctx, tweetID)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) DeleteTweetMedia(ctx context.Context, tweetID, userID int) error {
	args := m.Called(ctx, tweetID, userID)
	return args.Error(0)
}

func (m *mockMediaService) GetPresignedURL(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

type mockEventProducer struct {
	mock.Mock
}

func (m *mockEventProducer) Publish(ctx context.Context, topic string, event any) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}

func TestService_GetTweetById_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	tweet := &entity.Tweet{
		ID:      1,
		Content: "Test tweet",
		Author: &entity.SmallUser{
			ID: 1,
		},
	}

	author := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	counters := &entity.Counters{
		ReplyCount:   0,
		RetweetCount: 0,
		LikeCount:    0,
	}

	mockDB.On("GetTweetById", mock.Anything, 1).Return(tweet, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 1).Return(author, nil).Once()
	mockDB.On("GetCounts", mock.Anything, 1).Return(counters, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, 1).Return("", nil).Once()

	result, err := service.GetTweetById(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test tweet", result.Content)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetTweetById_NotFound(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("GetTweetById", mock.Anything, 1).Return(nil, errs.ErrTweetNotFound).Once()

	result, err := service.GetTweetById(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errs.ErrTweetNotFound, err)

	mockDB.AssertExpectations(t)
}

func TestService_DeleteTweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteTweetMedia", mock.Anything, 1, 1).Return(nil).Once()
	mockDB.On("DeleteTweet", mock.Anything, 1, 1).Return(nil).Once()
	mockProducer.On("Publish", mock.Anything, events.TopicTweet, mock.AnythingOfType("events.TweetDeleted")).Return(nil).Once()

	err := service.DeleteTweet(ctx, 1, 1)

	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteTweet_MediaNotFound(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteTweetMedia", mock.Anything, 1, 1).Return(errs.ErrTweetMediaNotFound).Once()
	mockDB.On("DeleteTweet", mock.Anything, 1, 1).Return(nil).Once()
	mockProducer.On("Publish", mock.Anything, events.TopicTweet, mock.AnythingOfType("events.TweetDeleted")).Return(nil).Once()

	err := service.DeleteTweet(ctx, 1, 1)

	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteTweet_OtherMediaError(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteTweetMedia", mock.Anything, 1, 1).Return(errors.New("other media error")).Once()
	mockDB.On("DeleteTweet", mock.Anything, 1, 1).Return(nil).Once()
	mockProducer.On("Publish", mock.Anything, events.TopicTweet, mock.AnythingOfType("events.TweetDeleted")).Return(nil).Once()

	err := service.DeleteTweet(ctx, 1, 1)

	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteTweet_DBError(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteTweetMedia", mock.Anything, 1, 1).Return(nil).Once()
	mockDB.On("DeleteTweet", mock.Anything, 1, 1).Return(errors.New("database error")).Once()

	err := service.DeleteTweet(ctx, 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_LikeTweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("LikeTweet", mock.Anything, 1, 1).Return(nil).Once()

	err := service.LikeTweet(ctx, 1, 1)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_LikeTweet_TweetNotFound(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("LikeTweet", mock.Anything, 1, 1).Return(errs.ErrTweetNotFound).Once()

	err := service.LikeTweet(ctx, 1, 1)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrTweetNotFound, err)

	mockDB.AssertExpectations(t)
}

func TestService_UnlikeTweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("UnLikeTweet", mock.Anything, 1, 1).Return(nil).Once()

	err := service.UnlikeTweet(ctx, 1, 1)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_GetTweetsAndRetweetsByUsername_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	tweets := []entity.Tweet{
		{
			ID:      1,
			Content: "Tweet 1",
			Author: &entity.SmallUser{
				ID: 1,
			},
		},
		{
			ID:      2,
			Content: "Tweet 2",
			Author: &entity.SmallUser{
				ID: 1,
			},
		},
	}

	author := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	counters := &entity.Counters{
		ReplyCount:   0,
		RetweetCount: 0,
		LikeCount:    0,
	}

	mockDB.On("GetTweetsAndRetweetsByUsername", mock.Anything, "testuser", 10, 0).Return(tweets, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 1).Return(author, nil).Times(2)
	mockDB.On("GetCounts", mock.Anything, mock.AnythingOfType("int")).Return(counters, nil).Times(2)
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Times(2)
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, mock.AnythingOfType("int")).Return("", nil).Times(2)

	result, err := service.GetTweetsAndRetweetsByUsername(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Tweet 1", result[0].Content)
	assert.Equal(t, "Tweet 2", result[1].Content)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetTweetsAndRetweetsByUsername_NoRows(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("GetTweetsAndRetweetsByUsername", mock.Anything, "testuser", 10, 0).Return([]entity.Tweet{}, sql.ErrNoRows).Once()

	result, err := service.GetTweetsAndRetweetsByUsername(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)

	mockDB.AssertExpectations(t)
}

func TestService_GetRepliesToTweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	replies := []entity.Tweet{
		{
			ID:      2,
			Content: "Reply 1",
			Author: &entity.SmallUser{
				ID: 2,
			},
		},
	}

	author := &entity.User{
		ID:       2,
		Username: "replyuser",
	}

	counters := &entity.Counters{
		ReplyCount:   0,
		RetweetCount: 0,
		LikeCount:    0,
	}

	mockDB.On("GetRepliesToTweet", mock.Anything, 1, 10, 0).Return(replies, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 2).Return(author, nil).Once()
	mockDB.On("GetCounts", mock.Anything, 2).Return(counters, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 2).Return("/avatars/2.jpg", nil).Once()
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, 2).Return("", nil).Once()

	result, err := service.GetRepliesToTweet(ctx, 1, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Reply 1", result[0].Content)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetLikes_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	users := []entity.SmallUser{
		{
			ID:       1,
			Username: "user1",
		},
		{
			ID:       2,
			Username: "user2",
		},
	}

	mockDB.On("GetLikes", mock.Anything, 1, 10, 0).Return(users, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 2).Return("/avatars/2.jpg", nil).Once()

	result, err := service.GetLikes(ctx, 1, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1", result[0].Username)
	assert.Equal(t, "/avatars/1.jpg", result[0].AvatarUrl)
	assert.Equal(t, "user2", result[1].Username)
	assert.Equal(t, "/avatars/2.jpg", result[1].AvatarUrl)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_CreateRetweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("Retweet", mock.Anything, 1, 1, mock.AnythingOfType("time.Time")).Return(nil).Once()

	err := service.CreateRetweet(ctx, 1, 1)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_DeleteRetweet_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("DeleteRetweet", mock.Anything, 1, 1).Return(nil).Once()

	err := service.DeleteRetweet(ctx, 1, 1)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_UpdateTweet_NotFound(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	req := &entity.Tweet{
		ID:      1,
		Content: "Updated content",
		Author: &entity.SmallUser{
			ID: 1,
		},
	}

	mockDB.On("GetTweetById", mock.Anything, 1).Return(nil, sql.ErrNoRows).Once()

	result, err := service.UpdateTweet(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errs.ErrTweetNotFound, err)

	mockDB.AssertExpectations(t)
}

func TestService_UpdateTweet_Unauthorized(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	existingTweet := &entity.Tweet{
		ID:      1,
		Content: "Old content",
		Author: &entity.SmallUser{
			ID: 2,
		},
	}

	req := &entity.Tweet{
		ID:      1,
		Content: "Updated content",
		Author: &entity.SmallUser{
			ID: 1,
		},
	}

	mockDB.On("GetTweetById", mock.Anything, 1).Return(existingTweet, nil).Once()

	result, err := service.UpdateTweet(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errs.ErrUnauthorizedUpdate, err)

	mockDB.AssertExpectations(t)
}

func TestService_BuildEntityTweetToResponse_Success(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	tweet := &entity.Tweet{
		ID:      1,
		Content: "Test tweet",
		Author: &entity.SmallUser{
			ID: 1,
		},
	}

	author := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	counters := &entity.Counters{
		ReplyCount:   0,
		RetweetCount: 0,
		LikeCount:    0,
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(author, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, 1).Return("/media/test.jpg", nil).Once()
	mockDB.On("GetCounts", mock.Anything, 1).Return(counters, nil).Once()

	result, err := service.BuildEntityTweetToResponse(ctx, tweet)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Author.Username)
	assert.Equal(t, "/avatars/1.jpg", result.Author.AvatarUrl)
	assert.Equal(t, "/media/test.jpg", result.MediaUrl)
	assert.Equal(t, 0, result.Counters.LikeCount)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_BuildEntityTweetToResponse_UserNotFound(t *testing.T) {
	mockDB := &mockTweetStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := tweets.NewTweetService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	tweet := &entity.Tweet{
		ID:      1,
		Content: "Test tweet",
		Author: &entity.SmallUser{
			ID: 1,
		},
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(nil, errors.New("user not found")).Once()

	result, err := service.BuildEntityTweetToResponse(ctx, tweet)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get tweet author")

	mockDB.AssertExpectations(t)
}
