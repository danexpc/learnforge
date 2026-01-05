//go:build go1.21
// +build go1.21

package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
