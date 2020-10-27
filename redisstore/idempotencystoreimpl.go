package redisstore

import (
	"errors"
	redis "gopkg.in/redis.v5"
	"time"
)

type RedisIdempotencyStore struct {
	client *redis.Client
	ttl    time.Duration
}

type RedisField struct {
	key   string
	value string
}

func NewRedisIdempotencyStore(protocol, serverAddress string, password string, maxIdle int, idleTimeout time.Duration, ttl time.Duration) *RedisIdempotencyStore {
	return &RedisIdempotencyStore{
		client: redis.NewClient(
			&redis.Options{
				Network:     protocol,
				Addr:        serverAddress,
				Password:    password,
				PoolSize:    maxIdle,
				PoolTimeout: idleTimeout,
				IdleTimeout: idleTimeout,
				DB:          0,
			}),
		ttl: ttl,
	}
}

func (rds *RedisIdempotencyStore) AddKey(key string) error {
	err := rds.HasKey(key)
	if err == nil {
		return errors.New("key exists")
	}

	_, err = rds.client.Set(key, nil, rds.ttl).Result()
	return err
}

func (rds *RedisIdempotencyStore) HasKey(key string) (err error) {
	_, err = rds.client.Get(key).Result()
	return
}
