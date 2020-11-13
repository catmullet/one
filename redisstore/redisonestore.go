package redisstore

import (
	"github.com/Kochava/collective2-ingestion-file-watcher/pkg/idempotency"
	"gopkg.in/redis.v5"
	"time"
)

// RedisStore struct to hold redis client and ttl.
type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStore returns a new RedisStore with redis client using https://gopkg.in/redis.v5.
func NewRedisStore(options *redis.Options, ttl time.Duration) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(options),
		ttl:    ttl,
	}
}

// AddKey adds key to redis if one does not exist otherwise returns error.
func (rds *RedisStore) AddKey(key string) (bool, error) {
	incr, err := rds.client.Incr(key).Result()
	if err != nil {
		return false, err
	}
	if incr > 1 {
		return false, idempotency.ErrKeyExists()
	}

	_, err = rds.client.Set(key, int64(1), rds.ttl).Result()

	return err == nil, err
}
