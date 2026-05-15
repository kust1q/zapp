package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kust1q/Zapp/main/internal/domain/entity"
	"github.com/kust1q/Zapp/main/internal/errs"
	conv "github.com/kust1q/Zapp/main/internal/providers/db/conv"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
)

// UpsertByTweetIdTx
func (pg *PostgresDB) UpsertByTweetIdTx(ctx context.Context, tx *sql.Tx, media *entity.TweetMedia) (*entity.TweetMedia, error) {
	mediaModel := conv.FromDomainToTweetMediaModel(media)
	if mediaModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (tweet_id, path, mime_type, size_bytes)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (tweet_id) DO UPDATE
        SET path = EXCLUDED.path,
            mime_type = EXCLUDED.mime_type,
            size_bytes = EXCLUDED.size_bytes
        RETURNING id
    `, TweetMediaTable)

	var id int
	if err := tx.QueryRowContext(ctx, query, mediaModel.TweetID, mediaModel.Path, mediaModel.MimeType, mediaModel.SizeBytes).Scan(&id); err != nil {
		return nil, err
	}

	mediaModel.ID = id
	updatedMedia := conv.FromTweetMediaModelToDomain(mediaModel)

	return updatedMedia, nil
}

func (pg *PostgresDB) GetMediaPathByTweetID(ctx context.Context, tweetID int) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE tweet_id = $1", TweetMediaTable)
	var path string
	if err := pg.db.GetContext(ctx, &path, query, tweetID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	return path, nil
}

func (pg *PostgresDB) GetMediaDataByTweetID(ctx context.Context, tweetID int) (*entity.TweetMedia, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE tweet_id = $1", TweetMediaTable)
	var mediaModel models.TweetMedia
	if err := pg.db.GetContext(ctx, &mediaModel, query, tweetID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTweetMediaNotFound
		}
		return nil, err
	}

	media := conv.FromTweetMediaModelToDomain(&mediaModel)
	return media, nil
}

func (pg *PostgresDB) DeleteMediaByTweetID(ctx context.Context, tweetID, userID int) error {
	query := fmt.Sprintf(`
        DELETE FROM %s m
        WHERE m.tweet_id = $1 
        AND EXISTS (
            SELECT 1 FROM %s t 
            WHERE t.id = m.tweet_id AND t.user_id = $2)`,
		TweetMediaTable,
		TweetsTable)

	result, err := pg.db.ExecContext(ctx, query, tweetID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrTweetMediaNotFound
	}

	return nil
}

func (pg *PostgresDB) UploadAvatarTx(ctx context.Context, tx *sql.Tx, avatar *entity.Avatar) (*entity.Avatar, error) {
	avatarModel := conv.FromDomainToAvatarModel(avatar)
	if avatarModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf("INSERT INTO %s (user_id, path, mime_type, size_bytes) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO NOTHING RETURNING id", AvatarsTable)
	var id int
	if err := tx.QueryRowContext(ctx, query, avatarModel.UserID, avatarModel.Path, avatarModel.MimeType, avatarModel.SizeBytes).Scan(&id); err != nil {
		return nil, err
	}

	avatarModel.ID = id
	updatedAvatar := conv.FromAvatarModelToDomain(avatarModel)

	return updatedAvatar, nil
}

func (pg *PostgresDB) GetAvatarPathByUserID(ctx context.Context, userID int) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE user_id = $1", AvatarsTable)
	var path string
	if err := pg.db.GetContext(ctx, &path, query, userID); err != nil {
		return "", err
	}
	return path, nil
}

func (pg *PostgresDB) GetAvatarDataByUserID(ctx context.Context, userID int) (*entity.Avatar, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", AvatarsTable)
	var avatarModel models.Avatar
	if err := pg.db.GetContext(ctx, &avatarModel, query, userID); err != nil {
		return nil, err
	}

	avatar := conv.FromAvatarModelToDomain(&avatarModel)
	return avatar, nil
}

func (pg *PostgresDB) DeleteAvatarByUserID(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", AvatarsTable)
	result, err := pg.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("avatar not found")
	}

	return nil
}

func (pg *PostgresDB) GetMediaUrlsByUserID(ctx context.Context, userID int) ([]string, error) {
	query := fmt.Sprintf("SELECT tm.path FROM %s tm JOIN %s t ON tm.tweet_id = t.id WHERE t.user_id = $1", TweetMediaTable, TweetsTable)

	var urls []string
	if err := pg.db.SelectContext(ctx, &urls, query, userID); err != nil {
		return nil, err
	}
	return urls, nil
}
