// +build !go1.21

package cache

import "fmt"

func NewRedisCache(redisURL string) (Cache, error) {
	return nil, fmt.Errorf("Redis support requires Go 1.21 or later. Please use in-memory cache or upgrade Go version")
}

