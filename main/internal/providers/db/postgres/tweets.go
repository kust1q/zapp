package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	conv "github.com/kust1q/Zapp/main/internal/providers/db/conv"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/sirupsen/logrus"
)

func (pg *PostgresDB) CreateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	tweetModel := conv.FromDomainToTweetModel(tweet)
	if tweetModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", TweetsTable)
	var id int
	if err := pg.db.QueryRowContext(ctx, query, tweetModel.UserID, tweetModel.ParentTweetID, tweetModel.Content, tweetModel.CreatedAt, tweetModel.UpdatedAt).Scan(&id); err != nil {
		return nil, err
	}
	tweetModel.ID = id
	go func(m *models.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if cErr := pg.Cache.SetTweet(cntx, m); cErr != nil {
			logrus.WithError(cErr).Warnf("set tweet to Cache failed")
		}

		if m.ParentTweetID != nil {
			if cErr := pg.Cache.InvalidateReplies(cntx, *m.ParentTweetID); cErr != nil {
				logrus.WithError(cErr).Warnf("invalidate Cached replies failed")
			}
			if cErr := pg.Cache.InvalidateTweetCounters(cntx, *m.ParentTweetID); cErr != nil {
				logrus.WithError(cErr).Warnf("invalidate Cached counters failed")
			}
		}
	}(tweetModel)
	createdtweet := conv.FromTweetModelToDomain(tweetModel)
	return createdtweet, nil
}

func (pg *PostgresDB) CreateTweetTx(ctx context.Context, tx *sql.Tx, tweet *entity.Tweet) (*entity.Tweet, error) {
	tweetModel := conv.FromDomainToTweetModel(tweet)
	if tweetModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, parent_tweet_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", TweetsTable)
	var id int
	err := tx.QueryRowContext(ctx, query, tweetModel.UserID, tweetModel.ParentTweetID, tweetModel.Content, tweetModel.CreatedAt, tweetModel.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}
	tweetModel.ID = id

	go func(m *models.Tweet) {
		if m.ParentTweetID != nil {
			cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if cErr := pg.Cache.InvalidateReplies(cntx, *m.ParentTweetID); cErr != nil {
				logrus.WithError(cErr).Warnf("invalidate Cached replies failed")
			}
			if cErr := pg.Cache.InvalidateTweetCounters(cntx, *m.ParentTweetID); cErr != nil {
				logrus.WithError(cErr).Warnf("invalidate Cached counters failed")
			}
		}
	}(tweetModel)

	createdtweet := conv.FromTweetModelToDomain(tweetModel)
	return createdtweet, nil
}

func (pg *PostgresDB) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	var cachedModel *models.Tweet
	var err error
	cachedModel, err = pg.Cache.GetTweet(ctx, tweetID)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).Warnf("get tweet from Cache failed")
	} else if err == nil {
		return conv.FromTweetModelToDomain(cachedModel), nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TweetsTable)
	var tweetModel models.Tweet
	err = pg.db.GetContext(ctx, &tweetModel, query, tweetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTweetNotFound
		}
		return nil, err
	}

	go func(m *models.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if cErr := pg.Cache.SetTweet(cntx, m); cErr != nil {
			logrus.WithError(cErr).Warnf("set tweet to Cache failed")
		}
	}(&tweetModel)

	return conv.FromTweetModelToDomain(&tweetModel), nil
}

func (pg *PostgresDB) UpdateTweet(ctx context.Context, tweet *entity.Tweet) (*entity.Tweet, error) {
	query := fmt.Sprintf("UPDATE %s SET content = $1, updated_at = $2 WHERE id = $3 RETURNING content, updated_at", TweetsTable)
	var content string
	var updatedAt time.Time
	err := pg.db.QueryRowContext(ctx, query, tweet.Content, tweet.UpdatedAt, tweet.ID).Scan(&content, &updatedAt)
	if err != nil {
		return nil, err
	}
	tweet.Content = content
	tweet.UpdatedAt = updatedAt

	go func(tweetID int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cErr := pg.Cache.InvalidateTweet(cntx, tweetID)
		if cErr != nil {
			logrus.WithError(cErr).Warn("invalidate tweet in Cache failed")
		}
	}(tweet.ID)
	return tweet, nil
}

func (pg *PostgresDB) DeleteTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", TweetsTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, tweetID, userID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrTweetNotFound
	}
	go func(tweetID int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cErr := pg.Cache.InvalidateTweet(cntx, tweetID)
		if cErr != nil {
			logrus.WithError(cErr).Warn("invalidate tweet in Cache failed")
		}
	}(tweetID)
	return nil
}

func (pg *PostgresDB) LikeTweet(ctx context.Context, userID, tweetID int) error {
	var exists bool
	var err error
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", TweetsTable)
	err = pg.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errs.ErrTweetNotFound
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id) VALUES ($1, $2) ON CONFLICT (user_id, tweet_id) DO NOTHING", LikesTable)
	_, err = pg.db.ExecContext(ctx, query, userID, tweetID)
	if err != nil {
		return err
	}

	go func(tweetID int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cErr := pg.Cache.InvalidateTweetLikers(cntx, tweetID)
		if cErr != nil {
			logrus.WithError(cErr).Warn("invalidate tweet in Cache failed")
		}
	}(tweetID)

	return nil
}

func (pg *PostgresDB) UnLikeTweet(ctx context.Context, userID, tweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", LikesTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, userID, tweetID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrTweetNotFound
	}

	go func(tweetID int) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cErr := pg.Cache.InvalidateTweetLikers(cntx, tweetID)
		if cErr != nil {
			logrus.WithError(cErr).Warn("invalidate tweet in Cache failed")
		}
	}(tweetID)

	return nil
}

func (pg *PostgresDB) Retweet(ctx context.Context, userID, tweetID int, createdAt time.Time) error {
	var exists bool
	var err error
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", TweetsTable)
	err = pg.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return errs.ErrTweetNotFound
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, tweet_id, created_at) VALUES ($1, $2, $3)", RetweetsTable)
	_, err = pg.db.ExecContext(ctx, query, userID, tweetID, createdAt)
	return err
}

func (pg *PostgresDB) DeleteRetweet(ctx context.Context, userID, retweetID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND tweet_id = $2", RetweetsTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, userID, retweetID)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrTweetNotFound
	}

	return nil
}

func (pg *PostgresDB) GetRepliesToTweet(ctx context.Context, parentTweetID, limit, offset int) ([]entity.Tweet, error) {
	var exists bool
	var err error
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", TweetsTable)
	err = pg.db.QueryRowContext(ctx, checkQuery, parentTweetID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.ErrTweetNotFound
	}

	var ids []int
	ids, err = pg.Cache.GetReplyIDs(ctx, parentTweetID)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).Warnf("get reply ids from Cache failed")
	} else if err == nil && len(ids) > 0 {
		var tweetsMap map[int]*models.Tweet
		tweetsMap, err = pg.Cache.MGetTweets(ctx, ids)
		if err == nil && len(tweetsMap) == len(ids) {
			var res []models.Tweet
			for _, id := range ids {
				res = append(res, *tweetsMap[id])
			}
			return conv.FromTweetModelToDomainList(res), nil
		}
	}

	query := fmt.Sprintf("SELECT id, user_id, parent_tweet_id, content, created_at, updated_at FROM %s WHERE parent_tweet_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", TweetsTable)
	var tweetModels []models.Tweet
	if err = pg.db.SelectContext(ctx, &tweetModels, query, parentTweetID, limit, offset); err != nil {
		return nil, err
	}

	go func(tweetModels []models.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var idList []int
		for _, t := range tweetModels {
			idList = append(idList, t.ID)
			cErr := pg.Cache.SetTweet(cntx, &t)
			if cErr != nil {
				logrus.WithError(cErr).Warnf("set tweet to Cache failed")
			}
		}
		cErr := pg.Cache.SetReplyIDs(cntx, parentTweetID, idList)
		if cErr != nil {
			logrus.WithError(cErr).Warnf("set reply ids to Cache failed")
		}
	}(tweetModels)

	return conv.FromTweetModelToDomainList(tweetModels), nil
}

func (pg *PostgresDB) GetTweetsAndRetweetsByUsername(ctx context.Context, username string, limit, offset int) ([]entity.Tweet, error) {
	var ids []int
	var err error
	ids, err = pg.Cache.GetUserTweetIDs(ctx, username)

	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).Warnf("get reply ids from Cache failed")
	} else if err == nil && len(ids) > 0 {
		var tweetsMap map[int]*models.Tweet
		tweetsMap, err = pg.Cache.MGetTweets(ctx, ids)
		if err == nil && len(tweetsMap) == len(ids) {
			var res []models.Tweet
			for _, id := range ids {
				res = append(res, *tweetsMap[id])
			}
			return conv.FromTweetModelToDomainList(res), nil
		}
	}

	query := fmt.Sprintf(`
        SELECT t.id, t.user_id, t.parent_tweet_id, t.content, t.created_at, t.updated_at 
        FROM %s t 
        JOIN %s u ON t.user_id = u.id 
        WHERE u.username = $1
        UNION ALL
        SELECT t.id, r.user_id, t.parent_tweet_id, t.content, r.created_at, t.updated_at
        FROM %s r
        JOIN %s t ON r.tweet_id = t.id 
        JOIN %s u ON r.user_id = u.id
        WHERE u.username = $1 
        ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		TweetsTable, UserTable, RetweetsTable, TweetsTable, UserTable)

	var tweetModels []models.Tweet
	if err = pg.db.SelectContext(ctx, &tweetModels, query, username, limit, offset); err != nil {
		return nil, err
	}

	go func(list []models.Tweet) {
		cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		idList := make([]int, 0, len(list))
		for _, t := range list {
			idList = append(idList, t.ID)
			if cErr := pg.Cache.SetTweet(cntx, &t); cErr != nil {
				logrus.WithError(cErr).Warnf("set tweet to Cache failed")
			}
		}
		if cErr := pg.Cache.SetUserTweetIDs(cntx, username, idList); cErr != nil {
			logrus.WithError(cErr).Warnf("set user tweets ids to Cache failed")
		}
	}(tweetModels)

	return conv.FromTweetModelToDomainList(tweetModels), nil
}

func (pg *PostgresDB) GetCounts(ctx context.Context, tweetID int) (*entity.Counters, error) {
	var Cached *models.Counters
	var err error
	Cached, err = pg.Cache.GetTweetCounters(ctx, tweetID)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).Warn("get tweet counters from Cache failed")
	} else if err == nil {
		return conv.FromCountersModelToDomain(Cached), nil
	}
	query := fmt.Sprintf(`
        SELECT 
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as likes,
            (SELECT COUNT(*) FROM %s WHERE tweet_id = $1) as retweets,
            (SELECT COUNT(*) FROM %s WHERE parent_tweet_id = $1) as replies
    `, LikesTable, RetweetsTable, TweetsTable)

	var likes, retweets, replies int
	err = pg.db.QueryRowContext(ctx, query, tweetID).Scan(&likes, &retweets, &replies)
	if err != nil {
		return nil, err
	}
	countersModel := &models.Counters{
		LikeCount:    likes,
		RetweetCount: retweets,
		ReplyCount:   replies,
	}

	go func(id int, m *models.Counters) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if cErr := pg.Cache.SetTweetCounters(cntx, id, m); cErr != nil {
			logrus.WithError(cErr).Warn("set tweet counters to Cache failed")
		}
	}(tweetID, countersModel)

	return conv.FromCountersModelToDomain(countersModel), nil
}

func (pg *PostgresDB) GetLikes(ctx context.Context, tweetID, limit, offset int) ([]entity.SmallUser, error) {
	var exists bool
	var err error
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", TweetsTable)
	err = pg.db.QueryRowContext(ctx, checkQuery, tweetID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.ErrTweetNotFound
	}

	query := fmt.Sprintf(
		`SELECT u.id, u.username, a.path as avatar_path
        FROM %s l
        JOIN %s u ON l.user_id = u.id
        JOIN %s a ON u.id = a.user_id
        WHERE l.tweet_id = $1
		LIMIT $2 OFFSET $3`,
		LikesTable, UserTable, AvatarsTable)

	type LikerRow struct {
		ID         int    `db:"id"`
		Username   string `db:"username"`
		AvatarPath string `db:"avatar_path"`
	}

	var rows []LikerRow
	if err = pg.db.SelectContext(ctx, &rows, query, tweetID, limit, offset); err != nil {
		return nil, err
	}

	res := make([]entity.SmallUser, 0, len(rows))
	ids := make([]int, 0, len(rows))

	for _, row := range rows {
		res = append(res, entity.SmallUser{
			ID:        row.ID,
			Username:  row.Username,
			AvatarUrl: row.AvatarPath,
		})
		ids = append(ids, row.ID)
	}

	go func() {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if cErr := pg.Cache.SetTweetLikerIDs(cntx, tweetID, ids); cErr != nil {
			logrus.WithError(cErr).Warn("set tweet likers to cache failed")
		}
	}()

	return res, nil
}
