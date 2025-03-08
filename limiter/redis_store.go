package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type RedisStore struct {
	rdb RedisClient
}

func NewRedisStore(rdb RedisClient) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Allow(key string, limit int, duration time.Duration) (bool, error) {
	ctx := context.Background()
	count, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		s.rdb.Expire(ctx, key, duration)
	}

	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}
