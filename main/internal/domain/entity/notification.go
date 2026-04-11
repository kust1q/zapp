package entity

import "time"

type NotificationType string

const (
	NotificationLike    NotificationType = "like"
	NotificationRetweet NotificationType = "retweet"
	NotificationReply   NotificationType = "reply"
	NotificationFollow  NotificationType = "follow"
)

type Notification struct {
	ID          string           `json:"id"`
	Type        NotificationType `json:"type"`
	RecipientID int              `json:"recipient_id"`
	ActorID     int              `json:"actor_id"`
	ActorName   string           `json:"actor_name"`
	ActorAvatar string           `json:"actor_avatar,omitempty"`
	TweetID     *int             `json:"tweet_id,omitempty"`
	TweetText   *string          `json:"tweet_text,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
	Read        bool             `json:"read"`
}
