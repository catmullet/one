package redisstore

import (
	"github.com/catmullet/one"
	"gopkg.in/redis.v5"
	"time"
)

// RedisOneStore struct to hold redis client and ttl
type RedisOneStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisOneStore returns a new RedisOneStore with redis client using gopkg.in/redis.v5
func NewRedisOneStore(options *redis.Options, ttl time.Duration) *RedisOneStore {
	return &RedisOneStore{
		client: redis.NewClient(options),
		ttl:    ttl,
	}
}

// AddKey adds key to redis if one does not exist otherwise returns error
func (rds *RedisOneStore) AddKey(key string) (ok bool, err error) {
	isExisting, err := rds.HasKey(key)
	if isExisting || err == nil {
		return
	}

	_, err = rds.client.Set(key, nil, rds.ttl).Result()
	if err != nil {
		return
	}
	return true, nil
}

// HasKey returns error if key does not exist
func (rds *RedisOneStore) HasKey(key string) (exists bool, err error) {
	_, err = rds.client.Get(key).Result()
	if err != nil {
		return false, one.ErrNoKeyExist
	}
	return true, nil
}
