package events

var (
	TweetCreateEvent EventType = "tweet.created"
	TweetUpdateEvent EventType = "tweet.updated"
	TweetDeleteEvent EventType = "tweet.deleted"
)

type (
	TweetEvent struct {
		EventType EventType `json:"event_type"`
		ID        int       `json:"id"`
		Content   string    `json:"content"`
		UserID    int       `json:"user_id"`
		Username  string    `json:"username"`
	}

	TweetDeleted struct {
		EventType EventType `json:"event_type"`
		ID        int       `json:"id"`
	}
)
