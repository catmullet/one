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
func (rds *RedisOneStore) AddKey(key string) error {
	err := rds.HasKey(key)
	if err == one.ErrKeyExist {
		return one.ErrKeyExist
	}

	_, err = rds.client.Set(key, nil, rds.ttl).Result()
	return err
}

// HasKey returns error if key does not exist
func (rds *RedisOneStore) HasKey(key string) (err error) {
	_, err = rds.client.Get(key).Result()
	if err != nil {
		return one.ErrNoKeyExist
	}
	return
}
