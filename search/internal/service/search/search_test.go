package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/search/internal/domain/entity"
	"github.com/kust1q/Zapp/search/internal/service/search"
)

type mockSearchRepo struct {
	mock.Mock
}

func (m *mockSearchRepo) SearchTweets(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockSearchRepo) SearchUsers(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockSearchRepo) IndexTweet(ctx context.Context, tweet *entity.Tweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}

func (m *mockSearchRepo) IndexUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockSearchRepo) DeleteTweet(ctx context.Context, tweetID int) error {
	args := m.Called(ctx, tweetID)
	return args.Error(0)
}

func (m *mockSearchRepo) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockSearchRepo) DeleteTweetsByUserID(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestSearchService_SearchTweets(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	ids := []int{1, 2, 3}
	repo.On("SearchTweets", mock.Anything, "query").Return(ids, nil)

	result, err := srv.SearchTweets(ctx, "query")
	assert.NoError(t, err)
	assert.Equal(t, ids, result)
	repo.AssertExpectations(t)
}

func TestSearchService_SearchTweets_Error(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	repo.On("SearchTweets", mock.Anything, "query").Return(nil, errors.New("err"))

	result, err := srv.SearchTweets(ctx, "query")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSearchService_SearchUsers(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	ids := []int{1, 2}
	repo.On("SearchUsers", mock.Anything, "query").Return(ids, nil)

	result, err := srv.SearchUsers(ctx, "query")
	assert.NoError(t, err)
	assert.Equal(t, ids, result)
}

func TestSearchService_IndexTweet(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	tweet := &entity.Tweet{ID: 1}
	repo.On("IndexTweet", mock.Anything, tweet).Return(nil)

	err := srv.IndexTweet(ctx, tweet)
	assert.NoError(t, err)
}

func TestSearchService_IndexUser(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	user := &entity.User{ID: 1}
	repo.On("IndexUser", mock.Anything, user).Return(nil)

	err := srv.IndexUser(ctx, user)
	assert.NoError(t, err)
}

func TestSearchService_DeleteTweet(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	repo.On("DeleteTweet", mock.Anything, 1).Return(nil)

	err := srv.DeleteTweet(ctx, 1)
	assert.NoError(t, err)
}

func TestSearchService_DeleteUserWithTweets(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	repo.On("DeleteUser", mock.Anything, 1).Return(nil)
	repo.On("DeleteTweetsByUserID", mock.Anything, 1).Return(nil)

	err := srv.DeleteUserWithTweets(ctx, 1)
	assert.NoError(t, err)
}

func TestSearchService_DeleteUserWithTweets_Error(t *testing.T) {
	repo := new(mockSearchRepo)
	srv := search.NewSearchService(repo)
	ctx := context.Background()

	repo.On("DeleteUser", mock.Anything, 1).Return(errors.New("err"))

	err := srv.DeleteUserWithTweets(ctx, 1)
	assert.Error(t, err)
}
