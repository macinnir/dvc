package nonce

import (
	bcrypt "golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash generates a password hash from a string password
func GeneratePasswordHash(password string) string {

	passwordBytes := []byte(password)

	// Hashing the password with the default cost of 10
	hashedPassword, e := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if e != nil {
		panic(e)
	}
	return string(hashedPassword)
}

// ComparePasswordHash compares a password and a password hash
func ComparePasswordHash(password string, passwordHash string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return e == nil // nil means it is a match
}
