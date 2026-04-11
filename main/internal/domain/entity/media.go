package entity

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeGIF   MediaType = "gif"
)

type (
	MediaPolicy struct {
		MaxSize       int64
		AllowedMime   []string
		AllowedExt    []string
		ForceMimeType string
	}

	TweetMedia struct {
		ID        int
		TweetID   int
		Path      string
		MimeType  string
		SizeBytes int64
	}

	Avatar struct {
		ID        int
		UserID    int
		Path      string
		MimeType  string
		SizeBytes int64
	}
)
