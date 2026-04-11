package events

var (
	UserCreateEvent EventType = "user.created"
	UserDeleteEvent EventType = "user.deleted"
)

type UserEvent struct {
	EventType EventType `json:"event_type"`
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
}

type UserDeleted struct {
	EventType EventType `json:"event_type"`
	ID        int       `json:"id"`
}
