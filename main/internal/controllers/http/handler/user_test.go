package http_test

import (
	"bytes"
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

func TestHandler_GetMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockUser := mocks.NewMockuserService(c)
	h := ctrl.NewHandler(mockAuth, nil, mockUser, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockUser.EXPECT().GetMe(gomock.Any(), userID, 10, 0).Return(&entity.UserProfile{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UpdateMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockUser := mocks.NewMockuserService(c)
	h := ctrl.NewHandler(mockAuth, nil, mockUser, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockUser.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/protected/users/me", bytes.NewBufferString(`{"bio":"New bio"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_FollowUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockUser := mocks.NewMockuserService(c)
	mockNotif := mocks.NewMocknotificationService(c)
	h := ctrl.NewHandler(mockAuth, nil, mockUser, nil, nil, nil, nil, mockNotif)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockUser.EXPECT().FollowToUser(gomock.Any(), userID, 2).Return(&entity.Follow{}, nil)
	mockNotif.EXPECT().NotifyFollow(gomock.Any(), userID, 2).Return(nil).AnyTimes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/protected/users/2/follow", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UnfollowUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockUser := mocks.NewMockuserService(c)
	h := ctrl.NewHandler(mockAuth, nil, mockUser, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockUser.EXPECT().UnfollowUser(gomock.Any(), userID, 2).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/users/2/follow", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetUserProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockUser := mocks.NewMockuserService(c)
	h := ctrl.NewHandler(nil, nil, mockUser, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockUser.EXPECT().GetUserProfile(gomock.Any(), "testuser", 10, 0).Return(&entity.UserProfile{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/users/testuser/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_DeleteMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockUser := mocks.NewMockuserService(c)
	h := ctrl.NewHandler(mockAuth, nil, mockUser, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockUser.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/protected/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
