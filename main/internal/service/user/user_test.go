package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kust1q/Zapp/main/internal/service/user"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/domain/events"
	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockUserStorage) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*entity.User), args.Error(1)
}

func (m *mockUserStorage) UpdateUserBio(ctx context.Context, userID int, bio string) error {
	args := m.Called(ctx, userID, bio)
	return args.Error(0)
}

func (m *mockUserStorage) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserStorage) FollowToUser(ctx context.Context, followerID, followingID int, createdAt time.Time) (*entity.Follow, error) {
	args := m.Called(ctx, followerID, followingID, createdAt)
	follow := args.Get(0)
	if follow == nil {
		return nil, args.Error(1)
	}
	return follow.(*entity.Follow), args.Error(1)
}

func (m *mockUserStorage) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	args := m.Called(ctx, followerID, followingID)
	return args.Error(0)
}

func (m *mockUserStorage) GetFollowersIds(ctx context.Context, username string, limit, offset int) ([]int, error) {
	args := m.Called(ctx, username, limit, offset)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockUserStorage) GetFollowingsIds(ctx context.Context, username string, limit, offset int) ([]int, error) {
	args := m.Called(ctx, username, limit, offset)
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockUserStorage) GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error) {
	args := m.Called(ctx, username, limit, offset)
	return args.Get(0).([]entity.Tweet), args.Error(1)
}

func (m *mockUserStorage) GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error) {
	args := m.Called(ctx, tweetID)
	counters := args.Get(0)
	if counters == nil {
		return nil, args.Error(1)
	}
	return counters.(*entity.Counters), args.Error(1)
}

type mockMediaService struct {
	mock.Mock
}

func (m *mockMediaService) GetAvatarUrlByUserID(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) DeleteAvatar(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockMediaService) GetMediaUrlByTweetID(ctx context.Context, tweetID int) (string, error) {
	args := m.Called(ctx, tweetID)
	return args.String(0), args.Error(1)
}

func (m *mockMediaService) DeleteMediasByUserID(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type mockEventProducer struct {
	mock.Mock
}

func (m *mockEventProducer) Publish(ctx context.Context, topic string, event any) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}

func TestService_DeleteUser_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteAvatar", mock.Anything, 1).Return(nil).Once()
	mockMedia.On("DeleteMediasByUserID", mock.Anything, 1).Return(nil).Once()
	mockDB.On("DeleteUser", mock.Anything, 1).Return(nil).Once()
	mockProducer.On("Publish", mock.Anything, events.TopicTweet, mock.AnythingOfType("events.UserDeleted")).Return(nil).Once()

	err := service.DeleteUser(ctx, 1)

	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteUser_MediaErrors(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteAvatar", mock.Anything, 1).Return(errors.New("avatar error")).Once()
	mockMedia.On("DeleteMediasByUserID", mock.Anything, 1).Return(errors.New("medias error")).Once()
	mockDB.On("DeleteUser", mock.Anything, 1).Return(nil).Once()
	mockProducer.On("Publish", mock.Anything, events.TopicTweet, mock.AnythingOfType("events.UserDeleted")).Return(nil).Once()

	err := service.DeleteUser(ctx, 1)

	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteUser_UserNotFound(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteAvatar", mock.Anything, 1).Return(nil).Once()
	mockMedia.On("DeleteMediasByUserID", mock.Anything, 1).Return(nil).Once()
	mockDB.On("DeleteUser", mock.Anything, 1).Return(errs.ErrUserNotFound).Once()

	err := service.DeleteUser(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrUserNotFound, err)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_DeleteUser_DBError(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockMedia.On("DeleteAvatar", mock.Anything, 1).Return(nil).Once()
	mockMedia.On("DeleteMediasByUserID", mock.Anything, 1).Return(nil).Once()
	mockDB.On("DeleteUser", mock.Anything, 1).Return(errors.New("db error")).Once()

	err := service.DeleteUser(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete user")

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestService_FollowToUser_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	follow := &entity.Follow{
		FollowerID:  1,
		FollowingID: 2,
		CreatedAt:   time.Now(),
	}

	mockDB.On("FollowToUser", mock.Anything, 1, 2, mock.AnythingOfType("time.Time")).Return(follow, nil).Once()

	result, err := service.FollowToUser(ctx, 1, 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.FollowerID)
	assert.Equal(t, 2, result.FollowingID)

	mockDB.AssertExpectations(t)
}

func TestService_FollowToUser_SelfFollow(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	result, err := service.FollowToUser(ctx, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "impossible to subscribe to yourself")

	mockDB.AssertExpectations(t)
}

func TestService_UnfollowUser_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("UnfollowUser", mock.Anything, 1, 2).Return(nil).Once()

	err := service.UnfollowUser(ctx, 1, 2)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_GetFollowers_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	followerIDs := []int{2, 3}
	user2 := &entity.User{ID: 2, Username: "user2"}
	user3 := &entity.User{ID: 3, Username: "user3"}

	mockDB.On("GetFollowersIds", mock.Anything, "testuser", 10, 0).Return(followerIDs, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 2).Return(user2, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 3).Return(user3, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 2).Return("/avatars/2.jpg", nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 3).Return("/avatars/3.jpg", nil).Once()

	result, err := service.GetFollowers(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "user2", result[0].Username)
	assert.Equal(t, "/avatars/2.jpg", result[0].AvatarUrl)
	assert.Equal(t, "user3", result[1].Username)
	assert.Equal(t, "/avatars/3.jpg", result[1].AvatarUrl)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetFollowers_GetIDsError(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("GetFollowersIds", mock.Anything, "testuser", 10, 0).Return([]int{}, errors.New("db error")).Once()

	result, err := service.GetFollowers(ctx, "testuser", 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get followers ids")

	mockDB.AssertExpectations(t)
}

func TestService_GetFollowers_GetUserError(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	followerIDs := []int{2}
	mockDB.On("GetFollowersIds", mock.Anything, "testuser", 10, 0).Return(followerIDs, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 2).Return(nil, errors.New("user not found")).Once()

	result, err := service.GetFollowers(ctx, "testuser", 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user by id")

	mockDB.AssertExpectations(t)
}

func TestService_GetFollowers_GetAvatarError(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	followerIDs := []int{2}
	user2 := &entity.User{ID: 2, Username: "user2"}

	mockDB.On("GetFollowersIds", mock.Anything, "testuser", 10, 0).Return(followerIDs, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 2).Return(user2, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 2).Return("", errors.New("avatar error")).Once()

	result, err := service.GetFollowers(ctx, "testuser", 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get avatar")

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetFollowings_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	followingIDs := []int{4, 5}
	user4 := &entity.User{ID: 4, Username: "user4"}
	user5 := &entity.User{ID: 5, Username: "user5"}

	mockDB.On("GetFollowingsIds", mock.Anything, "testuser", 10, 0).Return(followingIDs, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 4).Return(user4, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 5).Return(user5, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 4).Return("/avatars/4.jpg", nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 5).Return("/avatars/5.jpg", nil).Once()

	result, err := service.GetFollowings(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "user4", result[0].Username)
	assert.Equal(t, "/avatars/4.jpg", result[0].AvatarUrl)
	assert.Equal(t, "user5", result[1].Username)
	assert.Equal(t, "/avatars/5.jpg", result[1].AvatarUrl)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetUserByID_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Bio:      "test bio",
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()

	result, err := service.GetUserByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test bio", result.Bio)
	assert.Equal(t, "/avatars/1.jpg", result.AvatarUrl)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetUserByID_NotFound(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	mockDB.On("GetUserByID", mock.Anything, 1).Return(nil, errors.New("not found")).Once()

	result, err := service.GetUserByID(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user by id")

	mockDB.AssertExpectations(t)
}

func TestService_GetUserByID_AvatarError(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("", errors.New("avatar error")).Once()

	result, err := service.GetUserByID(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get avatar by user id")

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetUserByUsername_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Bio:      "test bio",
	}

	mockDB.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()

	result, err := service.GetUserByUsername(ctx, "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test bio", result.Bio)
	assert.Equal(t, "/avatars/1.jpg", result.AvatarUrl)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetMe_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	tweets := []entity.Tweet{
		{
			ID:      1,
			Content: "Tweet 1",
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

	mockDB.On("GetUserByID", mock.Anything, 1).Return(user, nil).Once()
	mockDB.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockDB.On("GetTweetsAndRetweetsByUsername", mock.Anything, "testuser", 10, 0).Return(tweets, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 1).Return(author, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, 1).Return("", nil).Once()
	mockDB.On("GetCounts", mock.Anything, 1).Return(counters, nil).Once()

	result, err := service.GetMe(ctx, 1, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Len(t, result.Tweets, 1)
	assert.Equal(t, "Tweet 1", result.Tweets[0].Content)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetUserProfile_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	tweets := []entity.Tweet{
		{
			ID:      1,
			Content: "Tweet 1",
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

	mockDB.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockDB.On("GetTweetsAndRetweetsByUsername", mock.Anything, "testuser", 10, 0).Return(tweets, nil).Once()
	mockDB.On("GetUserByID", mock.Anything, 1).Return(author, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockMedia.On("GetMediaUrlByTweetID", mock.Anything, 1).Return("", nil).Once()
	mockDB.On("GetCounts", mock.Anything, 1).Return(counters, nil).Once()

	result, err := service.GetUserProfile(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Len(t, result.Tweets, 1)
	assert.Equal(t, "Tweet 1", result.Tweets[0].Content)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_GetUserProfile_NoTweets(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	user := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	mockDB.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil).Once()
	mockMedia.On("GetAvatarUrlByUserID", mock.Anything, 1).Return("/avatars/1.jpg", nil).Once()
	mockDB.On("GetTweetsAndRetweetsByUsername", mock.Anything, "testuser", 10, 0).Return([]entity.Tweet{}, nil).Once()

	result, err := service.GetUserProfile(ctx, "testuser", 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Empty(t, result.Tweets)

	mockDB.AssertExpectations(t)
	mockMedia.AssertExpectations(t)
}

func TestService_Update_Success(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	req := &entity.UpdateBio{
		UserID: 1,
		Bio:    "new bio",
	}

	mockDB.On("UpdateUserBio", mock.Anything, 1, "new bio").Return(nil).Once()

	err := service.Update(ctx, req)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestService_Update_Error(t *testing.T) {
	mockDB := &mockUserStorage{}
	mockMedia := &mockMediaService{}
	mockProducer := &mockEventProducer{}

	service := user.NewUserService(mockDB, mockMedia, mockProducer)

	ctx := context.Background()

	req := &entity.UpdateBio{
		UserID: 1,
		Bio:    "new bio",
	}

	mockDB.On("UpdateUserBio", mock.Anything, 1, "new bio").Return(errors.New("db error")).Once()

	err := service.Update(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	mockDB.AssertExpectations(t)
}
