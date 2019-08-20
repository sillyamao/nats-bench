package bench

import (
	"fmt"
	"math/rand"
	"time"
)

type idBuilder struct {
	rand rand.Source
}

func (b *idBuilder) NewID(prefix string) string {
	return fmt.Sprintf("%v_%v", prefix, b.rand.Int63())
}

func NewID(prefix string) string {
	return defaultIDBuilder.NewID(prefix)
}

var defaultIDBuilder = idBuilder{rand: rand.NewSource(time.Now().UnixNano())}
