package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
)

type service struct {
	hub hubProvider
	db  db
}

func NewNotificationService(hub hubProvider, db db) *service {
	return &service{
		hub: hub,
		db:  db,
	}
}

func (s *service) NotifyLike(ctx context.Context, actorID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tweet, err := s.db.GetTweetById(ctx, tweetID)
	if err != nil {
		return fmt.Errorf("failed to get tweet by id: %w", err)
	}

	if tweet.Author.ID == actorID {
		return nil
	}

	actor, err := s.db.GetUserByID(ctx, actorID)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	tweetText := tweet.Content
	notification := &entity.Notification{
		ID:          uuid.New().String(),
		Type:        entity.NotificationLike,
		RecipientID: tweet.Author.ID,
		ActorID:     actorID,
		ActorName:   actor.Username,
		ActorAvatar: actor.AvatarUrl,
		TweetID:     &tweetID,
		TweetText:   &tweetText,
		Timestamp:   time.Now(),
		Read:        false,
	}

	s.hub.SendNotification(notification)
	return nil
}

func (s *service) NotifyRetweet(ctx context.Context, actorID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tweet, err := s.db.GetTweetById(ctx, tweetID)
	if err != nil {
		return fmt.Errorf("failed to get tweet by id: %w", err)
	}

	if tweet.Author.ID == actorID {
		return nil
	}

	actor, err := s.db.GetUserByID(ctx, actorID)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	tweetText := tweet.Content
	notification := &entity.Notification{
		ID:          uuid.New().String(),
		Type:        entity.NotificationRetweet,
		RecipientID: tweet.Author.ID,
		ActorID:     actorID,
		ActorName:   actor.Username,
		ActorAvatar: actor.AvatarUrl,
		TweetID:     &tweetID,
		TweetText:   &tweetText,
		Timestamp:   time.Now(),
		Read:        false,
	}

	s.hub.SendNotification(notification)
	return nil
}

func (s *service) NotifyReply(ctx context.Context, actorID, tweetID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tweet, err := s.db.GetTweetById(ctx, tweetID)
	if err != nil {
		return fmt.Errorf("failed to get tweet by id: %w", err)
	}

	if tweet.Author.ID == actorID {
		return nil
	}

	actor, err := s.db.GetUserByID(ctx, actorID)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	tweetText := tweet.Content
	notification := &entity.Notification{
		ID:          uuid.New().String(),
		Type:        entity.NotificationReply,
		RecipientID: tweet.Author.ID,
		ActorID:     actorID,
		ActorName:   actor.Username,
		ActorAvatar: actor.AvatarUrl,
		TweetID:     &tweetID,
		TweetText:   &tweetText,
		Timestamp:   time.Now(),
		Read:        false,
	}

	s.hub.SendNotification(notification)
	return nil
}

func (s *service) NotifyFollow(ctx context.Context, followerID, followingID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	follower, err := s.db.GetUserByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("failed to get user by id: %w", err)
	}

	notification := &entity.Notification{
		ID:          uuid.New().String(),
		Type:        entity.NotificationFollow,
		RecipientID: followingID,
		ActorID:     followerID,
		ActorName:   follower.Username,
		ActorAvatar: follower.AvatarUrl,
		Timestamp:   time.Now(),
		Read:        false,
	}

	s.hub.SendNotification(notification)
	return nil
}
