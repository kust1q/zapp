package conv

import (
	"time"

	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/response"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

// Requests
func FromSignUpRequestToDomain(req *request.SignUp) *entity.User {
	if req == nil {
		return nil
	}

	return &entity.User{
		Username:    req.Username,
		Gen:         req.Gen,
		Bio:         req.Bio,
		CreatedAt:   time.Now(),
		IsSuperuser: false,
		IsActive:    true,
		Credential: &entity.Credential{
			Email:    req.Email,
			Password: req.Password,
		},
	}
}

func FromSignInRequestToDomain(req *request.SignIn) *entity.Credential {
	if req == nil {
		return nil
	}

	return &entity.Credential{
		Email:    req.Email,
		Password: req.Password,
	}
}

func FromRefreshRequestToDomain(refreshToken string) *entity.Refresh {
	if refreshToken == "" {
		return nil
	}

	return &entity.Refresh{
		Refresh: refreshToken,
	}
}

func FromResetPasswordRequestToDomain(userID int, req *request.UpdatePassword) *entity.UpdatePassword {
	if req == nil {
		return nil
	}

	return &entity.UpdatePassword{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}

func FromForgotPasswordRequestToDomain(req *request.ForgotPassword) *entity.ForgotPassword {
	if req == nil {
		return nil
	}

	return &entity.ForgotPassword{
		Email: req.Email,
	}
}

func FromRecoveryPasswordRequestToDomain(req *request.RecoveryPassword) *entity.RecoveryPassword {
	if req == nil {
		return nil
	}

	return &entity.RecoveryPassword{
		RecoveryToken: req.RecoveryToken,
		NewPassword:   req.NewPassword,
	}
}

// Responses
func FromDomainToSignUpResponse(user *entity.User) *response.SignUp {
	if user == nil {
		return nil
	}

	var email string
	if user.Credential != nil {
		email = user.Credential.Email
	}

	return &response.SignUp{
		ID:        user.ID,
		Username:  user.Username,
		Email:     email,
		Bio:       user.Bio,
		Gen:       user.Gen,
		AvatarUrl: user.AvatarUrl,
		CreatedAt: user.CreatedAt,
	}
}

func FromDomainToAccessResponse(tokens *entity.Tokens) *response.Access {
	if tokens == nil || tokens.Access == nil {
		return nil
	}

	return &response.Access{
		Access: tokens.Access.Access,
	}
}

func FromDomainToRecoveryResponse(recovery *entity.Recovery) *response.Recovery {
	if recovery == nil {
		return nil
	}

	return &response.Recovery{
		RecoveryToken: recovery.Recovery,
	}
}
