package http_test

import (
	"bytes"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	ctrl "github.com/kust1q/Zapp/main/internal/controllers/http/handler"
	"github.com/kust1q/Zapp/main/internal/controllers/http/handler/mocks"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
)

func TestHandler_SignUp_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	user := &entity.User{ID: 1, Username: "test"}
	mockAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(user, nil)

	reqBody := `{"email":"test@test.com", "username":"test", "password":"password123", "gen":"male"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_SignUp_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SignUp_EmailUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errs.ErrEmailAlreadyUsed)

	reqBody := `{"email":"used@test.com", "username":"test", "password":"password123", "gen":"male"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestHandler_SignUp_UsernameUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errs.ErrUsernameAlreadyUsed)

	reqBody := `{"email":"test@test.com", "username":"used", "password":"password123", "gen":"male"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestHandler_SignUp_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errors.New("err"))

	reqBody := `{"email":"test@test.com", "username":"test", "password":"password123", "gen":"male"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_SignIn_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)

	mockAuth.EXPECT().GetRefreshTTL().Return(time.Hour).AnyTimes()

	_ = h.InitRouters()

	tokens := &entity.Tokens{
		Access:  &entity.Access{Access: "access-token-longer-than-ten-chars"},
		Refresh: &entity.Refresh{Refresh: "refresh-token-longer-than-ten-chars"},
	}
	mockAuth.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(tokens, nil).AnyTimes()
}

func TestHandler_Refresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	tokens := &entity.Tokens{
		Access:  &entity.Access{Access: "access-token-longer-than-ten-chars"},
		Refresh: &entity.Refresh{Refresh: "refresh-token-longer-than-ten-chars"},
	}
	mockAuth.EXPECT().Refresh(gomock.Any(), gomock.Any()).Return(tokens, nil)
	mockAuth.EXPECT().GetRefreshTTL().Return(time.Hour).AnyTimes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "old-refresh-token-longer-than-ten-chars"})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_SignOut_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignOut(gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/auth/sign-out", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token-longer-than-ten-chars"})
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_ForgotPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().ForgotPassword(gomock.Any(), gomock.Any()).Return(&entity.Recovery{Recovery: "token"}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewBufferString(`{"email":"test@test.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_UpdatePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockAuth.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/protected/reset-password", bytes.NewBufferString(`{"old_password":"oldpassword123","new_password":"newpassword123"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_SignUp_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errors.New("db err"))

	reqBody := `{"email":"test@test.com", "username":"test", "password":"password123", "gen":"male"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_SignIn_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, nil, nil)
	r := h.InitRouters()

	mockAuth.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(nil, errs.ErrInvalidCredentials)

	reqBody := `{"email":"test@test.com", "password":"wrongpassword123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
