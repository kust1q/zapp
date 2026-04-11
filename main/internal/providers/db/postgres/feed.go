package postgres

import (
	"context"
	"fmt"

	conv "github.com/kust1q/Zapp/main/internal/providers/db/conv"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/lib/pq"
)

func (s *PostgresDB) GetFeedByAuthorsIds(ctx context.Context, userIDs []int, limit, offset int) ([]entity.Tweet, error) {
	query := fmt.Sprintf(`
        SELECT id, user_id, parent_tweet_id, content, created_at, updated_at
        FROM %s
        WHERE user_id = ANY($1)
		UNION ALL
		SELECT t.id, t.user_id, t.parent_tweet_id, t.content, r.created_at, t.updated_at
		FROM %s r
		JOIN %s t ON r.tweet_id = t.id
		WHERE r.user_id = ANY($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		TweetsTable, RetweetsTable, TweetsTable)

	var tweetModels []models.Tweet
	if err := s.db.SelectContext(ctx, &tweetModels, query, pq.Array(userIDs), limit, offset); err != nil {
		return nil, err
	}
	return conv.FromTweetModelToDomainList(tweetModels), nil
}

func (pg *PostgresDB) GetAllTweets(ctx context.Context, limit, offset int) ([]entity.Tweet, error) {
	query := fmt.Sprintf(`
        SELECT id, user_id, parent_tweet_id, content, created_at, updated_at 
        FROM %s
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		TweetsTable)

	var tweetModels []models.Tweet

	err := pg.db.SelectContext(ctx, &tweetModels, query, limit, offset)

	if err != nil {
		return nil, err
	}

	return conv.FromTweetModelToDomainList(tweetModels), nil
}
