package one

import "fmt"

var (
	// ErrKeyExist Error for notifying of existing key.
	ErrKeyExist = fmt.Errorf("error: key exists")
)

// OneStore interface for existing and custom key storage.
type Store interface {
	AddKey(key string) (ok bool, err error)
}
