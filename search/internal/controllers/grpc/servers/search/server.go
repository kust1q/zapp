package searchgrpc

import (
	"context"

	"github.com/kust1q/Zapp/search/internal/controllers/grpc/conv"
	searchproto "github.com/kust1q/Zapp/search/pkg/gen/proto/search"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type searchServiceAPI struct {
	searchproto.UnimplementedSearchServiceServer
	searchService searchService
}

func NewSearchServer(searchService searchService) *searchServiceAPI {
	return &searchServiceAPI{
		searchService: searchService,
	}
}

func (s *searchServiceAPI) SearchUsers(ctx context.Context, req *searchproto.SearchUsersRequest) (*searchproto.SearchUsersResponse, error) {
	users, err := s.searchService.SearchUsers(ctx, req.Query)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"query":   req.Query,
			"service": "search",
		}).Error("failed to search users")
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.ToSearchUserProtoResponse(users), nil
}

func (s *searchServiceAPI) SearchTweets(ctx context.Context, req *searchproto.SearchTweetsRequest) (*searchproto.SearchTweetsResponse, error) {
	tweets, err := s.searchService.SearchTweets(ctx, req.Query)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"query":   req.Query,
			"service": "search",
		}).Error("failed to search tweets")
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.ToSearchTweetProtoResponse(tweets), nil
}
