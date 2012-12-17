package stomp

import (
	"fmt"
	"math/rand"
	"time"
)

type messageIdGenerator struct {
	counter int64
	rand    *rand.Rand
}

func newMessageIdGenerator() *messageIdGenerator {
	source := rand.NewSource(time.Now().Unix())
	mig := &messageIdGenerator{rand: rand.New(source)}

	return mig
}

func (mig *messageIdGenerator) Generate() string {
	mig.counter = mig.counter + int64(mig.rand.Int31n(64) + 1)
	return fmt.Sprintf("%d", mig.counter)
}
