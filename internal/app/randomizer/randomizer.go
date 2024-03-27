package randomizer

import (
	"math/rand"
	"time"
)

const baseCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	random_bytes := make([]byte, length)
	for i := range random_bytes {
		random_char_index := seededRand.Intn(len(charset))
		random_bytes[i] = charset[random_char_index]
	}
	return string(random_bytes)
}

func String(length int) string {
	return StringWithCharset(length, baseCharset)
}
