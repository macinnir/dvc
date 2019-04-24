package apitest

import "math/rand"

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyz_1234567890")
var numberRunes = []rune("1234567890")
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// RandString returns a random string based on a set of runes
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

// RandInt generates a random number where 0 <= n <= max
func RandInt(max int) int64 {
	return int64(rand.Intn(max))
}
