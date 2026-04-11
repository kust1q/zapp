package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

const (
	UserTable           = "users"
	FollowsTable        = "follows"
	LikesTable          = "likes"
	TweetsTable         = "tweets"
	RetweetsTable       = "retweets"
	TweetMediaTable     = "tweet_media"
	SecretQuestionTable = "secret_questions"
	AvatarsTable        = "avatars"
)

type PostgresDB struct {
	db    *sqlx.DB
	Cache cache
}

func NewPostgresDB(db *sqlx.DB, cache cache) *PostgresDB {
	return &PostgresDB{
		db:    db,
		Cache: cache,
	}
}

func (pg *PostgresDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return pg.db.BeginTx(ctx, nil)
}
