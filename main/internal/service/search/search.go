package search

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type service struct {
	db     searchStorage
	media  mediaService
	tweet  tweetService
	search searchProvider
}

func NewSearchService(db searchStorage, media mediaService, tweet tweetService, search searchProvider) *service {
	return &service{
		db:     db,
		media:  media,
		tweet:  tweet,
		search: search,
	}
}

func (s *service) SearchTweets(ctx context.Context, query string) ([]entity.Tweet, error) {
	ids, err := s.search.SearchTweets(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tweets: %w", err)
	}
	if len(ids) == 0 {
		return []entity.Tweet{}, nil
	}

	tweets, err := s.db.GetTweetsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tweets by ids: %w", err)
	}

	res := make([]entity.Tweet, 0, len(tweets))
	for _, t := range tweets {
		tr, err := s.tweet.BuildEntityTweetToResponse(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}

	return res, nil
}

func (s *service) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	ids, err := s.search.SearchUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	users, err := s.db.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by ids: %w", err)
	}

	for i := range users {
		users[i].AvatarUrl, err = s.media.GetAvatarUrlByUserID(ctx, users[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user avatar by id: %w", err)
		}
	}

	return users, nil
}
