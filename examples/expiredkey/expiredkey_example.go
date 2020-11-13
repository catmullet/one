package main

import (
	"github.com/catmullet/one"
	"github.com/catmullet/one/redisstore"
	"gopkg.in/redis.v5"
	"log"
	"time"
)

type Event struct {
	Bucket     string
	Object     string
	Version    int
	UpdateTime time.Time
}

func main() {
	evt := Event{
		Bucket:     "gcs://test-bucket",
		Object:     "test_data.txt",
		Version:    1,
		UpdateTime: time.Now(),
	}

	var redisStore one.Store
	redisStore = redisstore.NewRedisStore(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	}, time.Second*2)

	log.Println("PubSub messages coming in for storage events...")

	for i := 0; i < 10; i++ {
		evt.Version = i
		oneKey := one.MakeKey(evt.Bucket, evt.Object, evt.Version)

		ok, err := redisStore.AddKey(oneKey)
		if err != nil {
			log.Println("error", err)
			continue
		}
		if !ok {
			log.Println("key already exists")
			continue
		}
		log.Println("key was added")
	}

	log.Println("Sleep to let ttl pass on keys...")
	time.Sleep(time.Second * 5)

	log.Println("Duplicate PubSub messages coming in for storage events...")

	for i := 0; i < 10; i++ {
		evt.Version = i
		oneKey := one.MakeKey(evt.Bucket, evt.Object, evt.Version)

		ok, err := redisStore.AddKey(oneKey)
		if err != nil {
			log.Println("error", err)
			continue
		}
		if !ok {
			log.Println("key already exists")
			continue
		}
		log.Println("key was added")
	}

	log.Println("finished")
}
