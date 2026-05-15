package postgres

import (
	"context"
	"fmt"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	conv "github.com/kust1q/Zapp/main/internal/providers/db/conv"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/lib/pq"
)

func (pg *PostgresDB) GetTweetsByIDs(ctx context.Context, ids []int) ([]entity.Tweet, error) {
	if len(ids) == 0 {
		return []entity.Tweet{}, nil
	}

	query := fmt.Sprintf(`
			SELECT id, user_id, parent_tweet_id, content, created_at, updated_at 
			FROM %s 
			WHERE id = ANY($1)`,
		TweetsTable)

	var tweetModels []models.Tweet

	err := pg.db.SelectContext(ctx, &tweetModels, query, pq.Array(ids))

	if err != nil {
		return nil, err
	}

	tweetMap := make(map[int]models.Tweet)
	for _, t := range tweetModels {
		tweetMap[t.ID] = t
	}

	var result []entity.Tweet
	for _, id := range ids {
		if t, ok := tweetMap[id]; ok {
			result = append(result, *conv.FromTweetModelToDomain(&t))
		}
	}

	return result, nil
}

func (pg *PostgresDB) GetUsersByIDs(ctx context.Context, ids []int) ([]entity.User, error) {
	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	query := fmt.Sprintf(`
			SELECT id, username, email, password, bio, gen, created_at, is_active, is_superuser 
			FROM %s 
			WHERE id = ANY($1)`,
		UserTable)

	var userModels []models.User
	if err := pg.db.SelectContext(ctx, &userModels, query, pq.Array(ids)); err != nil {
		return nil, err
	}

	userMap := make(map[int]models.User)
	for _, u := range userModels {
		userMap[u.ID] = u
	}

	var result []entity.User
	for _, id := range ids {
		if u, ok := userMap[id]; ok {
			result = append(result, *conv.FromUserModelToDomain(&u))
		}
	}

	return result, nil
}
