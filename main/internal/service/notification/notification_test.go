package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/service/notification"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	args := m.Called(ctx, tweetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tweet), args.Error(1)
}

func (m *mockDB) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type mockHub struct {
	mock.Mock
}

func (m *mockHub) SendNotification(notif *entity.Notification) {
	m.Called(notif)
}

func TestService_NotifyLike_Success(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	tweet := &entity.Tweet{
		ID:      1,
		Content: "test",
		Author:  &entity.SmallUser{ID: 2},
	}
	user := &entity.User{
		ID:        3,
		Username:  "actor",
		AvatarUrl: "url",
	}

	dbMock.On("GetTweetById", mock.Anything, 1).Return(tweet, nil)
	dbMock.On("GetUserByID", mock.Anything, 3).Return(user, nil)
	hubMock.On("SendNotification", mock.AnythingOfType("*entity.Notification")).Return()

	err := srv.NotifyLike(context.Background(), 3, 1)

	assert.NoError(t, err)
	dbMock.AssertExpectations(t)
	hubMock.AssertExpectations(t)
}

func TestService_NotifyLike_Self(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	tweet := &entity.Tweet{
		ID:      1,
		Content: "test",
		Author:  &entity.SmallUser{ID: 2},
	}

	dbMock.On("GetTweetById", mock.Anything, 1).Return(tweet, nil)

	err := srv.NotifyLike(context.Background(), 2, 1)

	assert.NoError(t, err)
	dbMock.AssertExpectations(t)
	hubMock.AssertExpectations(t)
}

func TestService_NotifyLike_ErrTweet(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)
	srv := notification.NewNotificationService(hubMock, dbMock)

	dbMock.On("GetTweetById", mock.Anything, 1).Return(nil, errors.New("err"))

	err := srv.NotifyLike(context.Background(), 3, 1)

	assert.Error(t, err)
}

func TestService_NotifyRetweet_Success(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	tweet := &entity.Tweet{
		ID:      1,
		Content: "test",
		Author:  &entity.SmallUser{ID: 2},
	}
	user := &entity.User{
		ID:        3,
		Username:  "actor",
		AvatarUrl: "url",
	}

	dbMock.On("GetTweetById", mock.Anything, 1).Return(tweet, nil)
	dbMock.On("GetUserByID", mock.Anything, 3).Return(user, nil)
	hubMock.On("SendNotification", mock.AnythingOfType("*entity.Notification")).Return()

	err := srv.NotifyRetweet(context.Background(), 3, 1)

	assert.NoError(t, err)
}

func TestService_NotifyReply_Success(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	tweet := &entity.Tweet{
		ID:      1,
		Content: "test",
		Author:  &entity.SmallUser{ID: 2},
	}
	user := &entity.User{
		ID:        3,
		Username:  "actor",
		AvatarUrl: "url",
	}

	dbMock.On("GetTweetById", mock.Anything, 1).Return(tweet, nil)
	dbMock.On("GetUserByID", mock.Anything, 3).Return(user, nil)
	hubMock.On("SendNotification", mock.AnythingOfType("*entity.Notification")).Return()

	err := srv.NotifyReply(context.Background(), 3, 1)

	assert.NoError(t, err)
}

func TestService_NotifyFollow_Success(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	user := &entity.User{
		ID:        3,
		Username:  "actor",
		AvatarUrl: "url",
	}

	dbMock.On("GetUserByID", mock.Anything, 3).Return(user, nil)
	hubMock.On("SendNotification", mock.AnythingOfType("*entity.Notification")).Return()

	err := srv.NotifyFollow(context.Background(), 3, 1)

	assert.NoError(t, err)
}

func TestService_NotifyFollow_ErrUser(t *testing.T) {
	dbMock := new(mockDB)
	hubMock := new(mockHub)

	srv := notification.NewNotificationService(hubMock, dbMock)

	dbMock.On("GetUserByID", mock.Anything, 3).Return(nil, errors.New("err"))

	err := srv.NotifyFollow(context.Background(), 3, 1)

	assert.Error(t, err)
}
