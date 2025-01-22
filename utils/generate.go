package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(strLength int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Create a slice of random bytes
	b := make([]byte, strLength)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}
