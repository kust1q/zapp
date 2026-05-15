package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	ctrl "github.com/kust1q/Zapp/main/internal/controllers/http/handler"
	"github.com/kust1q/Zapp/main/internal/controllers/http/handler/mocks"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func TestHandler_DeleteTweetMedia_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	userID := 1
	tweetID := 2

	mockAuth.EXPECT().VerifyAccessToken("valid").Return(userID, nil)
	mockMedia.EXPECT().DeleteTweetMedia(gomock.Any(), tweetID, userID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/2/media", nil)
	req.Header.Set("Authorization", "Bearer valid")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_DeleteTweetMedia_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/2/media", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_DeleteTweetMedia_InvalidTweetID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	userID := 1

	mockAuth.EXPECT().VerifyAccessToken("valid").Return(userID, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/invalid/media", nil)
	req.Header.Set("Authorization", "Bearer valid")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_DeleteTweetMedia_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	userID := 1
	tweetID := 2

	mockAuth.EXPECT().VerifyAccessToken("valid").Return(userID, nil)
	mockMedia.EXPECT().DeleteTweetMedia(gomock.Any(), tweetID, userID).Return(errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/2/media", nil)
	req.Header.Set("Authorization", "Bearer valid")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_GetTweetMedia_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	tweetID := 2
	media := &entity.TweetMedia{ID: 1}

	mockMedia.EXPECT().GetMediaDataByTweetID(gomock.Any(), tweetID).Return(media, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/media/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetTweetMedia_InvalidTweetID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/media/invalid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetTweetMedia_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	tweetID := 2

	mockMedia.EXPECT().GetMediaDataByTweetID(gomock.Any(), tweetID).Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/media/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_GetAvatar_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	userID := 1
	avatar := &entity.Avatar{ID: 2}

	mockMedia.EXPECT().GetAvatarDataByUserID(gomock.Any(), userID).Return(avatar, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/users/avatar/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetAvatar_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/users/avatar/invalid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_GetAvatar_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockMedia := mocks.NewMockmediaService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, nil, mockMedia, nil, nil)
	r := h.InitRouters()

	userID := 1

	mockMedia.EXPECT().GetAvatarDataByUserID(gomock.Any(), userID).Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/users/avatar/1", nil)
	r.ServeHTTP(w, req)

}
