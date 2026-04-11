package elastic

type (
	tweetDoc struct {
		Content  string `json:"content"`
		Username string `json:"username"`
		UserID   int    `json:"user_id"`
	}

	userDoc struct {
		Username string `json:"username"`
		Bio      string `json:"bio"`
	}
)
