package apitest

import (
	"fmt"
	"math/rand"
)

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

// RandLCLetters returns a random string of lower case letters from the
// English alphabet
func RandLCLetters(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandEmail returns a random email with the username and the length of the domain name provided
// Domain names and tlds are randomized
func RandEmail(username string, domainLength int) string {
	return fmt.Sprintf("%s@%s.%s", username, RandLCLetters(domainLength), RandDataString("domainSuffixes"))
}

var dataStrings = map[string][]string{
	"domainSuffixes": {"com", "biz", "info", "name", "net", "org", "io"},
}

// RandDataString returns a random data string from `name` set in dataStrings
func RandDataString(name string) string {
	return dataStrings[name][rand.Intn(len(dataStrings[name]))]
}
