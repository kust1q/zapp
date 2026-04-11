package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *tokensDB) StoreRecovery(ctx context.Context, recoveryToken, userID string, ttl time.Duration) error {
	recoveryKey := s.buildRecoveryKey(recoveryToken)
	return s.redis.Set(ctx, recoveryKey, userID, ttl).Err()
}

func (s *tokensDB) GetUserIdByRecoveryToken(ctx context.Context, recoveryToken string) (string, error) {
	recoveryKey := s.buildRecoveryKey(recoveryToken)
	userID, err := s.redis.Get(ctx, recoveryKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis error: %w", err)
	}
	return userID, nil
}

func (s *tokensDB) RemoveRecovery(ctx context.Context, recoveryToken string) error {
	recoveryKey := s.buildRecoveryKey(recoveryToken)
	err := s.redis.Del(ctx, recoveryKey).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to remove recovery token: %w", err)
	}
	return nil
}

func (s *tokensDB) buildRecoveryKey(recoveryToken string) string {
	return prefixRecoveryToken + recoveryToken
}
