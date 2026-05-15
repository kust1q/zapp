package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/service/search"
)

type mockMediaService struct {
	mock.Mock
}

func (m *mockMediaService) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	args := m.Called(ctx, tweetID)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

type mockTweetService struct {
	mock.Mock
}

func (m *mockTweetService) BuildEntityTweetToResponse(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	args := m.Called(ctx, tweet)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tweet), args.Error(1)
}

type mockSearchStorage struct {
	mock.Mock
}

func (m *mockSearchStorage) GetTweetsByIDs(ctx context.Context, ids []int) ([]entity.Tweet, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

func (m *mockSearchStorage) GetUsersByIDs(ctx context.Context, ids []int) ([]entity.User, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.User), args.Error(1)
}

type mockSearchProvider struct {
	mock.Mock
}

func (m *mockSearchProvider) SearchTweets(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockSearchProvider) SearchUsers(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func TestService_SearchTweets_Success(t *testing.T) {
	mockDB := new(mockSearchStorage)
	mockMedia := new(mockMediaService)
	mockTweet := new(mockTweetService)
	mockSearch := new(mockSearchProvider)

	srv := search.NewSearchService(mockDB, mockMedia, mockTweet, mockSearch)

	ids := []int{1, 2}
	tweets := []entity.Tweet{{ID: 1}, {ID: 2}}

	mockSearch.On("SearchTweets", mock.Anything, "q").Return(ids, nil)
	mockDB.On("GetTweetsByIDs", mock.Anything, ids).Return(tweets, nil)
	mockTweet.On("BuildEntityTweetToResponse", mock.Anything, mock.AnythingOfType("*entity.Tweet")).Return(&entity.Tweet{}, nil).Twice()

	res, err := srv.SearchTweets(context.Background(), "q")

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestService_SearchTweets_Empty(t *testing.T) {
	mockDB := new(mockSearchStorage)
	mockMedia := new(mockMediaService)
	mockTweet := new(mockTweetService)
	mockSearch := new(mockSearchProvider)

	srv := search.NewSearchService(mockDB, mockMedia, mockTweet, mockSearch)

	mockSearch.On("SearchTweets", mock.Anything, "q").Return([]int{}, nil)

	res, err := srv.SearchTweets(context.Background(), "q")

	assert.NoError(t, err)
	assert.Empty(t, res)
}

func TestService_SearchTweets_Err(t *testing.T) {
	mockDB := new(mockSearchStorage)
	mockMedia := new(mockMediaService)
	mockTweet := new(mockTweetService)
	mockSearch := new(mockSearchProvider)
	srv := search.NewSearchService(mockDB, mockMedia, mockTweet, mockSearch)
	mockSearch.On("SearchTweets", mock.Anything, "q").Return(nil, errors.New("err"))

	res, err := srv.SearchTweets(context.Background(), "q")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_SearchUsers_Success(t *testing.T) {
	mockDB := new(mockSearchStorage)
	mockMedia := new(mockMediaService)
	mockTweet := new(mockTweetService)
	mockSearch := new(mockSearchProvider)

	srv := search.NewSearchService(mockDB, mockMedia, mockTweet, mockSearch)

	ids := []int{1, 2}
	users := []entity.User{{ID: 1}, {ID: 2}}

	mockSearch.On("SearchUsers", mock.Anything, "q").Return(ids, nil)
	mockDB.On("GetUsersByIDs", mock.Anything, ids).Return(users, nil)
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, mock.AnythingOfType("int")).Return("url", nil).Twice()

	res, err := srv.SearchUsers(context.Background(), "q")

	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestService_SearchUsers_Empty(t *testing.T) {
	mockDB := new(mockSearchStorage)
	mockMedia := new(mockMediaService)
	mockTweet := new(mockTweetService)
	mockSearch := new(mockSearchProvider)

	srv := search.NewSearchService(mockDB, mockMedia, mockTweet, mockSearch)

	mockSearch.On("SearchUsers", mock.Anything, "q").Return([]int{}, nil)

	res, err := srv.SearchUsers(context.Background(), "q")

	assert.NoError(t, err)
	assert.Empty(t, res)
}
