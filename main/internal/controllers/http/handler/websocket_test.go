package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	ctrl "github.com/kust1q/Zapp/main/internal/controllers/http/handler"
	"github.com/kust1q/Zapp/main/internal/controllers/http/handler/mocks"
)

func TestHandler_ServeWs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockAuth := mocks.NewMockauthService(c)
	mockWs := mocks.NewMockwebSocketService(c)
	h := ctrl.NewHandler(mockAuth, nil, nil, nil, nil, nil, mockWs, nil)
	r := h.InitRouters()

	userID := 1
	token := "token"

	mockAuth.EXPECT().VerifyAccessToken(token).Return(userID, nil)
	mockWs.EXPECT().HandleConnection(gomock.Any(), gomock.Any(), userID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected/ws", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
