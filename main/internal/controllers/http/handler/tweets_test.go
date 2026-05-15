package http_test

import (
	"bytes"
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
	"github.com/kust1q/Zapp/main/internal/errs"
)

func TestHandler_CreateTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "valid-token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	
	tweet := &entity.Tweet{ID: 10, Content: "Hello", Author: &entity.SmallUser{ID: userID}}
	mockTweets.EXPECT().CreateTweet(gomock.Any(), gomock.Any()).Return(tweet, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/protected/tweets", bytes.NewBufferString(`{"content":"Hello"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_GetTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(nil, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	tweet := &entity.Tweet{ID: 1, Content: "Hello", Author: &entity.SmallUser{ID: 1}}
	mockTweets.EXPECT().GetTweetById(gomock.Any(), 1).Return(tweet, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_DeleteTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockTweets.EXPECT().DeleteTweet(gomock.Any(), userID, 10).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_LikeTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	mockNotif := mocks.NewMocknotificationService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, mockNotif)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockTweets.EXPECT().LikeTweet(gomock.Any(), userID, 10).Return(nil)
	mockNotif.EXPECT().NotifyLike(gomock.Any(), userID, 10).Return(nil).AnyTimes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/protected/tweets/10/like", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UnlikeTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockTweets.EXPECT().UnlikeTweet(gomock.Any(), userID, 10).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/10/like", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetRepliesToTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(nil, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockTweets.EXPECT().GetRepliesToTweet(gomock.Any(), 1, 10, 0).Return([]entity.Tweet{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/1/replies", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_CreateRetweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	mockNotif := mocks.NewMocknotificationService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, mockNotif)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockTweets.EXPECT().CreateRetweet(gomock.Any(), userID, 10).Return(nil)
	mockNotif.EXPECT().NotifyRetweet(gomock.Any(), userID, 10).Return(nil).AnyTimes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/protected/tweets/10/retweet", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetTweetsByUsername_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(nil, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockTweets.EXPECT().GetTweetsAndRetweetsByUsername(gomock.Any(), "testuser", 10, 0).Return([]entity.Tweet{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/users/testuser/tweets", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetTweet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(nil, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockTweets.EXPECT().GetTweetById(gomock.Any(), 1).Return(nil, errs.ErrTweetNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/tweets/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_DeleteTweet_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockTweets := mocks.NewMocktweetService(c)
	h := ctrl.NewHandler(mockAuth, mockTweets, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockTweets.EXPECT().DeleteTweet(gomock.Any(), userID, 10).Return(errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/tweets/10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
