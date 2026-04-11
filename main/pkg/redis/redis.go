package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	redis := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := redis.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redis, nil
}
