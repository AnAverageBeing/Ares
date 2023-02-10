package utils

import (
	"math/rand"
	"time"
)

var (
	validNameRunes    = []rune("aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789_")
	validNameRunesLen = len(validNameRunes)
	randCh            = make(chan int, 1024)
	src               = rand.NewSource(time.Now().UnixNano())
)

func Init() {
	go func() {
		for {
			randCh <- int(src.Int63() % int64(validNameRunesLen))
		}
	}()
}

func RandomName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = validNameRunes[<-randCh]
	}
	return string(b)
}
