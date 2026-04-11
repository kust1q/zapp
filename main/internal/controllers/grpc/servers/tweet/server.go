package tweetgrpc

import (
	"context"
	"errors"

	"github.com/kust1q/Zapp/main/internal/controllers/grpc/conv"
	"github.com/kust1q/Zapp/main/internal/errs"
	tweetproto "github.com/kust1q/Zapp/main/pkg/gen/proto/tweet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tweetServerAPI struct {
	tweetproto.UnimplementedTweetServiceServer
	tweetService tweetService
}

func NewTweetServer(tweetService tweetService) *tweetServerAPI {
	return &tweetServerAPI{
		tweetService: tweetService,
	}
}

func (s *tweetServerAPI) GetTweetById(ctx context.Context, req *tweetproto.GetTweetByIdRequest) (*tweetproto.Tweet, error) {
	tweet, err := s.tweetService.GetTweetById(ctx, int(req.TweetId))
	if err != nil {
		if !errors.Is(err, errs.ErrTweetNotFound) {
			return nil, status.Error(codes.NotFound, "tweet not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToTweetProto(tweet), nil
}

func (s *tweetServerAPI) GetRepliesToTweet(ctx context.Context, req *tweetproto.GetRepliesToTweetRequest) (*tweetproto.TweetList, error) {
	replies, err := s.tweetService.GetRepliesToTweet(ctx, int(req.TweetId), int(req.Limit), int(req.Offset))
	if err != nil {
		if !errors.Is(err, errs.ErrTweetNotFound) {
			return nil, status.Error(codes.NotFound, "tweet not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToTweetListTweetProto(replies), nil
}

func (s *tweetServerAPI) GetTweetsAndRetweetsByUsername(ctx context.Context, req *tweetproto.GetTweetsAndRetweetsByUsernameRequest) (*tweetproto.TweetList, error) {
	tweets, err := s.tweetService.GetTweetsAndRetweetsByUsername(ctx, req.Username, int(req.Limit), int(req.Offset))
	if err != nil {
		if !errors.Is(err, errs.ErrTweetNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToTweetListTweetProto(tweets), nil
}

func (s *tweetServerAPI) GetLikes(ctx context.Context, req *tweetproto.GetTweetLikesRequest) (*tweetproto.LikersList, error) {
	likers, err := s.tweetService.GetLikes(ctx, int(req.TweetId), int(req.Limit), int(req.Offset))
	if err != nil {
		if !errors.Is(err, errs.ErrTweetNotFound) {
			return nil, status.Error(codes.NotFound, "tweet not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToSmallUserListTweetProto(likers), nil
}
