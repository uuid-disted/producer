package generator

import (
	"crypto/sha512"
	"fmt"
	"sync"
	"time"

	"github.com/uuid-disted/producer/internal/utils"
)

type SnowflakeGenerator struct {
	mu            sync.Mutex
	id            int
	epoch         time.Time
	sequence      int64
	lastTimeStamp int64
}

func New(id int, epoch time.Time) *SnowflakeGenerator {
	return &SnowflakeGenerator{
		id:    id,
		epoch: epoch,
	}
}

func (g *SnowflakeGenerator) updateSequence(t time.Time, now int64) {
	if g.lastTimeStamp == now {
		g.sequence = (g.sequence + 1) & 0xFFF
		if g.sequence == 0 {
			for now <= g.lastTimeStamp {
				now = t.UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		g.sequence = 0
	}
	g.lastTimeStamp = now
}

func (g *SnowflakeGenerator) construct(parts ...int64) string {
	var result string
	for _, part := range parts {
		result += fmt.Sprintf("%d", part)
	}
	return result
}

func (g *SnowflakeGenerator) hash(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (g *SnowflakeGenerator) Generate(t time.Time) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := t.UnixNano() / int64(time.Millisecond)
	g.updateSequence(t, now)

	random, err := utils.GenerateCryptoRandomNumber(0, 1)
	if err != nil {
		panic(fmt.Sprintf("Error generating random number: %v", err))
	}

	time.Sleep(2 * time.Millisecond)
	constructed := g.construct(now, int64(g.id), g.sequence, random)
	return g.hash(constructed)
}
