package conv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kust1q/Zapp/main/internal/controllers/http/conv"
	"github.com/kust1q/Zapp/main/internal/controllers/http/dto/request"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

func TestConv_Auth(t *testing.T) {
	req := &request.SignUp{Username: "test", Email: "test@test.com"}
	user := conv.FromSignUpRequestToDomain(req)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "test@test.com", user.Credential.Email)
	assert.Nil(t, conv.FromSignUpRequestToDomain(nil))

	loginReq := &request.SignIn{Email: "test@test.com"}
	cred := conv.FromSignInRequestToDomain(loginReq)
	assert.Equal(t, "test@test.com", cred.Email)
	assert.Nil(t, conv.FromSignInRequestToDomain(nil))

	refresh := conv.FromRefreshRequestToDomain("token")
	assert.Equal(t, "token", refresh.Refresh)
	assert.Nil(t, conv.FromRefreshRequestToDomain(""))

	updatePwd := &request.UpdatePassword{OldPassword: "old"}
	reset := conv.FromResetPasswordRequestToDomain(1, updatePwd)
	assert.Equal(t, 1, reset.UserID)
	assert.Equal(t, "old", reset.OldPassword)
	assert.Nil(t, conv.FromResetPasswordRequestToDomain(1, nil))

	forgot := conv.FromForgotPasswordRequestToDomain(&request.ForgotPassword{Email: "e"})
	assert.Equal(t, "e", forgot.Email)
	assert.Nil(t, conv.FromForgotPasswordRequestToDomain(nil))

	recovery := conv.FromRecoveryPasswordRequestToDomain(&request.RecoveryPassword{RecoveryToken: "t"})
	assert.Equal(t, "t", recovery.RecoveryToken)
	assert.Nil(t, conv.FromRecoveryPasswordRequestToDomain(nil))

	resp := conv.FromDomainToSignUpResponse(&entity.User{ID: 1})
	assert.Equal(t, 1, resp.ID)
	assert.Nil(t, conv.FromDomainToSignUpResponse(nil))

	tokensResp := conv.FromDomainToAccessResponse(&entity.Tokens{Access: &entity.Access{Access: "a"}})
	assert.Equal(t, "a", tokensResp.Access)
	assert.Nil(t, conv.FromDomainToAccessResponse(nil))

	recResp := conv.FromDomainToRecoveryResponse(&entity.Recovery{Recovery: "r"})
	assert.Equal(t, "r", recResp.RecoveryToken)
	assert.Nil(t, conv.FromDomainToRecoveryResponse(nil))
}

func TestConv_Media(t *testing.T) {
	m := conv.FromDomainToMediaResponse(&entity.TweetMedia{Path: "url"})
	assert.Equal(t, "url", m.MediaUrl)
	assert.Nil(t, conv.FromDomainToMediaResponse(nil))

	a := conv.FromDomainToAvatarResponse(&entity.Avatar{Path: "url"})
	assert.Equal(t, "url", a.AvatarUrl)
	assert.Nil(t, conv.FromDomainToAvatarResponse(nil))
}

func TestConv_Tweets(t *testing.T) {
	tr := conv.FromTweetRequestToDomain(1, nil, nil, &request.Tweet{Content: "c"})
	assert.Equal(t, 1, tr.Author.ID)
	assert.Equal(t, "c", tr.Content)

	tw := conv.FromDomainToTweetResponse(&entity.Tweet{ID: 1, Author: &entity.SmallUser{ID: 1}})
	assert.Equal(t, 1, tw.ID)
	assert.Nil(t, conv.FromDomainToTweetResponse(nil))

	tws := conv.FromDomainToTweetListResponse([]entity.Tweet{{ID: 1, Author: &entity.SmallUser{ID: 1}}})
	assert.Len(t, tws, 1)
}

func TestConv_User(t *testing.T) {
	ur := conv.FromDomainToUserResponse(&entity.User{ID: 1})
	assert.Equal(t, 1, ur.ID)
	assert.Nil(t, conv.FromDomainToUserResponse(nil))

	urs := conv.FromDomainToSmallUserListResponse([]entity.SmallUser{{ID: 1}})
	assert.Len(t, urs, 1)

	upr := conv.FromDomainToUserProfileResponse(&entity.UserProfile{User: &entity.User{ID: 1}})
	assert.Equal(t, 1, upr.User.ID)
	assert.Nil(t, conv.FromDomainToUserProfileResponse(nil))
}
