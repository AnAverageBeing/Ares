package utils

import (
	"math/rand"
	"time"
)

var (
	validNameRunes    = []rune("aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789_")
	validNameRunesLen = len(validNameRunes)
	src               = rand.NewSource(time.Now().UnixNano())
)

func RandomName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = validNameRunes[src.Int63()%int64(validNameRunesLen)]
	}
	return string(b)
}
