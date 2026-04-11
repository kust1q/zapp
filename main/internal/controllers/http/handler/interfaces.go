package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type (
	authService interface {
		SignUp(ctx context.Context, req *entity.User) (*entity.User, error)
		SignIn(ctx context.Context, req *entity.Credential) (*entity.Tokens, error)
		Refresh(ctx context.Context, req *entity.Refresh) (*entity.Tokens, error)
		SignOut(ctx context.Context, req *entity.Refresh) error
		VerifyAccessToken(tokenString string) (int, error)
		UpdatePassword(ctx context.Context, req *entity.UpdatePassword) error
		ForgotPassword(ctx context.Context, req *entity.ForgotPassword) (*entity.Recovery, error)
		RecoveryPassword(ctx context.Context, req *entity.RecoveryPassword) error
		GetRefreshTTL() time.Duration
	}

	tweetService interface {
		CreateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error)
		GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error)
		UpdateTweet(ctx context.Context, req *entity.Tweet) (*entity.Tweet, error)
		LikeTweet(ctx context.Context, userID, tweetID int) error
		UnlikeTweet(ctx context.Context, userID, tweetID int) error
		GetRepliesToTweet(ctx context.Context, tweetID, limit, offset int) ([]entity.Tweet, error)
		CreateRetweet(ctx context.Context, userID, tweetID int) error
		DeleteRetweet(ctx context.Context, userID, retweetID int) error
		GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error)
		GetLikes(ctx context.Context, tweetID, limit, offset int) ([]entity.SmallUser, error)
		DeleteTweet(ctx context.Context, userID, tweetID int) error
	}

	userService interface {
		Update(ctx context.Context, req *entity.UpdateBio) error
		FollowToUser(ctx context.Context, followerID, followingID int) (*entity.Follow, error)
		UnfollowUser(ctx context.Context, followerID, followingID int) error
		GetFollowers(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
		GetFollowings(ctx context.Context, username string, limit, offset int) ([]entity.SmallUser, error)
		GetUserProfile(ctx context.Context, username string, limit, offset int) (*entity.UserProfile, error)
		GetMe(ctx context.Context, userID, limit, offset int) (*entity.UserProfile, error)
		DeleteUser(ctx context.Context, userID int) error
	}

	clientSearchService interface {
		SearchTweets(ctx context.Context, query string) ([]entity.Tweet, error)
		SearchUsers(ctx context.Context, query string) ([]entity.User, error)
	}

	feedService interface {
		GetUserFeedByUserId(ctx context.Context, userID, limit, offset int) ([]entity.Tweet, error)
		GetDeafultFeed(ctx context.Context, limit, offset int) ([]entity.Tweet, error)
	}

	mediaService interface {
		GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error)
		GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error)
		DeleteTweetMedia(ctx context.Context, tweetID, userID int) error
	}

	webSocketService interface {
		HandleConnection(w http.ResponseWriter, r *http.Request, userID int) error
	}

	notificationService interface {
		NotifyLike(ctx context.Context, actorID, tweetID int) error
		NotifyRetweet(ctx context.Context, actorID, tweetID int) error
		NotifyReply(ctx context.Context, actorID, tweetID int) error
		NotifyFollow(ctx context.Context, followerID, followingID int) error
	}
)
