package conv

import (
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	userproto "github.com/kust1q/Zapp/main/pkg/gen/proto/user"
)

func FromDomainToUserProto(user *entity.User) *userproto.User {
	return &userproto.User{
		Id:          int64(user.ID),
		Username:    user.Username,
		Bio:         user.Bio,
		Gen:         user.Gen,
		Email:       user.Credential.Email,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		AvatarUrl:   user.AvatarUrl,
		IsSuperuser: user.IsSuperuser,
		IsActive:    user.IsActive,
	}
}

func FromDomainToUserProfileProto(profile *entity.UserProfile) *userproto.UserProfile {
	return &userproto.UserProfile{
		User:   FromDomainToUserProto(profile.User),
		Tweets: FromDomainToTweetListUserProto(profile.Tweets),
	}
}

func FromDomainToSmallUserProto(user *entity.SmallUser) *userproto.SmallUser {
	return &userproto.SmallUser{
		Id:        int64(user.ID),
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToSmallUserListUserProto(users []entity.SmallUser) *userproto.SmallUserList {
	res := make([]*userproto.SmallUser, 0, len(users))
	for i := range users {
		res = append(res, FromDomainToSmallUserProto(&users[i]))
	}
	return &userproto.SmallUserList{
		Users: res,
	}
}

func FromDomainToTweetUserProto(tweet *entity.Tweet) *userproto.Tweet {
	return &userproto.Tweet{
		Id:            int64(tweet.ID),
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     tweet.UpdatedAt.Format(time.RFC3339),
		ParentTweetId: int64(ptrOrZero(tweet.ParentTweetID)),
		MediaUrl:      tweet.MediaUrl,
		Author:        FromDomainToTweetAuthorUserProto(tweet.Author),
		Counters:      FromDomainToTweetCountersUserProto(tweet.Counters),
	}
}

func FromDomainToTweetListUserProto(tweets []entity.Tweet) []*userproto.Tweet {
	res := make([]*userproto.Tweet, 0, len(tweets))
	for i := range tweets {
		res = append(res, FromDomainToTweetUserProto(&tweets[i]))
	}
	return res
}

func FromDomainToTweetAuthorUserProto(user *entity.SmallUser) *userproto.TweetAuthor {
	return &userproto.TweetAuthor{
		Id:        int64(user.ID),
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToTweetCountersUserProto(counters *entity.Counters) *userproto.TweetCounters {
	return &userproto.TweetCounters{
		ReplyCount:   int64(counters.ReplyCount),
		RetweetCount: int64(counters.RetweetCount),
		LikeCount:    int64(counters.LikeCount),
	}
}
