package utils

import (
	"fmt"
	"testing"
)

func TestGenerateCryptoRandomNumber(t *testing.T) {
	tests := []struct {
		min         int64
		max         int64
		expectError bool
	}{
		{0, 10, false},
		{10, 20, false},
		{-10, 10, false},
		{0, 0, true},   // invalid range
		{10, 10, true}, // invalid range
		{20, 10, true}, // invalid range
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("min=%d,max=%d", tt.min, tt.max), func(t *testing.T) {
			result, err := GenerateCryptoRandomNumber(tt.min, tt.max)
			if (err != nil) != tt.expectError {
				t.Errorf("GenerateCryptoRandomNumber(%d, %d) error = %v, expectError %v", tt.min, tt.max, err, tt.expectError)
			}

			if !tt.expectError {
				if result < tt.min || result >= tt.max {
					t.Errorf("GenerateCryptoRandomNumber(%d, %d) = %d, out of range [%d, %d)", tt.min, tt.max, result, tt.min, tt.max)
				}
			}
		})
	}
}
