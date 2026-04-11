package usergrpc

import (
	"context"
	"errors"

	"github.com/kust1q/Zapp/main/internal/controllers/grpc/conv"
	"github.com/kust1q/Zapp/main/internal/errs"
	userproto "github.com/kust1q/Zapp/main/pkg/gen/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServerAPI struct {
	userproto.UnimplementedUserServiceServer
	userService userService
}

func NewUserServer(userService userService) *userServerAPI {
	return &userServerAPI{
		userService: userService,
	}
}

func (s *userServerAPI) GetUserByID(ctx context.Context, req *userproto.GetUserByIDRequest) (*userproto.User, error) {
	user, err := s.userService.GetUserByID(ctx, int(req.UserId))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return conv.FromDomainToUserProto(user), nil
}

func (s *userServerAPI) GetUserByUsername(ctx context.Context, req *userproto.GetUserByUsernameRequest) (*userproto.User, error) {
	user, err := s.userService.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return conv.FromDomainToUserProto(user), nil
}

func (s *userServerAPI) GetUserProfile(ctx context.Context, req *userproto.GetUserProfileRequest) (*userproto.UserProfile, error) {
	profile, err := s.userService.GetUserProfile(ctx, req.Username, int(req.Limit), int(req.Offset))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToUserProfileProto(profile), nil
}

func (s *userServerAPI) GetFollowers(ctx context.Context, req *userproto.GetFollowersRequest) (*userproto.SmallUserList, error) {
	users, err := s.userService.GetFollowers(ctx, req.Username, int(req.Limit), int(req.Offset))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToSmallUserListUserProto(users), nil
}

func (s *userServerAPI) GetFollowings(ctx context.Context, req *userproto.GetFollowingsRequest) (*userproto.SmallUserList, error) {
	users, err := s.userService.GetFollowings(ctx, req.Username, int(req.Limit), int(req.Offset))
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return conv.FromDomainToSmallUserListUserProto(users), nil
}
