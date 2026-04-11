package conv

import (
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/response"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

// Requests
func FromUpdateBioRequestToDomain(userID int, req *request.UpdateBio) *entity.UpdateBio {
	if req == nil {
		return nil
	}

	return &entity.UpdateBio{
		UserID: userID,
		Bio:    req.Bio,
	}
}

// Responses
func FromDomainToSmallUserResponse(user *entity.SmallUser) *response.SmallUser {
	if user == nil {
		return nil
	}

	return &response.SmallUser{
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToUserResponse(user *entity.User) *response.User {
	if user == nil {
		return nil
	}

	var email string
	if user.Credential != nil {
		email = user.Credential.Email
	}

	return &response.User{
		ID:        user.ID,
		Username:  user.Username,
		Bio:       user.Bio,
		Gen:       user.Gen,
		Email:     email,
		CreatedAt: user.CreatedAt,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToUserListResponse(users []entity.User) []response.User {
	if users == nil {
		return nil
	}

	res := make([]response.User, 0, len(users))
	for _, t := range users {
		userResponse := FromDomainToUserResponse(&t)
		if userResponse != nil {
			res = append(res, *userResponse)
		}
	}
	return res
}

func FromDomainToUserProfileResponse(userProfile *entity.UserProfile) *response.UserProfile {
	if userProfile == nil {
		return nil
	}

	return &response.UserProfile{
		User:   FromDomainToUserResponse(userProfile.User),
		Tweets: FromDomainToTweetListResponse(userProfile.Tweets),
	}
}

func FromDomainToFollow(follow *entity.Follow) *response.Follow {
	if follow == nil {
		return nil
	}

	return &response.Follow{
		FollowerID:  follow.FollowerID,
		FollowingID: follow.FollowingID,
		CreatedAt:   follow.CreatedAt,
	}
}

func FromDomainToSmallUserListResponse(users []entity.SmallUser) []response.SmallUser {
	if users == nil {
		return nil
	}

	responses := make([]response.SmallUser, 0, len(users))
	for _, u := range users {
		userResponse := FromDomainToSmallUserResponse(&u)
		if userResponse != nil {
			responses = append(responses, *userResponse)
		}
	}
	return responses
}
