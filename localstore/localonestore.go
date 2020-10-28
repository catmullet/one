package localstore

import (
	"context"
	"github.com/catmullet/one"
	"sync"
	"time"
)

const (
	ttlKeyTimeFormat = "02150405"
)

type LocalOneStore struct {
	ctx      context.Context
	ttl      time.Duration
	store    map[string]struct{}
	ttlStore map[string][]string
	mu       *sync.RWMutex
}

func NewLocalOneStore(ctx context.Context, ttl time.Duration) *LocalOneStore {
	store := &LocalOneStore{
		ctx:      ctx,
		ttl:      ttl,
		store:    make(map[string]struct{}),
		ttlStore: make(map[string][]string),
		mu:       new(sync.RWMutex),
	}
	go store.cleanup()
	return store
}

func (lis *LocalOneStore) AddKey(key string) error {
	if lis.HasKey(key) == one.ErrNoKeyExist {
		lis.mu.Lock()
		defer lis.mu.Unlock()
		lis.store[key] = struct{}{}
		lis.addTtlKey(key)
		return nil
	}
	return one.ErrKeyExist
}

func (lis *LocalOneStore) HasKey(key string) error {
	if _, ok := lis.store[key]; !ok {
		return one.ErrNoKeyExist
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

func (lis *LocalOneStore) addTtlKey(key string) {
	ttlKey := getTtlKey(lis.ttl)
	if _, ok := lis.ttlStore[ttlKey]; ok {
		lis.ttlStore[ttlKey] = append(lis.ttlStore[ttlKey], key)
	} else {
		lis.ttlStore[ttlKey] = []string{key}
	}
}

func (lis *LocalOneStore) cleanup() {
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

func (lis *LocalOneStore) findAndRemoveExpired(tc time.Time) {
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
