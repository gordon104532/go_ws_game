package pkg

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewRedisConn(options *redis.Options) (cache *redis.Client, err error) {
	cache = redis.NewClient(options)
	err = cache.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return cache, nil
}
