package conv

import (
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func FromDomainToUserModel(user *entity.User) *models.User {
	if user == nil {
		return nil
	}

	var email, password string
	if user.Credential != nil {
		email = user.Credential.Email
		password = user.Credential.Password
	}

	return &models.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       email,
		Password:    password,
		Bio:         user.Bio,
		Gen:         user.Gen,
		CreatedAt:   user.CreatedAt,
		IsSuperuser: user.IsSuperuser,
	}
}

func FromUserModelToDomain(user *models.User) *entity.User {
	if user == nil {
		return nil
	}

	return &entity.User{
		ID:          user.ID,
		Username:    user.Username,
		Gen:         user.Gen,
		Bio:         user.Bio,
		CreatedAt:   user.CreatedAt,
		IsSuperuser: user.IsSuperuser,
		Credential: &entity.Credential{
			Email:    user.Email,
			Password: user.Password,
		},
	}
}

func FromDomainToFollowModel(follow *entity.Follow) *models.Follow {
	if follow == nil {
		return nil
	}

	return &models.Follow{
		FollowerID:  follow.FollowerID,
		FollowingID: follow.FollowingID,
		CreatedAt:   follow.CreatedAt,
	}
}

func FromFollowModelToDomain(follow *models.Follow) *entity.Follow {
	if follow == nil {
		return nil
	}

	return &entity.Follow{
		FollowerID:  follow.FollowerID,
		FollowingID: follow.FollowingID,
		CreatedAt:   follow.CreatedAt,
	}
}
