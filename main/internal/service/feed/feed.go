package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type service struct {
	db           db
	tweetService tweetService
}

func NewFeedService(db db, tweetService tweetService) *service {
	return &service{
		db:           db,
		tweetService: tweetService,
	}
}

func (s *service) GetUserFeedByUserId(ctx context.Context, userID, limit, offset int) ([]entity.Tweet, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	followings, err := s.db.GetFollowingsIds(ctx, user.Username, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followings by username: %w", err)
	}

	feed, err := s.db.GetFeedByAuthorsIds(ctx, followings, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed by ids: %w", err)
	}

	dflt, err := s.db.GetAllTweets(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get default feed: %w", err)
	}

	feed = append(feed, dflt...)
	res := make([]entity.Tweet, 0, len(feed))
	for _, t := range feed {
		tr, err := s.tweetService.BuildEntityTweetToResponse(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}

	return res, nil
}

func (s *service) GetDeafultFeed(ctx context.Context, limit, offset int) ([]entity.Tweet, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	feed, err := s.db.GetAllTweets(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get default feed: %w", err)
	}
	res := make([]entity.Tweet, 0, len(feed))
	for _, t := range feed {
		tr, err := s.tweetService.BuildEntityTweetToResponse(ctx, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to change tweet entity to response: %w", err)
		}
		res = append(res, *tr)
	}
	return res, nil
}
