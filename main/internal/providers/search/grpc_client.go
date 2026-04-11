package search

import (
	"context"

	searchproto "github.com/kust1q/Zapp/main/pkg/gen/proto/search"
	"google.golang.org/grpc"
)

type clientSearchService struct {
	client searchproto.SearchServiceClient
}

func NewClientSearchService(conn *grpc.ClientConn) *clientSearchService {
	return &clientSearchService{
		client: searchproto.NewSearchServiceClient(conn),
	}
}

func (s *clientSearchService) SearchTweets(ctx context.Context, query string) ([]int, error) {
	resp, err := s.client.SearchTweets(ctx, &searchproto.SearchTweetsRequest{Query: query})
	if err != nil {
		return nil, err
	}
	res := make([]int, 0, len(resp.TweetIds))
	for i := range resp.TweetIds {
		res = append(res, int(resp.TweetIds[i]))
	}

	return res, nil
}

func (s *clientSearchService) SearchUsers(ctx context.Context, query string) ([]int, error) {
	resp, err := s.client.SearchUsers(ctx, &searchproto.SearchUsersRequest{Query: query})
	if err != nil {
		return nil, err
	}
	res := make([]int, 0, len(resp.UserIds))
	for i := range resp.UserIds {
		res = append(res, int(resp.UserIds[i]))
	}
	return res, nil
}
