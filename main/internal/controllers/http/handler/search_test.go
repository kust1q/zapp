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

func TestHandler_Search_MissingQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockSearch := mocks.NewMockclientSearchService(c)
	h := ctrl.NewHandler(nil, nil, nil, mockSearch, nil, nil, nil, nil)
	r := h.InitRouters()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/search", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Search_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockSearch := mocks.NewMockclientSearchService(c)
	h := ctrl.NewHandler(nil, nil, nil, mockSearch, nil, nil, nil, nil)
	r := h.InitRouters()

	users := []entity.User{{ID: 1, Username: "test"}}
	tweets := []entity.Tweet{{ID: 2, Content: "twt"}}

	mockSearch.EXPECT().SearchUsers(gomock.Any(), "hello").Return(users, nil)
	mockSearch.EXPECT().SearchTweets(gomock.Any(), "hello").Return(tweets, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/search?query=hello", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Search_ErrorUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockSearch := mocks.NewMockclientSearchService(c)
	h := ctrl.NewHandler(nil, nil, nil, mockSearch, nil, nil, nil, nil)
	r := h.InitRouters()

	mockSearch.EXPECT().SearchUsers(gomock.Any(), "hello").Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/search?query=hello", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_Search_ErrorTweets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()

	mockSearch := mocks.NewMockclientSearchService(c)
	h := ctrl.NewHandler(nil, nil, nil, mockSearch, nil, nil, nil, nil)
	r := h.InitRouters()

	mockSearch.EXPECT().SearchUsers(gomock.Any(), "hello").Return([]entity.User{}, nil)
	mockSearch.EXPECT().SearchTweets(gomock.Any(), "hello").Return(nil, errors.New("err"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/public/search?query=hello", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
