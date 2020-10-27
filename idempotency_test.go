package one

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestObject struct {
	Name    string
	Number  int
	Floater float64
	Message []byte
	Date    time.Time
}

var (
	objectList = []TestObject{
		TestObject{"Test Object 1", 34, 45.3, []byte("Object 1"), time.Now()},
		TestObject{"Test Object 2", 3545, 0.3, []byte("Object 2"), time.Now().Add(time.Minute * 40)},
	}
)

func TestCheckIdempotencyB(t *testing.T) {
	object1 := objectList[0]
	object2 := objectList[1]

	key1 := CreateIdempotencyKey(object1.Name, object1.Number, object1.Floater, object1.Message, object1.Date)
	key2 := CreateIdempotencyKey(object2.Name, object1.Number, object1.Floater, object2.Message, object2.Date)

	assert.NotEqual(t, key1, key2)
}
