package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrCreateFailure is returned wrapped from generate when a token fails
	// to create
	ErrCreateFailure = errors.New("an error occurred while creating a token")

	// ErrInvalidToken is returned if a provided security token is not legit
	ErrInvalidToken = errors.New("the provided token is not valid")

	// ErrTokenExpired is returned when a token is not longer available for use
	ErrTokenExpired = errors.New("the provided token has expired")

	// ErrTokenNotFound when the token is not found
	ErrTokenNotFound = errors.New("the provided token does not exist")
)

const (
	size  = 32
	split = size / 2
)

// GenerateToken generates pieces needed for user confirm
// selector: hash of the first half of a 64 byte value
// (to be stored in the database and used in SELECT query)
// verifier: hash of the second half of a 64 byte value
// (to be stored in database but never used in SELECT query)
// token: the user-facing base64 encoded selector+verifier
func GenerateToken() (selector, verifier, token string, err error) {
	rawToken := make([]byte, size)
	if _, err = io.ReadFull(rand.Reader, rawToken); err != nil {
		return "", "", "", fmt.Errorf("%w : %v", ErrCreateFailure, err)
	}

	selectorBytes := rawToken[:split]
	verifierBytes := sha256.Sum256(rawToken[split:])

	return base64.StdEncoding.EncodeToString(selectorBytes[:]),
		base64.StdEncoding.EncodeToString(verifierBytes[:]),
		base64.URLEncoding.EncodeToString(rawToken),
		nil
}

func ExtractTokenPartsFromToken(token string) (selector, verifier string, err error) {
	rawToken, _ := base64.URLEncoding.DecodeString(token)
	if len(rawToken) != size {
		return "", "", ErrInvalidToken
	}

	selectorBytes := rawToken[:split]
	verifierBytes := sha256.Sum256(rawToken[split:])

	return base64.StdEncoding.EncodeToString(selectorBytes[:]),
		base64.StdEncoding.EncodeToString(verifierBytes[:]),
		nil
}
