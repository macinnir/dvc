package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthRequest represents the JWT json data sent back to a successfully authenticated client/user
// @model AuthRequest
// Another comment after the symbol
type AuthRequest struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// JWT is a JSON Web Token
// @model JWT
type JWT struct {
	Header    string
	Payload   string
	Signature string
}

// JWTResponse represents the JWT json data sent back to a successfully authenticated client/user
// @model JWTResponse
type JWTResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       int64  `json:"userId"`
	RefreshToken string `json:"refresh_token"`
}

// Token represents the JWT json data sent back to a successfully authenticated client/user
// @model AuthRequest
// Another comment after the symbol
type Token struct {
	ID        int64
	Val       string
	UserID    int64
	IsDeleted int
	CreatedAt string
	UpdatedAt string
}

// Public & private auth key
var authKey, authKeyPriv []byte

// GetUserIDFromAuthHeader validates the authHeaders and fetches the userID
func GetUserIDFromAuthHeader(r *http.Request) (userID int64, e error) {

	authorizationHeader := r.Header.Get("Authorization")

	if len(authorizationHeader) == 0 {
		e = errors.New("Missing authorization header")
		return
	}

	prefix := "Bearer "

	if !strings.HasPrefix(authorizationHeader, prefix) {
		e = errors.New("Invalid authorization header")
		return
	}

	if userID, e = GetUserIDFromAuthString(authorizationHeader[len(prefix):]); e != nil {
		return
	}

	return

}

// GetUserIDFromAuthString gets a userID from an authentication string
func GetUserIDFromAuthString(authString string) (userID int64, e error) {

	var authToken jwt.JWT
	var subject string
	var ok bool

	log.Printf(authString)
	if authToken, e = validateToken(authString); e != nil {
		return
	}

	if subject, ok = authToken.Claims().Subject(); !ok {
		e = errors.New("Missing subject")
		return
	}

	if userID, e = strconv.ParseInt(subject, 10, 64); e != nil {
		e = errors.New("Invalid subject")
		return
	}

	return
}

// validateToken validates and returns the JWT authentication token
func validateToken(authToken string) (auth jwt.JWT, e error) {

	// TODO read this file into memory when the application starts
	// bytes, _ := ioutil.ReadFile("./auth_key.pub")
	bytes := getPublicAuthKey()
	rsaPublic, _ := crypto.ParseRSAPublicKeyFromPEM(bytes)

	// fmt.Printf("Authorization Token %s \n", authToken)
	auth, e = jws.ParseJWT([]byte(authToken))

	if e != nil {
		return
	}

	// Validate token
	if e = auth.Validate(rsaPublic, crypto.SigningMethodRS256); e != nil {
		return
	}

	// iss Issuer
	// issuer, ok := jwt.Claims().Issuer()

	// sub Subject - Identifies the principal of the subject of the jwt
	// https://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#rfc.section.4.1.2
	// subject, ok := jwt.Claims().Subject()

	return
}

func getPublicAuthKey() []byte {

	var err error

	if authKey == nil {
		authKey, err = ioutil.ReadFile("./auth_key.pub")
		if err != nil {
			panic(err)
		}
	}

	return authKey
}

func getPrivateAuthKey() []byte {
	var err error

	if authKeyPriv == nil {
		authKeyPriv, err = ioutil.ReadFile("./auth_key.priv")
		if err != nil {
			panic(err)
		}
	}

	return authKeyPriv
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)[0:s], err
}

func generateRefreshToken() string {
	n := 32
	s, _ := GenerateRandomString(n)
	return s
}

// BuildJWT creates a JWT token based on userId and refreshToken. If RefreshToken is an empty string (len() == 0), generates a new refreshToken
func BuildJWT(userID int64, refreshToken string, expiryMinutes int, issuer string) (jwt JWTResponse) {

	// TODO Put in configuration
	// https://www.goinggo.net/2013/06/gos-duration-type-unravelled.html
	expires := time.Duration(expiryMinutes) * time.Minute
	expiresTime := time.Now().Add(expires)
	token := generateKey(expiresTime, strconv.FormatInt(userID, 10), issuer)

	if len(refreshToken) == 0 {
		refreshToken = generateRefreshToken()
	}

	jwt = JWTResponse{
		UserID:       userID,
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    int(expires.Seconds()),
		RefreshToken: refreshToken,
	}

	return
}

// ValidateToken validates a JWT authorization token
func ValidateToken(authToken string) error {

	bytes := getPublicAuthKey()
	rsaPublic, _ := crypto.ParseRSAPublicKeyFromPEM(bytes)

	// fmt.Printf("Authorization Token %s \n", authToken)

	jwt, err := jws.ParseJWT([]byte(authToken))

	if err != nil {
		return err
	}

	// Validate token
	if err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256); err != nil {
		return err
	}

	// iss Issuer
	// issuer, ok := jwt.Claims().Issuer()

	// sub Subject - Identifies the principal of the subject of the jwt
	// https://self-issued.info/docs/draft-ietf-oauth-json-web-token.html#rfc.section.4.1.2
	// subject, ok := jwt.Claims().Subject()

	return nil
}

// JOSE Standards
// Javascript Object Signing and Encryption
// http://jose.readthedocs.io/en/latest/
func generateKey(expires time.Time, subject string, issuer string) string {
	bytes := getPrivateAuthKey()
	claims := jws.Claims{}
	claims.SetExpiration(expires)
	claims.SetSubject(subject)
	claims.SetIssuer(issuer)
	rsaPrivate, _ := crypto.ParseRSAPrivateKeyFromPEM(bytes)
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)
	b, _ := jwt.Serialize(rsaPrivate)
	return string(b)
}

// EncryptPassword encrypts a password string
func EncryptPassword(password string) (hash string, e error) {

	var hashBytes []byte

	hashBytes, e = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if e != nil {
		return
	}

	hash = string(hashBytes)
	return
}

// http://jose.readthedocs.io/en/latest/#jwk
// JWK - JSON Web Key
// JSON data structure that represents a cryptographic key
// jwk = { 'k': <password> }
