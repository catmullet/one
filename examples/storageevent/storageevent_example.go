package main

import (
	"context"
	"github.com/catmullet/one"
	"github.com/catmullet/one/localstore"
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
	ctx := context.Background()

	evt := Event{
		Bucket:     "gcs://test-bucket",
		Object:     "test_data.txt",
		Version:    1,
		UpdateTime: time.Now(),
	}

	var localStore one.OneStore
	localStore = localstore.NewLocalOneStore(ctx, time.Second*15)

	log.Println("PubSub messages coming in for storage events...")

	for i := 0; i < 10; i++ {
		evt.Version = i
		oneKey := one.MakeKey(evt.Bucket, evt.Object, evt.Version)

		err := localStore.AddKey(oneKey)
		if err == one.ErrKeyExist {
			log.Println("key already exists")
		} else {
			log.Println("key was added")
		}
	}

	log.Println("Duplicate PubSub messages coming in for storage events...")

	for i := 0; i < 10; i++ {
		evt.Version = i
		oneKey := one.MakeKey(evt.Bucket, evt.Object, evt.Version)

		err := localStore.AddKey(oneKey)
		if err == one.ErrKeyExist {
			log.Println("key already exists")
		} else {
			log.Println("key was added")
		}
	}

	log.Println("finished")
}
