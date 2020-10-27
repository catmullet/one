package one

type IdempotencyStore interface {
	AddKey(key string) error
	HasKey(key string) error
}
