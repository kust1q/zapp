package response

type TweetMedia struct {
	MediaUrl  string `json:"media_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type Avatar struct {
	AvatarUrl string `json:"avatar_url"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
