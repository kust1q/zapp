package feed_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/kust1q/Zapp/main/internal/service/feed"
)

type mockFeedStorage struct {
	mock.Mock
}

func (m *mockFeedStorage) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockFeedStorage) GetFollowingsIds(ctx context.Context, username string, limit, offset int) ([]int, error) {
	args := m.Called(ctx, username, limit, offset)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockFeedStorage) GetFeedByAuthorsIds(ctx context.Context, userIDs []int, limit, offset int) ([]entity.Tweet, error) {
	args := m.Called(ctx, userIDs, limit, offset)
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

func (m *mockFeedStorage) GetAllTweets(ctx context.Context, limit, offset int) ([]entity.Tweet, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

type mockTweetService struct {
	mock.Mock
}

func (m *mockTweetService) BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	args := m.Called(ctx, tweet)
	val := args.Get(0)
	if val == nil {
		return nil, args.Error(1)
	}
	return val.(*entity.Tweet), args.Error(1)
}

func TestService_GetUserFeedByUserId_Success(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	followings := []int{2, 3}
	feedTweets := []entity.Tweet{
		{ID: 1, Content: "t1"},
	}
	defaultTweets := []entity.Tweet{
		{ID: 2, Content: "t2"},
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockDB.On("GetFollowingsIds", mock.Anything, "testuser", 10, 0).Return(followings, nil).Once()
	mockDB.On("GetFeedByAuthorsIds", mock.Anything, followings, 10, 0).Return(feedTweets, nil).Once()
	mockDB.On("GetAllTweets", mock.Anything, 10, 0).Return(defaultTweets, nil).Once()

	mockTweet.On("BuildEntityTweetToResponse", mock.Anything, mock.AnythingOfType("*entity.Tweet")).Return(&entity.Tweet{}, nil).Twice()

	res, err := service.GetUserFeedByUserId(ctx, 1, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}

func TestService_GetUserFeedByUserId_FindUserError(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()

	mockDB.On("GetUserByID", mock.Anything, 1).Return(nil, errs.ErrUserNotFound).Once()

	res, err := service.GetUserFeedByUserId(ctx, 1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}

func TestService_GetUserFeedByUserId_GetFollowingsError(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()
	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockDB.On("GetFollowingsIds", mock.Anything, "testuser", 10, 0).Return([]int{}, errors.New("db error")).Once()

	res, err := service.GetUserFeedByUserId(ctx, 1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}

func TestService_GetUserFeedByUserId_GetFeedError(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()
	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	followings := []int{2, 3}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockDB.On("GetFollowingsIds", mock.Anything, "testuser", 10, 0).Return(followings, nil).Once()
	mockDB.On("GetFeedByAuthorsIds", mock.Anything, followings, 10, 0).Return([]entity.Tweet{}, errors.New("db error")).Once()

	res, err := service.GetUserFeedByUserId(ctx, 1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}

func TestService_GetDeafultFeed_Success(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()

	defaultTweets := []entity.Tweet{
		{ID: 2, Content: "t2"},
	}

	mockDB.On("GetAllTweets", mock.Anything, 10, 0).Return(defaultTweets, nil).Once()

	mockTweet.On("BuildEntityTweetToResponse", mock.Anything, mock.AnythingOfType("*entity.Tweet")).Return(&entity.Tweet{}, nil).Once()

	res, err := service.GetDeafultFeed(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}

func TestService_GetDeafultFeed_GetError(t *testing.T) {
	mockDB := &mockFeedStorage{}
	mockTweet := &mockTweetService{}

	service := feed.NewFeedService(mockDB, mockTweet)

	ctx := context.Background()

	mockDB.On("GetAllTweets", mock.Anything, 10, 0).Return([]entity.Tweet{}, errors.New("db error")).Once()

	res, err := service.GetDeafultFeed(ctx, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, res)

	mockDB.AssertExpectations(t)
	mockTweet.AssertExpectations(t)
}
