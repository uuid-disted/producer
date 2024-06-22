package generator

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"math/big"
	"time"
)

type SnowflakeGenerator struct {
	ID            int
	Epoch         time.Time
	Sequence      int64
	LastTimeStamp int64
	UseBuffer     bool
	BufferChan    chan int64
}

type SnowflakeGeneratorConfig struct {
	GeneratorConfig
	ID    int
	Epoch time.Time
}

func NewSnowflakeGenerator(config SnowflakeGeneratorConfig) *SnowflakeGenerator {
	gen := &SnowflakeGenerator{
		ID:        config.ID,
		Epoch:     config.Epoch,
		UseBuffer: config.UseBuffer,
	}

	if config.UseBuffer {
		gen.BufferChan = make(chan int64, 1000)
		go func() {
			for {
				n, err := gen.random()
				if err != nil {
					continue
				}
				gen.BufferChan <- n
			}
		}()
	}

	return gen
}

func (gen *SnowflakeGenerator) updateSequence(t time.Time, now int64) {
	if gen.LastTimeStamp == now {
		gen.Sequence = (gen.Sequence + 1) & 0xFFF
		if gen.Sequence == 0 {
			for now <= gen.LastTimeStamp {
				now = t.UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		gen.Sequence = 0
	}
	gen.LastTimeStamp = now
}

func (gen *SnowflakeGenerator) construct(parts ...int64) string {
	var result string
	for _, part := range parts {
		result += fmt.Sprintf("%d", part)
	}
	return result
}

func (gen *SnowflakeGenerator) hash(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (gen *SnowflakeGenerator) random() (int64, error) {
	max := big.NewInt(1<<63 - 1)
	b, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return b.Int64(), nil
}

func (gen *SnowflakeGenerator) GetRandom() (int64, error) {
	if gen.UseBuffer {
		return <-gen.BufferChan, nil
	} else {
		return gen.random()
	}
}

func (gen *SnowflakeGenerator) Generate(t time.Time) (string, error) {
	now := t.UnixNano() / int64(time.Millisecond)
	gen.updateSequence(t, now)

	random, err := gen.random()
	if err != nil {
		return "", err
	}
	constructed := gen.construct(now, int64(gen.ID), gen.Sequence, random)
	return gen.hash(constructed), nil
}
