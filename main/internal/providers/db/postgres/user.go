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

func (pg *PostgresDB) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	userModel := conv.FromDomainToUserModel(user)
	if userModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (username, email, password, bio, gen, created_at, is_active, is_superuser) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`, UserTable)

	var id int
	err := tx.QueryRowContext(ctx, query,
		userModel.Username, userModel.Email, userModel.Password,
		userModel.Bio, userModel.Gen, userModel.CreatedAt,
		userModel.IsActive, userModel.IsSuperuser).Scan(&id)
	if err != nil {
		return nil, err
	}
	userModel.ID = id
	return conv.FromUserModelToDomain(userModel), nil
}

func (pg *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var cachedModel *models.User
	var err error
	cachedModel, err = pg.Cache.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).WithField("email", email).Warn("Cache get failed, falling back to DB")
	} else if err == nil {
		return conv.FromUserModelToDomain(cachedModel), nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", UserTable)
	var userModel models.User
	err = pg.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	go func(model *models.User) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if CacheErr := pg.Cache.SetUser(cntx, model); CacheErr != nil {
			logrus.WithError(CacheErr).WithField("email", email).Warn("failed to set user in Cache")
		}
	}(&userModel)

	user := conv.FromUserModelToDomain(&userModel)
	return user, nil
}

func (pg *PostgresDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var cachedModel *models.User
	var err error
	cachedModel, err = pg.Cache.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).WithField("username", username).Warn("Cache get failed, falling back to DB")
	} else if err == nil {
		userEntity := conv.FromUserModelToDomain(cachedModel)
		return userEntity, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", UserTable)
	var userModel models.User
	err = pg.db.GetContext(ctx, &userModel, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	go func(model *models.User) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if CacheErr := pg.Cache.SetUser(cntx, model); CacheErr != nil {
			logrus.WithError(CacheErr).WithField("username", username).Warn("failed to set user in Cache")
		}
	}(&userModel)
	return conv.FromUserModelToDomain(&userModel), nil
}

func (pg *PostgresDB) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	var cachedModel *models.User
	var err error
	cachedModel, err = pg.Cache.GetUserByID(ctx, userID)
	if err != nil && !errors.Is(err, errs.ErrCacheKeyNotFound) {
		logrus.WithError(err).WithField("user_id", userID).Warn("Cache get failed, falling back to DB")
	} else if err == nil {
		return conv.FromUserModelToDomain(cachedModel), nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", UserTable)
	var userModel models.User
	err = pg.db.GetContext(ctx, &userModel, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	go func(model *models.User) {
		cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if CacheErr := pg.Cache.SetUser(cntx, model); CacheErr != nil {
			logrus.WithError(CacheErr).WithField("user_id", userID).Warn("failed to set user in Cache")
		}
	}(&userModel)
	return conv.FromUserModelToDomain(&userModel), nil
}

func (pg *PostgresDB) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", UserTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, password, userID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return pg.Cache.InvalidateUser(ctx, userID)
}

func (pg *PostgresDB) UpdateUserBio(ctx context.Context, userID int, bio string) error {
	query := fmt.Sprintf("UPDATE %s SET bio = $1 WHERE id = $2", UserTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, bio, userID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return pg.Cache.InvalidateUser(ctx, userID)
}

func (pg *PostgresDB) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	var err error
	exists, err = pg.Cache.ExistsByEmail(ctx, email)
	if err != nil {
		logrus.WithError(err).Warn("user exists by email check failed")
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", UserTable)
	var count int
	err = pg.db.GetContext(ctx, &count, query, email)
	return count > 0, err
}

func (pg *PostgresDB) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	var err error
	exists, err = pg.Cache.ExistsByUsername(ctx, username)
	if err != nil {
		logrus.WithError(err).Warn("user exists by username check failed")
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = $1", UserTable)
	var count int
	err = pg.db.GetContext(ctx, &count, query, username)
	return count > 0, err
}

func (pg *PostgresDB) FollowToUser(ctx context.Context, followerID, followingID int, createdAt time.Time) (*entity.Follow, error) {
	var followerExists, followingExists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", UserTable)
	var err error
	err = pg.db.QueryRowContext(ctx, checkQuery, followerID).Scan(&followerExists)
	if err != nil {
		return nil, err
	}
	if !followerExists {
		return nil, errs.ErrUserNotFound
	}

	err = pg.db.QueryRowContext(ctx, checkQuery, followingID).Scan(&followingExists)
	if err != nil {
		return nil, err
	}
	if !followingExists {
		return nil, errs.ErrUserNotFound
	}

	query := fmt.Sprintf("INSERT INTO %s (follower_id, following_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", FollowsTable)
	_, err = pg.db.ExecContext(ctx, query, followerID, followingID, createdAt)
	if err != nil {
		return nil, err
	}

	followEntity := &entity.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   createdAt,
	}
	return followEntity, nil
}

func (pg *PostgresDB) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	var followerExists, followingExists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", UserTable)
	var err error
	err = pg.db.QueryRowContext(ctx, checkQuery, followerID).Scan(&followerExists)
	if err != nil {
		return err
	}
	if !followerExists {
		return errs.ErrUserNotFound
	}

	err = pg.db.QueryRowContext(ctx, checkQuery, followingID).Scan(&followingExists)
	if err != nil {
		return err
	}
	if !followingExists {
		return errs.ErrUserNotFound
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE follower_id = $1 AND following_id = $2", FollowsTable)
	var result sql.Result
	result, err = pg.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (pg *PostgresDB) GetFollowersIds(ctx context.Context, username string, limit, offset int) ([]int, error) {
	var userExists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE username = $1)", UserTable)
	var err error
	err = pg.db.QueryRowContext(ctx, checkQuery, username).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errs.ErrUserNotFound
	}

	query := fmt.Sprintf(`
        SELECT f.follower_id
        FROM %s f
        JOIN %s u ON f.following_id = u.id
        WHERE u.username = $1
		LIMIT $2 OFFSET $3`,
		FollowsTable, UserTable)

	var res []int
	err = pg.db.SelectContext(ctx, &res, query, username, limit, offset)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pg *PostgresDB) GetFollowingsIds(ctx context.Context, username string, limit, offset int) ([]int, error) {
	var userExists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE username = $1)", UserTable)
	var err error
	err = pg.db.QueryRowContext(ctx, checkQuery, username).Scan(&userExists)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errs.ErrUserNotFound
	}

	query := fmt.Sprintf(`
        SELECT f.following_id
        FROM %s f
        JOIN %s u ON f.follower_id = u.id
        WHERE u.username = $1
		LIMIT $2 OFFSET $3`,
		FollowsTable, UserTable)

	var res []int
	err = pg.db.SelectContext(ctx, &res, query, username, limit, offset)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pg *PostgresDB) DeleteUser(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", UserTable)
	var result sql.Result
	var err error
	result, err = pg.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}

	return pg.Cache.InvalidateUser(ctx, userID)
}
