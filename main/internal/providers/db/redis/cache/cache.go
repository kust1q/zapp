package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/main/internal/errs"
	"github.com/kust1q/Zapp/main/internal/providers/db/models"
	"github.com/redis/go-redis/v9"
)

const (
	tweetIdCachePrefix    = "tweet:id:"
	userTweetsCachePrefix = "user:tweets:"
	repliesCachePrefix    = "tweet:replies:"
	likesCachePrefis      = "tweet:likes:"
	countersCachePrefix   = "tweet:counters:"

	userIdCachePrefix   = "user:id:"
	usernameCachePrefix = "user:username:"
	userEmaiCachePrefix = "user:email:"
)

type cache struct {
	client      *redis.Client
	defaultTtl  time.Duration
	countersTtl time.Duration
}

func NewCache(client *redis.Client, defaultTtl, countersTtl time.Duration) *cache {
	return &cache{
		client:      client,
		defaultTtl:  defaultTtl,
		countersTtl: countersTtl,
	}
}

func (c *cache) SetTweet(ctx context.Context, tweet *models.Tweet) error {
	data, err := json.Marshal(tweet)
	if err != nil {
		return fmt.Errorf("marshal tweet error: %w", err)
	}
	key := fmt.Sprintf("%s%d", tweetIdCachePrefix, tweet.ID)
	return c.client.Set(ctx, key, data, c.defaultTtl).Err()
}

func (c *cache) GetTweet(ctx context.Context, tweetID int) (*models.Tweet, error) {
	key := fmt.Sprintf("%s%d", tweetIdCachePrefix, tweetID)
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, err
	}
	var tweet models.Tweet
	if err := json.Unmarshal(val, &tweet); err != nil {
		return nil, fmt.Errorf("unmarshal tweet: %w", err)
	}
	return &tweet, nil
}

func (c *cache) MGetTweets(ctx context.Context, tweetIDs []int) (map[int]*models.Tweet, error) {
	if len(tweetIDs) == 0 {
		return nil, nil
	}

	keys := make([]string, len(tweetIDs))
	for i, id := range tweetIDs {
		keys[i] = fmt.Sprintf("%s%d", tweetIdCachePrefix, id)
	}

	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis mget error: %w", err)
	}

	result := make(map[int]*models.Tweet)

	for i, val := range values {
		if val == nil {
			continue
		}

		strVal, ok := val.(string)
		if !ok {
			continue
		}

		var tweet models.Tweet
		if err := json.Unmarshal([]byte(strVal), &tweet); err != nil {
			continue
		}
		result[tweetIDs[i]] = &tweet
	}

	return result, nil
}

func (c *cache) InvalidateTweet(ctx context.Context, tweetID int) error {
	return c.client.Del(ctx, fmt.Sprintf("%s%d", tweetIdCachePrefix, tweetID)).Err()
}

func (c *cache) setIntList(ctx context.Context, key string, ids []int) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.defaultTtl).Err()
}

func (c *cache) getIntList(ctx context.Context, key string) ([]int, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal(val, &ids); err != nil {
		return nil, fmt.Errorf("unmarshal list error: %w", err)
	}
	return ids, nil
}

func (c *cache) SetUserTweetIDs(ctx context.Context, username string, ids []int) error {
	return c.setIntList(ctx, fmt.Sprintf("%s%s", userTweetsCachePrefix, username), ids)
}

func (c *cache) GetUserTweetIDs(ctx context.Context, username string) ([]int, error) {
	return c.getIntList(ctx, fmt.Sprintf("%s%s", userTweetsCachePrefix, username))
}

func (c *cache) InvalidateUserTweets(ctx context.Context, username string) error {
	return c.client.Del(ctx, fmt.Sprintf("%s%s", userTweetsCachePrefix, username)).Err()
}

func (c *cache) SetReplyIDs(ctx context.Context, parentTweetID int, ids []int) error {
	return c.setIntList(ctx, fmt.Sprintf("%s%d", repliesCachePrefix, parentTweetID), ids)
}

func (c *cache) GetReplyIDs(ctx context.Context, parentTweetID int) ([]int, error) {
	return c.getIntList(ctx, fmt.Sprintf("%s%d", repliesCachePrefix, parentTweetID))
}

func (c *cache) InvalidateReplies(ctx context.Context, parentTweetID int) error {
	return c.client.Del(ctx, fmt.Sprintf("%s%d", repliesCachePrefix, parentTweetID)).Err()
}

func (c *cache) SetTweetLikerIDs(ctx context.Context, tweetID int, userIDs []int) error {
	return c.setIntList(ctx, fmt.Sprintf("%s%d", likesCachePrefis, tweetID), userIDs)
}

func (c *cache) GetTweetLikerIDs(ctx context.Context, tweetID int) ([]int, error) {
	return c.getIntList(ctx, fmt.Sprintf("%s%d", likesCachePrefis, tweetID))
}

func (c *cache) InvalidateTweetLikers(ctx context.Context, tweetID int) error {
	return c.client.Del(ctx, fmt.Sprintf("%s%d", likesCachePrefis, tweetID)).Err()
}

func (c *cache) SetTweetCounters(ctx context.Context, tweetID int, counters *models.Counters) error {
	data, err := json.Marshal(counters)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s%d", countersCachePrefix, tweetID)
	return c.client.Set(ctx, key, data, c.countersTtl).Err()
}

func (c *cache) GetTweetCounters(ctx context.Context, tweetID int) (*models.Counters, error) {
	key := fmt.Sprintf("%s%d", countersCachePrefix, tweetID)
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, err
	}
	var counters models.Counters
	if err := json.Unmarshal(val, &counters); err != nil {
		return nil, fmt.Errorf("unmarshal counters error: %w", err)
	}
	return &counters, nil
}

func (c *cache) InvalidateTweetCounters(ctx context.Context, tweetID int) error {
	return c.client.Del(ctx, fmt.Sprintf("%s%d", countersCachePrefix, tweetID)).Err()
}

func (c *cache) SetUser(ctx context.Context, user *models.User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}
	idKey := fmt.Sprintf("%s%d", userIdCachePrefix, user.ID)
	usernameKey := fmt.Sprintf("%s%s", usernameCachePrefix, user.Username)
	emailKey := fmt.Sprintf("%s%s", userEmaiCachePrefix, user.Email)

	pipe := c.client.Pipeline()
	pipe.Set(ctx, idKey, userData, c.defaultTtl)
	pipe.Set(ctx, usernameKey, user.ID, c.defaultTtl)
	pipe.Set(ctx, emailKey, user.ID, c.defaultTtl)

	_, err = pipe.Exec(ctx)
	return err
}

func (c *cache) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	key := fmt.Sprintf("%s%d", userIdCachePrefix, id)
	return c.fetchUser(ctx, key)
}

func (c *cache) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	indexKey := fmt.Sprintf("%s%s", usernameCachePrefix, username)

	userID, err := c.client.Get(ctx, indexKey).Int()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, fmt.Errorf("failed to get id by username: %w", err)
	}

	return c.GetUserByID(ctx, userID)
}

func (c *cache) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	indexKey := fmt.Sprintf("%s%s", userEmaiCachePrefix, email)

	userID, err := c.client.Get(ctx, indexKey).Int()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, fmt.Errorf("failed to get id by email: %w", err)
	}

	return c.GetUserByID(ctx, userID)
}

func (c *cache) InvalidateUser(ctx context.Context, userID int) error {
	pipe := c.client.Pipeline()

	pipe.Del(ctx, fmt.Sprintf("%s%d", userIdCachePrefix, userID))

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to invalidate user cache: %w", err)
	}
	return nil
}

func (c *cache) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	indexKey := fmt.Sprintf("%s%s", usernameCachePrefix, username)

	count, err := c.client.Exists(ctx, indexKey).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return count > 0, nil
}

func (c *cache) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	indexKey := fmt.Sprintf("%s%s", userEmaiCachePrefix, email)

	count, err := c.client.Exists(ctx, indexKey).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}

	return count > 0, nil
}

func (c *cache) fetchUser(ctx context.Context, key string) (*models.User, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrCacheKeyNotFound
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var user models.User
	if err := json.Unmarshal(val, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data for key %s: %w", key, err)
	}
	return &user, nil
}
