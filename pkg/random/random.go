package random

import (
	"time"

	"golang.org/x/exp/rand"
)

const (
	AvailableCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// NewRandomString generates random string with given size.
func NewRandomString(size uint) string {
	if size == 0 {
		return ""
	}

	rnd := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	runes := []rune(AvailableCharacters)

	b := make([]rune, size)
	for i := range b {
		randSymb := rnd.Intn(len(runes))
		b[i] = runes[randSymb]
	}

	return string(b)
}
