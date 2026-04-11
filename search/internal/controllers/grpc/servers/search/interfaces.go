package searchgrpc

import (
	"context"
)

type (
	searchService interface {
		SearchTweets(ctx context.Context, query string) ([]int, error)
		SearchUsers(ctx context.Context, query string) ([]int, error)
	}
)
