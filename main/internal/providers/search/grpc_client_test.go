package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	searchproto "github.com/kust1q/Zapp/main/pkg/gen/proto/search"
)

type mockSearchClient struct {
	mock.Mock
}

func (m *mockSearchClient) SearchTweets(ctx context.Context, in *searchproto.SearchTweetsRequest, opts ...grpc.CallOption) (*searchproto.SearchTweetsResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchproto.SearchTweetsResponse), args.Error(1)
}

func (m *mockSearchClient) SearchUsers(ctx context.Context, in *searchproto.SearchUsersRequest, opts ...grpc.CallOption) (*searchproto.SearchUsersResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*searchproto.SearchUsersResponse), args.Error(1)
}

func TestClientSearch_SearchTweets(t *testing.T) {
	mockClient := new(mockSearchClient)
	s := &clientSearchService{client: mockClient}
	ctx := context.Background()

	mockClient.On("SearchTweets", mock.Anything, &searchproto.SearchTweetsRequest{Query: "q"}).
		Return(&searchproto.SearchTweetsResponse{TweetIds: []int64{1, 2}}, nil)

	res, err := s.SearchTweets(ctx, "q")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, res)
}

func TestClientSearch_SearchUsers(t *testing.T) {
	mockClient := new(mockSearchClient)
	s := &clientSearchService{client: mockClient}
	ctx := context.Background()

	mockClient.On("SearchUsers", mock.Anything, &searchproto.SearchUsersRequest{Query: "q"}).
		Return(&searchproto.SearchUsersResponse{UserIds: []int64{10}}, nil)

	res, err := s.SearchUsers(ctx, "q")
	assert.NoError(t, err)
	assert.Equal(t, []int{10}, res)
}
