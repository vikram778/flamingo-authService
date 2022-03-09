package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTPCode(length int) string {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, _ := rand.Int(rand.Reader, max)
		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice)
}
