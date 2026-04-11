package conv

import (
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	tweetproto "github.com/kust1q/Zapp/main/pkg/gen/proto/tweet"
)

func FromDomainToTweetProto(tweet *entity.Tweet) *tweetproto.Tweet {
	return &tweetproto.Tweet{
		Id:            int64(tweet.ID),
		Content:       tweet.Content,
		CreatedAt:     tweet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     tweet.UpdatedAt.Format(time.RFC3339),
		ParentTweetId: int64(ptrOrZero(tweet.ParentTweetID)),
		MediaUrl:      tweet.MediaUrl,
		Author:        FromDomainToTweetAuthorTweetProto(tweet.Author),
		Counters:      FromDomainToTweetCountersTweetProto(tweet.Counters),
	}
}

func FromDomainToTweetAuthorTweetProto(user *entity.SmallUser) *tweetproto.TweetAuthor {
	return &tweetproto.TweetAuthor{
		Id:        int64(user.ID),
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToTweetCountersTweetProto(counters *entity.Counters) *tweetproto.TweetCounters {
	return &tweetproto.TweetCounters{
		ReplyCount:   int64(counters.ReplyCount),
		RetweetCount: int64(counters.RetweetCount),
		LikeCount:    int64(counters.LikeCount),
	}
}

func FromDomainToTweetListTweetProto(tweets []entity.Tweet) *tweetproto.TweetList {
	res := make([]*tweetproto.Tweet, 0, len(tweets))
	for i := range tweets {
		res = append(res, FromDomainToTweetProto(&tweets[i]))
	}
	return &tweetproto.TweetList{
		Tweets: res,
	}
}

func FromDomainToSmallTweetProto(user *entity.SmallUser) *tweetproto.Liker {
	return &tweetproto.Liker{
		Id:        int64(user.ID),
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}

func FromDomainToSmallUserListTweetProto(users []entity.SmallUser) *tweetproto.LikersList {
	res := make([]*tweetproto.Liker, 0, len(users))
	for i := range users {
		res = append(res, FromDomainToSmallTweetProto(&users[i]))
	}
	return &tweetproto.LikersList{
		Users: res,
	}
}
