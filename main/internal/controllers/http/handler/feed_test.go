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

func TestHandler_GetDefaultFeed_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockFeed := mocks.NewMockfeedService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, mockFeed, nil, nil, nil)
	r := h.InitRouters()

	tweets := []entity.Tweet{{ID: 1, Content: "test tweet"}}
	mockFeed.EXPECT().GetDeafultFeed(gomock.Any(), 10, 0).Return(tweets, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetDefaultFeed_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockFeed := mocks.NewMockfeedService(c)
	h := ctrl.NewHandler(nil, nil, nil, nil, mockFeed, nil, nil, nil)
	r := h.InitRouters()

	mockFeed.EXPECT().GetDeafultFeed(gomock.Any(), 10, 0).Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_GetFeed_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockFeed := mocks.NewMockfeedService(c)
	mockAuth := mocks.NewMockauthService(c)

	h := ctrl.NewHandler(mockAuth, nil, nil, nil, mockFeed, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	mockAuth.EXPECT().VerifyAccessToken("valid-token").Return(userID, nil)

	tweets := []entity.Tweet{{ID: 2, Content: "test feed"}}
	mockFeed.EXPECT().GetUserFeedByUserId(gomock.Any(), userID, 10, 0).Return(tweets, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected/feed", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetFeed_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockFeed := mocks.NewMockfeedService(c)
	mockAuth := mocks.NewMockauthService(c)

	h := ctrl.NewHandler(mockAuth, nil, nil, nil, mockFeed, nil, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected/feed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_GetFeed_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockFeed := mocks.NewMockfeedService(c)
	mockAuth := mocks.NewMockauthService(c)

	h := ctrl.NewHandler(mockAuth, nil, nil, nil, mockFeed, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	mockAuth.EXPECT().VerifyAccessToken("valid-token").Return(userID, nil)
	mockFeed.EXPECT().GetUserFeedByUserId(gomock.Any(), userID, 10, 0).Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected/feed", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
