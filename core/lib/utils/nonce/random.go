package nonce

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomString returns a randomized string
func GenerateRandomString(stringLen int) string {
	b := make([]byte, stringLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

const upperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var upperSeededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomStringUpperOnly returns a randomized string
func GenerateRandomStringUpperOnly(stringLen int) string {
	b := make([]byte, stringLen)
	for i := range b {
		b[i] = upperCharset[upperSeededRand.Intn(len(upperCharset))]
	}
	return string(b)
}
