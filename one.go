package one

import (
	"crypto/sha256"
	"fmt"
)

func MakeKey(fields ...interface{}) (key string) {
	var sha = sha256.New()
	sha.Write([]byte(fmt.Sprintf("%v", fields...)))
	key = fmt.Sprintf("%x", sha.Sum(nil))
	return
}
