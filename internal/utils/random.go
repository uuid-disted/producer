package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateCryptoRandomNumber(min, max int64) (int64, error) {
	if min >= max {
		return 0, fmt.Errorf("invalid range [%d, %d)", min, max)
	}

	rangeSize := max - min
	nBig, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return 0, err
	}

	return nBig.Int64() + min, nil
}
