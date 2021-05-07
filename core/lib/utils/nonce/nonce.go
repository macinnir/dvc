package nonce

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	nonceLen = 64
)

// GenerateNonceFromString generates a sha1 encoded string from another string
func GenerateNonceFromString(str string) string {
	s := sha256.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

// ShortNonceWithPadding generates a padded nonce
func ShortNonceWithPadding(intPrefix int64, tokenLen int) string {
	baseToken := GenerateRandomStringUpperOnly(tokenLen)
	return BuildShortNonceWithPadding(intPrefix, baseToken)
}

// BuildShortNonceWithPadding generates a nonce from a baseToken and an integer prefix
func BuildShortNonceWithPadding(intPrefix int64, baseToken string) string {
	tokenLen := len(baseToken)
	intPrefixStr := strconv.FormatInt(intPrefix, 10)
	tokenRest := nonceLen - tokenLen - len(intPrefixStr)
	return intPrefixStr + strings.Repeat("0", tokenRest) + baseToken
}
