package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *tokensDB) StoreRefresh(ctx context.Context, refreshToken, userID string, ttl time.Duration) error {
	refreshKey := s.buildRefreshKey(refreshToken)
	_, err := s.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, refreshKey, refreshToken)
		pipe.Expire(ctx, s.buildSessionKey(userID), ttl)
		pipe.Set(ctx, refreshKey, userID, ttl)
		return nil
	})
	return err
}

func (s *tokensDB) GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshKey := s.buildRefreshKey(refreshToken)
	userID, err := s.redis.Get(ctx, refreshKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis error: %w", err)
	}
	return userID, nil
}

func (s *tokensDB) RemoveRefresh(ctx context.Context, refreshToken string) error {
	userID, err := s.GetUserIdByRefreshToken(ctx, refreshToken)
	if userID == "" {
		return nil
	}
	if err != nil {
		return err
	}
	refreshKey := s.buildRefreshKey(refreshToken)
	pipe := s.redis.Pipeline()
	pipe.Del(ctx, refreshKey)
	pipe.SRem(ctx, s.buildSessionKey(userID), refreshToken)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokensDB) CloseAllSessions(ctx context.Context, userID string) error {
	tokens, err := s.redis.SMembers(ctx, s.buildSessionKey(userID)).Result()
	if err != nil {
		return err
	}
	pipe := s.redis.Pipeline()
	for _, token := range tokens {
		refreshKey := s.buildRefreshKey(token)
		pipe.Del(ctx, refreshKey)
		pipe.SRem(ctx, s.buildSessionKey(userID), token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokensDB) buildRefreshKey(refreshToken string) string {
	return prefixRefreshToken + refreshToken
}

func (s *tokensDB) buildSessionKey(userID string) string {
	return prefixUserSessions + userID
}
