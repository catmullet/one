package localstore

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	ttlKeyTimeFormat = "02150405"
)

type LocalIdempotencyStore struct {
	ctx      context.Context
	ttl      time.Duration
	store    map[string]struct{}
	ttlStore map[string][]string
	mu       *sync.RWMutex
}

func NewLocalIdempotencyStore(ctx context.Context, ttl time.Duration) *LocalIdempotencyStore {
	store := &LocalIdempotencyStore{
		ctx:      ctx,
		ttl:      ttl,
		store:    make(map[string]struct{}),
		ttlStore: make(map[string][]string),
		mu:       new(sync.RWMutex),
	}
	go store.cleanup()
	return store
}

func (lis *LocalIdempotencyStore) AddKey(key string) error {
	if lis.HasKey(key) != nil {
		lis.mu.Lock()
		defer lis.mu.Unlock()
		lis.store[key] = struct{}{}
		lis.addTtlKey(key)
	}
	return fmt.Errorf("key already exists")
}

func (lis *LocalIdempotencyStore) HasKey(key string) error {
	if _, ok := lis.store[key]; !ok {
		return fmt.Errorf("no key found")
	}
	return nil
}

func getTtlKey(ttl time.Duration) string {
	t := time.Now().Add(ttl)
	return t.Format(ttlKeyTimeFormat)
}

func getTimeFromKey(t string) (time.Time, error) {
	return time.Parse(ttlKeyTimeFormat, t)
}

func (lis *LocalIdempotencyStore) addTtlKey(key string) {
	ttlKey := getTtlKey(lis.ttl)
	if _, ok := lis.ttlStore[ttlKey]; ok {
		lis.ttlStore[ttlKey] = append(lis.ttlStore[ttlKey], key)
	} else {
		lis.ttlStore[ttlKey] = []string{key}
	}
}

func (lis *LocalIdempotencyStore) cleanup() {
	timer := time.NewTicker(lis.ttl)

	for {
		select {
		case <-lis.ctx.Done():
			timer.Stop()
			return
		case tc := <-timer.C:
			go lis.findAndRemoveExpired(tc)
		}
	}
}

func (lis *LocalIdempotencyStore) findAndRemoveExpired(tc time.Time) {
	lis.mu.Lock()
	defer lis.mu.Unlock()
	ct := tc.Format(ttlKeyTimeFormat)
	t, _ := time.Parse(ttlKeyTimeFormat, ct)
	for k, v := range lis.ttlStore {
		keyTime, err := getTimeFromKey(k)

		if err != nil {
			continue
		}

		if t.After(keyTime) {
			for _, storeKey := range v {
				delete(lis.store, storeKey)
			}
			delete(lis.ttlStore, k)
		}
	}
}
