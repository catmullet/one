package one

import "fmt"

var (
	// ErrKeyExist Error for notifying of existing key
	ErrKeyExist = fmt.Errorf("Error: Key Exists")

	// ErrNoKeyExist Error for notifying of no key exists
	ErrNoKeyExist = fmt.Errorf("Error: No Key Exists")
)

// OneStore interface for existing and custom key storage
type OneStore interface {
	AddKey(key string) (ok bool, err error)
	HasKey(key string) (exists bool, err error)
}
