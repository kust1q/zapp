package search

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/search/internal/domain/entity"
)

type searchService struct {
	searchRepo searchRepository
}

func NewSearchService(searchRepo searchRepository) *searchService {
	return &searchService{
		searchRepo: searchRepo,
	}
}

func (s *searchService) SearchTweets(ctx context.Context, query string) ([]int, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ids, err := s.searchRepo.SearchTweets(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tweets: %w", err)
	}
	return ids, nil
}

func (s *searchService) SearchUsers(ctx context.Context, query string) ([]int, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ids, err := s.searchRepo.SearchUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return ids, nil
}

func (s *searchService) IndexTweet(ctx context.Context, tweet *entity.Tweet) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.searchRepo.IndexTweet(ctx, tweet)
}

func (s *searchService) IndexUser(ctx context.Context, user *entity.User) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.searchRepo.IndexUser(ctx, user)
}

func (s *searchService) DeleteTweet(ctx context.Context, tweetID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.searchRepo.DeleteTweet(ctx, tweetID)
}

func (s *searchService) DeleteUserWithTweets(ctx context.Context, userID int) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.searchRepo.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return s.searchRepo.DeleteTweetsByUserID(ctx, userID)
}
