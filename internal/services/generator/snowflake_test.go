package generator

import (
	"crypto/sha512"
	"encoding/hex"
	"sync"
	"testing"
	"time"
)

func TestNewSnowflakeGenerator(t *testing.T) {
	epoch := time.Now()
	generator := NewSnowflakeGenerator(1, epoch)

	if generator.ID != 1 {
		t.Errorf("Expected generator ID to be 1, got %d", generator.ID)
	}

	if !generator.Epoch.Equal(epoch) {
		t.Errorf("Expected generator epoch to be %v, got %v", epoch, generator.Epoch)
	}
}

func TestUpdateSequence(t *testing.T) {
	epoch := time.Now()
	generator := NewSnowflakeGenerator(1, epoch)
	now := epoch.UnixNano() / int64(time.Millisecond)

	// Test initial sequence
	generator.updateSequence(epoch, now)
	if generator.Sequence != 0 {
		t.Errorf("Expected sequence to be 0, got %d", generator.Sequence)
	}

	// Test sequence increment
	generator.updateSequence(epoch, now)
	if generator.Sequence != 1 {
		t.Errorf("Expected sequence to be 1, got %d", generator.Sequence)
	}

	// Test sequence reset on new timestamp
	newTime := epoch.Add(1 * time.Millisecond)
	newNow := newTime.UnixNano() / int64(time.Millisecond)
	generator.updateSequence(newTime, newNow)
	if generator.Sequence != 0 {
		t.Errorf("Expected sequence to be reset to 0, got %d", generator.Sequence)
	}
}

func TestConstruct(t *testing.T) {
	epoch := time.Now()
	generator := NewSnowflakeGenerator(1, epoch)
	constructed := generator.construct(123, 456, 789)
	expected := "123456789"
	if constructed != expected {
		t.Errorf("Expected %s, got %s", expected, constructed)
	}
}

func TestHash(t *testing.T) {
	epoch := time.Now()
	generator := NewSnowflakeGenerator(1, epoch)
	input := "test"
	expected := sha512.Sum512([]byte(input))
	expectedHex := hex.EncodeToString(expected[:])
	hashed := generator.hash(input)
	if hashed != expectedHex {
		t.Errorf("Expected %s, got %s", expectedHex, hashed)
	}
}

func TestGenerate(t *testing.T) {
	epoch := time.Now()
	generator := NewSnowflakeGenerator(1, epoch)

	generator.mu = sync.Mutex{} // Reset mutex
	generator.Sequence = 0
	generator.LastTimeStamp = 0

	now := time.Now()
	generated := generator.Generate(now)

	// Construct the expected string to hash
	nowMillis := now.UnixNano() / int64(time.Millisecond)
	constructed := generator.construct(nowMillis, int64(generator.ID), generator.Sequence, 0)
	expectedHash := sha512.Sum512([]byte(constructed))
	expectedHex := hex.EncodeToString(expectedHash[:])

	if generated != expectedHex {
		t.Errorf("Expected %s, got %s", expectedHex, generated)
	}
}
