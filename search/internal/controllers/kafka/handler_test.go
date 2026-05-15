package kafka_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/search/internal/controllers/kafka"
	"github.com/kust1q/Zapp/search/internal/domain/entity"
	"github.com/kust1q/Zapp/search/internal/domain/events"
)

type mockSearchService struct {
	mock.Mock
}

func (m *mockSearchService) IndexTweet(ctx context.Context, tweet *entity.Tweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}

func (m *mockSearchService) IndexUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockSearchService) DeleteTweet(ctx context.Context, tweetID int) error {
	args := m.Called(ctx, tweetID)
	return args.Error(0)
}

func (m *mockSearchService) DeleteUserWithTweets(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockSearchService) SearchTweets(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockSearchService) SearchUsers(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]int), args.Error(1)
}

func TestSearchHandler_HandleTweetCreate(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	event := events.TweetEvent{
		EventType: events.TweetCreateEvent,
		ID:        1,
		Content:   "test",
		UserID:    1,
		Username:  "user",
	}
	data, _ := json.Marshal(event)

	svc.On("IndexTweet", mock.Anything, mock.Anything).Return(nil)

	err := handler.Handle(ctx, events.TopicTweet, data)
	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSearchHandler_HandleUserCreate(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	event := events.UserEvent{
		EventType: events.UserCreateEvent,
		ID:        1,
		Username:  "user",
	}
	data, _ := json.Marshal(event)

	svc.On("IndexUser", mock.Anything, mock.Anything).Return(nil)

	err := handler.Handle(ctx, events.TopicUser, data)
	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSearchHandler_HandleTweetDelete(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	event := events.TweetDeleted{
		EventType: events.TweetDeleteEvent,
		ID:        1,
	}
	data, _ := json.Marshal(event)

	svc.On("DeleteTweet", mock.Anything, 1).Return(nil)

	err := handler.Handle(ctx, events.TopicTweet, data)
	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSearchHandler_HandleUserDelete(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	event := events.UserDeleted{
		EventType: events.UserDeleteEvent,
		ID:        1,
	}
	data, _ := json.Marshal(event)

	svc.On("DeleteUserWithTweets", mock.Anything, 1).Return(nil)

	err := handler.Handle(ctx, events.TopicUser, data)
	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

func TestSearchHandler_HandleInvalidTopic(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	err := handler.Handle(ctx, "invalid-topic", []byte("data"))
	assert.NoError(t, err)
}

func TestSearchHandler_HandleInvalidData(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	err := handler.Handle(ctx, events.TopicTweet, []byte("invalid-json"))
	assert.Error(t, err)

	err = handler.Handle(ctx, events.TopicUser, []byte("invalid-json"))
	assert.Error(t, err)
}

func TestSearchHandler_HandleTweetUnknownType(t *testing.T) {
	svc := new(mockSearchService)
	handler := kafka.NewSearchHandler(svc)
	ctx := context.Background()

	data := []byte(`{"event_type":"unknown"}`)
	err := handler.Handle(ctx, events.TopicTweet, data)
	assert.NoError(t, err)
}
