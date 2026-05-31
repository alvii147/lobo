package enigma

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/alvii147/lobo/errfmt"
	"github.com/alvii147/lobo/isotope"
	"github.com/alvii147/lobo/jxon"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWTType is a string representing type of JWT.
// Allowed strings are "access", "refresh", and "activation".
type JWTType string

const (
	// HashingCost is the cost for cryptographic hashing.
	HashingCost = 14
	// APIKeyPrefixLength is the length of API key prefixes.
	APIKeyPrefixLength = 8
	// APIKeySecretNBytes is the number of bytes in API key secrets.
	APIKeySecretNBytes = 32
	// IDPrefixJWT is the prefix for JWT IDs.
	IDPrefixJWT = "jwt"
)

// UserJWTClaims represents claims in JWTs for users.
type UserJWTClaims struct {
	Subject   string         `json:"sub"`
	TokenType string         `json:"token_type"`
	IssuedAt  jxon.Timestamp `json:"iat"`
	ExpiresAt jxon.Timestamp `json:"exp"`
	JWTID     string         `json:"jti"`
	jwt.RegisteredClaims
}

// enigma implements various cryptographic operations and helpers.
type enigma struct {
	timeProvider TimeProvider
	secretKey    string
}

// NewEnigma returns a new enigma.
func NewEnigma(
	timeProvider TimeProvider,
	secretKey string,
) *enigma {
	return &enigma{
		timeProvider: timeProvider,
		secretKey:    secretKey,
	}
}

// GenerateID generates a prefixed ID. These IDs are time-sortable.
// For example, "foo_AGLDCCTD5R5L5FHKFNYJSHKQZ4".
func (e *enigma) GenerateID(prefix string) (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", errfmt.Error(err, "uuid.NewV7 failed")
	}

	id := fmt.Sprintf("%s_%s", prefix, strings.TrimRight(base32.StdEncoding.EncodeToString(u[:]), "="))

	return id, nil
}

// HashPassword hashes a given password.
func (e *enigma) HashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), HashingCost)
	if err != nil {
		return "", errfmt.Error(err, "bcrypt.GenerateFromPassword failed")
	}
	hashedPassword := string(hashedPasswordBytes)

	return hashedPassword, nil
}

// CheckPassword checks if a given hashed password matches a given plaintext password.
func (e *enigma) CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

// CreateUserJWT creates a JWT for a user with a given type and lifetime.
func (e *enigma) CreateUserJWT(userID string, tokenType string, lifetime time.Duration) (string, error) {
	jwtID, err := e.GenerateID(IDPrefixJWT)
	if err != nil {
		return "", errfmt.Error(err)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&UserJWTClaims{
			Subject:   userID,
			TokenType: tokenType,
			IssuedAt:  jxon.Timestamp(e.timeProvider.Now()),
			ExpiresAt: jxon.Timestamp(e.timeProvider.Now().Add(lifetime)),
			JWTID:     jwtID,
		},
	)
	signedToken, err := token.SignedString([]byte(e.secretKey))
	if err != nil {
		return "", errfmt.Errorf(
			err,
			"jwt.Token.SignedString failed for user.ID %s of token type %s",
			userID,
			tokenType,
		)
	}

	return signedToken, nil
}

// ValidateUserJWT validates JWT for a user, checks that
// the JWT is not expired, and returns parsed JWT claims.
func (e *enigma) ValidateUserJWT(token string, tokenType JWTType) (*UserJWTClaims, bool) {
	claims := &UserJWTClaims{}
	ok := true

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(e.secretKey), nil
	})
	if err != nil {
		ok = false
	}

	if parsedToken == nil || !parsedToken.Valid {
		ok = false
	}

	if subtle.ConstantTimeCompare([]byte(claims.TokenType), []byte(tokenType)) == 0 {
		ok = false
	}

	if e.timeProvider.Now().After(time.Time(claims.ExpiresAt)) {
		ok = false
	}

	if !ok {
		return nil, false
	}

	return claims, true
}

// CreateAPIKey creates prefix, secret, and hashed key for API key.
func (e *enigma) CreateAPIKey() (string, string, string, error) {
	prefix, err := isotope.SecureString(APIKeyPrefixLength, true, true, true)
	if err != nil {
		return "", "", "", errfmt.Error(err)
	}

	secretBytes := make([]byte, APIKeySecretNBytes)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return "", "", "", errfmt.Error(err)
	}

	secret := base64.StdEncoding.EncodeToString(secretBytes)

	rawKey := fmt.Sprintf("%s.%s", prefix, secret)

	hashedKey, err := e.HashPassword(rawKey)
	if err != nil {
		return "", "", "", errfmt.Error(err)
	}

	return prefix, rawKey, hashedKey, nil
}

// ParseAPIKey parses API key and returns prefix and secret.
func (e *enigma) ParseAPIKey(key string) (string, string, error) {
	prefix, secret, ok := strings.Cut(key, ".")
	if !ok {
		return "", "", errfmt.Errorf(nil, "failed to parse prefix and secret for API key %s", key)
	}

	return prefix, secret, nil
}
