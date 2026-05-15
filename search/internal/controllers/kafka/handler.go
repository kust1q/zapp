package kafka

import (
	"context"
	"encoding/json"

	"github.com/kust1q/Zapp/search/internal/domain/entity"
	"github.com/kust1q/Zapp/search/internal/domain/events"
)

type eventSearchHandler struct {
	searchService searchService
}

func NewSearchHandler(service searchService) *eventSearchHandler {
	return &eventSearchHandler{
		searchService: service,
	}
}

func (h *eventSearchHandler) Handle(ctx context.Context, topic string, data []byte) error {
	switch topic {
	case events.TopicTweet:
		return h.handleTweet(ctx, data)
	case events.TopicUser:
		return h.handleUser(ctx, data)
	}
	return nil
}

func (h *eventSearchHandler) handleTweet(ctx context.Context, data []byte) error {
	var meta struct {
		EventType events.EventType `json:"event_type"`
	}
	var err error
	if err = json.Unmarshal(data, &meta); err != nil {
		return err
	}

	switch meta.EventType {
	case events.TweetCreateEvent, events.TweetUpdateEvent:
		var ev events.TweetEvent
		if err = json.Unmarshal(data, &ev); err != nil {
			return err
		}
		tweet := entity.Tweet{
			ID:      ev.ID,
			Content: ev.Content,
			Author: &entity.SmallUser{
				ID:       ev.UserID,
				Username: ev.Username,
			},
		}
		return h.searchService.IndexTweet(ctx, &tweet)

	case events.TweetDeleteEvent:
		var ev events.TweetDeleted
		if err = json.Unmarshal(data, &ev); err != nil {
			return err
		}
		return h.searchService.DeleteTweet(ctx, ev.ID)
	}
	return nil
}

func (h *eventSearchHandler) handleUser(ctx context.Context, data []byte) (err error) {
	var meta struct {
		EventType events.EventType `json:"event_type"`
	}
	if err = json.Unmarshal(data, &meta); err != nil {
		return err
	}

	switch meta.EventType {
	case events.UserCreateEvent:
		var ev events.UserEvent
		if err = json.Unmarshal(data, &ev); err != nil {
			return err
		}
		user := entity.User{
			ID:       ev.ID,
			Username: ev.Username,
			Bio:      ev.Bio,
		}
		return h.searchService.IndexUser(ctx, &user)

	case events.UserDeleteEvent:
		var ev events.UserDeleted
		if err = json.Unmarshal(data, &ev); err != nil {
			return err
		}
		return h.searchService.DeleteUserWithTweets(ctx, ev.ID)
	}


	return nil
}
