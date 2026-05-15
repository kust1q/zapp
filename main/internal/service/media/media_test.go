package media_test

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/service/media"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, tweetMedia *entity.TweetMedia) (*entity.TweetMedia, error) {
	args := m.Called(ctx, tx, tweetMedia)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TweetMedia), args.Error(1)
}

func (m *mockDB) GetMediaPathByTweetID(ctx context.Context, tweetID int) (string, error) {
	args := m.Called(ctx, tweetID)
	return args.String(0), args.Error(1)
}

func (m *mockDB) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	args := m.Called(ctx, tweetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TweetMedia), args.Error(1)
}

func (m *mockDB) DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error {
	args := m.Called(ctx, tweetID, userID)
	return args.Error(0)
}

func (m *mockDB) UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error) {
	args := m.Called(ctx, tx, avatar)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Avatar), args.Error(1)
}

func (m *mockDB) GetAvatarPathByUserID(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockDB) GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Avatar), args.Error(1)
}

func (m *mockDB) DeleteAvatarByUserID(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockDB) GetMediaUrlsByUserID(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type mockObject struct {
	mock.Mock
}

func (m *mockObject) Upload(ctx context.Context, data io.Reader, mt entity.MediaType, name string) (string, string, error) {
	args := m.Called(ctx, data, mt, name)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockObject) GetPresignedURL(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

func (m *mockObject) Remove(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func TestService_GetMediaUrlByTweetID_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetMediaPathByTweetID", mock.Anything, 1).Return("somepath", nil)
	objMock.On("GetPresignedURL", mock.Anything, "somepath").Return("http://url", nil)

	url, err := srv.GetMediaUrlByTweetID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "http://url", url)
}

func TestService_GetMediaDataByTweetID_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	mediaEntity := &entity.TweetMedia{Path: "somepath"}

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(mediaEntity, nil)
	objMock.On("GetPresignedURL", mock.Anything, "somepath").Return("http://url", nil)

	res, err := srv.GetMediaDataByTweetID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "http://url", res.Path)
}

func TestService_UploadAndAttachTweetMediaTx_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	file := strings.NewReader("dummy content")
	filename := "test.jpg"

	db, smock, _ := sqlmock.New()
	defer db.Close()
	tx, _ := db.Begin()

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeImage, filename).Return("path/to/media", "image/jpeg", nil)
	dbMock.On("UpsertByTweetIdTx", mock.Anything, tx, mock.Anything).Return(&entity.TweetMedia{Path: "path/to/media"}, nil)
	objMock.On("GetPresignedURL", mock.Anything, "path/to/media").Return("http://url", nil)

	url, err := srv.UploadAndAttachTweetMediaTx(ctx, 1, file, filename, tx)

	assert.NoError(t, err)
	assert.Equal(t, "http://url", url)
	smock.ExpectRollback()
}

func TestService_UploadAvatarTx_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	file := strings.NewReader("dummy content")
	filename := "avatar.png"

	db, smock, _ := sqlmock.New()
	defer db.Close()
	tx, _ := db.Begin()

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeImage, filename).Return("path/to/avatar", "image/png", nil)
	dbMock.On("UploadAvatarTx", mock.Anything, tx, mock.Anything).Return(&entity.Avatar{Path: "path/to/avatar"}, nil)
	objMock.On("GetPresignedURL", mock.Anything, "path/to/avatar").Return("http://url", nil)

	res, err := srv.UploadAvatarTx(ctx, 1, file, filename, tx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "http://url", res.Path)
	smock.ExpectRollback()
}

func TestService_GetAvatarUrlByUserID_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetAvatarPathByUserID", mock.Anything, 1).Return("avatarpath", nil)
	objMock.On("GetPresignedURL", mock.Anything, "avatarpath").Return("http://url", nil)

	url, err := srv.GetAvatarUrlByUserID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "http://url", url)
}

func TestService_DeleteAvatar_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	avatarEntity := &entity.Avatar{Path: "avatarpath"}

	dbMock.On("GetAvatarDataByUserID", mock.Anything, 1).Return(avatarEntity, nil)
	dbMock.On("DeleteAvatarByUserID", mock.Anything, 1).Return(nil)
	objMock.On("Remove", mock.Anything, "avatarpath").Return(nil)

	err := srv.DeleteAvatar(ctx, 1)
	assert.NoError(t, err)
	
	time.Sleep(100 * time.Millisecond)
}

func TestService_DeleteMediasByUserID_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	urls := []string{"path1", "path2"}

	dbMock.On("GetMediaUrlsByUserID", mock.Anything, 1).Return(urls, nil)
	objMock.On("Remove", mock.Anything, "path1").Return(nil)
	objMock.On("Remove", mock.Anything, "path2").Return(nil)

	err := srv.DeleteMediasByUserID(ctx, 1)
	assert.NoError(t, err)
}

func TestService_GetPresignedURL(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	objMock.On("GetPresignedURL", mock.Anything, "path").Return("http://url", nil)

	url, err := srv.GetPresignedURL(ctx, "path")
	assert.NoError(t, err)
	assert.Equal(t, "http://url", url)
}

func TestService_GetMediaDataByTweetID_Error(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(nil, errors.New("db error"))

	res, err := srv.GetMediaDataByTweetID(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_GetMediaDataByTweetID_NoRows(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(nil, sql.ErrNoRows)

	res, err := srv.GetMediaDataByTweetID(ctx, 1)

	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestService_DeleteTweetMedia_Error(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(nil, errors.New("err"))

	err := srv.DeleteTweetMedia(ctx, 1, 2)
	assert.Error(t, err)
}

func TestService_UploadAvatarTx_UploadError(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	file := strings.NewReader("dummy content")
	filename := "avatar.png"

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeImage, filename).Return("", "", errors.New("upload err"))

	res, err := srv.UploadAvatarTx(ctx, 1, file, filename, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestService_DeleteAvatar_Error(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	dbMock.On("GetAvatarDataByUserID", mock.Anything, 1).Return(nil, errors.New("err"))

	err := srv.DeleteAvatar(ctx, 1)
	assert.Error(t, err)
}

func TestService_GetAvatarDataByUserID_Success(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	avatarEntity := &entity.Avatar{ID: 1, UserID: 1, Path: "path"}

	dbMock.On("GetAvatarDataByUserID", mock.Anything, 1).Return(avatarEntity, nil)
	dbMock.On("GetAvatarPathByUserID", mock.Anything, 1).Return("path", nil)
	objMock.On("GetPresignedURL", mock.Anything, "path").Return("http://url", nil)

	res, err := srv.GetAvatarDataByUserID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "http://url", res.Path)
}

func TestService_DetectMediaType(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)
	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	tests := []struct {
		filename string
		expected entity.MediaType
		wantErr  bool
	}{
		{"test.jpg", entity.MediaTypeImage, false},
		{"test.jpeg", entity.MediaTypeImage, false},
		{"test.png", entity.MediaTypeImage, false},
		{"test.webp", entity.MediaTypeImage, false},
		{"test.gif", entity.MediaTypeGIF, false},
		{"test.mp4", entity.MediaTypeVideo, false},
		{"test.mov", entity.MediaTypeVideo, false},
		{"test.m4v", entity.MediaTypeVideo, false},
		{"test.mp3", entity.MediaTypeAudio, false},
		{"test.wav", entity.MediaTypeAudio, false},
		{"test.ogg", entity.MediaTypeAudio, false},
		{"test.flac", entity.MediaTypeAudio, false},
		{"test.aac", entity.MediaTypeAudio, false},
		{"test.m4a", entity.MediaTypeAudio, false},
		{"test.webm", entity.MediaTypeAudio, false},
		{"test.txt", "", true},
	}

	db, smock, _ := sqlmock.New()
	defer db.Close()
	tx, _ := db.Begin()
	smock.ExpectRollback()

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			if tt.wantErr {
				url, err := srv.UploadAndAttachTweetMediaTx(ctx, 1, strings.NewReader("content"), tt.filename, tx)
				assert.Error(t, err)
				assert.Empty(t, url)
				return
			}

			objMock.On("Upload", mock.Anything, mock.Anything, tt.expected, tt.filename).Return("path", "mime", nil).Once()
			dbMock.On("UpsertByTweetIdTx", mock.Anything, tx, mock.Anything).Return(&entity.TweetMedia{Path: "path"}, nil).Once()
			objMock.On("GetPresignedURL", mock.Anything, "path").Return("url", nil).Once()

			url, err := srv.UploadAndAttachTweetMediaTx(ctx, 1, strings.NewReader("content"), tt.filename, tx)
			assert.NoError(t, err)
			assert.NotEmpty(t, url)
		})
	}
}

func TestService_UploadAndAttachTweetMediaTx_VariousTypes(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)
	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()
	
	db, smock, _ := sqlmock.New()
	defer db.Close()
	tx, _ := db.Begin()
	smock.ExpectRollback()

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeVideo, "test.mp4").Return("vpath", "video/mp4", nil).Once()
	dbMock.On("UpsertByTweetIdTx", mock.Anything, tx, mock.Anything).Return(&entity.TweetMedia{Path: "vpath"}, nil).Once()
	objMock.On("GetPresignedURL", mock.Anything, "vpath").Return("vurl", nil).Once()

	url, err := srv.UploadAndAttachTweetMediaTx(ctx, 1, strings.NewReader("content"), "test.mp4", tx)
	assert.NoError(t, err)
	assert.Equal(t, "vurl", url)

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeGIF, "test.gif").Return("gpath", "image/gif", nil).Once()
	dbMock.On("UpsertByTweetIdTx", mock.Anything, tx, mock.Anything).Return(&entity.TweetMedia{Path: "gpath"}, nil).Once()
	objMock.On("GetPresignedURL", mock.Anything, "gpath").Return("gurl", nil).Once()

	url, err = srv.UploadAndAttachTweetMediaTx(ctx, 1, strings.NewReader("content"), "test.gif", tx)
	assert.NoError(t, err)
	assert.Equal(t, "gurl", url)

	objMock.On("Upload", mock.Anything, mock.Anything, entity.MediaTypeAudio, "test.mp3").Return("apath", "audio/mpeg", nil).Once()
	dbMock.On("UpsertByTweetIdTx", mock.Anything, tx, mock.Anything).Return(&entity.TweetMedia{Path: "apath"}, nil).Once()
	objMock.On("GetPresignedURL", mock.Anything, "apath").Return("aurl", nil).Once()

	url, err = srv.UploadAndAttachTweetMediaTx(ctx, 1, strings.NewReader("content"), "test.mp3", tx)
	assert.NoError(t, err)
	assert.Equal(t, "aurl", url)
}

func TestService_DeleteTweetMedia_DeleteError(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	mediaEntity := &entity.TweetMedia{Path: "somepath"}

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(mediaEntity, nil)
	dbMock.On("DeleteMediaByTweetID", mock.Anything, 1, 2).Return(errors.New("delete err"))

	err := srv.DeleteTweetMedia(ctx, 1, 2)
	assert.Error(t, err)
}

func TestService_DeleteTweetMedia_EmptyPath(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	mediaEntity := &entity.TweetMedia{Path: ""}

	dbMock.On("GetMediaDataByTweetID", mock.Anything, 1).Return(mediaEntity, nil)
	dbMock.On("DeleteMediaByTweetID", mock.Anything, 1, 2).Return(nil)

	err := srv.DeleteTweetMedia(ctx, 1, 2)
	assert.NoError(t, err)
}

func TestService_CleanUpMedia_Error(t *testing.T) {
	dbMock := new(mockDB)
	objMock := new(mockObject)

	srv := media.NewMediaService(dbMock, objMock)
	ctx := context.Background()

	objMock.On("Remove", mock.Anything, "path").Return(errors.New("remove err"))

	dbMock.On("GetMediaUrlsByUserID", mock.Anything, 1).Return([]string{"path"}, nil)
	
	err := srv.DeleteMediasByUserID(ctx, 1)
	assert.NoError(t, err)
}
