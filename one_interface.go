package one

import "fmt"

var (
	ErrKeyExist   = fmt.Errorf("Error: Key Exists")
	ErrNoKeyExist = fmt.Errorf("Error: No Key Exists")
)

type OneStore interface {
	AddKey(key string) error
	HasKey(key string) error
}
