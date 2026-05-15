package searchgrpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	searchgrpc "github.com/kust1q/Zapp/search/internal/controllers/grpc/servers/search"
	searchproto "github.com/kust1q/Zapp/search/pkg/gen/proto/search"
)

type mockSearchService struct {
	mock.Mock
}

func (m *mockSearchService) SearchTweets(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockSearchService) SearchUsers(ctx context.Context, query string) ([]int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func TestSearchServer_SearchUsers_Success(t *testing.T) {
	svc := new(mockSearchService)
	server := searchgrpc.NewSearchServer(svc)
	ctx := context.Background()

	ids := []int{1, 2}
	svc.On("SearchUsers", mock.Anything, "query").Return(ids, nil)

	resp, err := server.SearchUsers(ctx, &searchproto.SearchUsersRequest{Query: "query"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.UserIds[0])
	assert.Equal(t, int64(2), resp.UserIds[1])
}

func TestSearchServer_SearchUsers_Error(t *testing.T) {
	svc := new(mockSearchService)
	server := searchgrpc.NewSearchServer(svc)
	ctx := context.Background()

	svc.On("SearchUsers", mock.Anything, "query").Return(nil, errors.New("err"))

	resp, err := server.SearchUsers(ctx, &searchproto.SearchUsersRequest{Query: "query"})
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestSearchServer_SearchTweets_Success(t *testing.T) {
	svc := new(mockSearchService)
	server := searchgrpc.NewSearchServer(svc)
	ctx := context.Background()

	ids := []int{10, 20}
	svc.On("SearchTweets", mock.Anything, "query").Return(ids, nil)

	resp, err := server.SearchTweets(ctx, &searchproto.SearchTweetsRequest{Query: "query"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(10), resp.TweetIds[0])
	assert.Equal(t, int64(20), resp.TweetIds[1])
}

func TestSearchServer_SearchTweets_Error(t *testing.T) {
	svc := new(mockSearchService)
	server := searchgrpc.NewSearchServer(svc)
	ctx := context.Background()

	svc.On("SearchTweets", mock.Anything, "query").Return(nil, errors.New("err"))

	resp, err := server.SearchTweets(ctx, &searchproto.SearchTweetsRequest{Query: "query"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}
