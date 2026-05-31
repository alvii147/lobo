package enigma_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alvii147/lobo/enigma"
	"github.com/alvii147/lobo/isotope"
	"github.com/alvii147/lobo/jxon"
	"github.com/alvii147/lobo/timekeeper"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestEnigmaGenerateIDSuccess(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	e := enigma.NewEnigma(timeProvider, secretKey)

	id, err := e.GenerateID("foo")
	require.NoError(t, err)
	require.Len(t, id, 30)
}

func TestEnigmaHashPasswordSuccess(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	e := enigma.NewEnigma(timeProvider, secretKey)

	password := isotope.Password()
	hashedPassword, err := e.HashPassword(password)
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	require.NoError(t, err)
}

func TestEnigmaHashPasswordError(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	e := enigma.NewEnigma(timeProvider, secretKey)

	password := strings.Repeat("X", 73)
	_, err := e.HashPassword(password)
	require.Error(t, err)
}

func TestEnigmaCheckPassword(t *testing.T) {
	t.Parallel()

	secretKey := "deadbeef"
	correctPassword := isotope.Password()
	incorrectPassword := isotope.Password()

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(correctPassword), enigma.HashingCost)
	require.NoError(t, err)

	hashedPassword := string(hashedPasswordBytes)

	testcases := map[string]struct {
		password string
		wantOk   bool
	}{
		"Correct password": {
			password: correctPassword,
			wantOk:   true,
		},
		"Incorrect password": {
			password: incorrectPassword,
			wantOk:   false,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			timeProvider := timekeeper.NewFrozenProvider()
			e := enigma.NewEnigma(timeProvider, secretKey)

			ok := e.CheckPassword(hashedPassword, testcase.password)
			require.Equal(t, testcase.wantOk, ok)
		})
	}
}

func TestEnigmaCreateUserJWTSuccess(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	e := enigma.NewEnigma(timeProvider, secretKey)

	testcases := map[string]struct {
		tokenType string
		lifetime  time.Duration
	}{
		"Access token": {
			tokenType: "access",
			lifetime:  time.Hour,
		},
		"Refresh token": {
			tokenType: "refresh",
			lifetime:  30 * 24 * time.Hour,
		},
		"Activation token": {
			tokenType: "activation",
			lifetime:  30 * 24 * time.Hour,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userID := "user_XYZ"
			token, err := e.CreateUserJWT(
				userID,
				testcase.tokenType,
				testcase.lifetime,
			)
			require.NoError(t, err)

			claims := &enigma.UserJWTClaims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
				return []byte(secretKey), nil
			})
			require.NoError(t, err)

			require.NotNil(t, parsedToken)
			require.True(t, parsedToken.Valid)
			require.Equal(t, userID, claims.Subject)
			require.Equal(t, testcase.tokenType, claims.TokenType)
			require.WithinDuration(t, timeProvider.Now(), time.Time(claims.IssuedAt), time.Minute)
			require.WithinDuration(
				t,
				timeProvider.Now().Add(testcase.lifetime),
				time.Time(claims.ExpiresAt),
				time.Minute,
			)
			require.True(t, strings.HasPrefix(claims.JWTID, "jwt_"))
		})
	}
}

func TestEnigmaValidateUserJWT(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	userID := "user_XYZ"
	jti := "jwt_XYZ"
	oneDayAgo := timeProvider.Now().Add(-24 * time.Hour)
	validSecretKey := "deadbeef"

	validAccessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&enigma.UserJWTClaims{
			Subject:   userID,
			TokenType: "access",
			IssuedAt:  jxon.Timestamp(timeProvider.Now()),
			ExpiresAt: jxon.Timestamp(timeProvider.Now().Add(time.Hour)),
			JWTID:     jti,
		},
	).SignedString([]byte(validSecretKey))
	require.NoError(t, err)

	tokenOfInvalidType, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&enigma.UserJWTClaims{
			Subject:   userID,
			TokenType: string("invalidtype"),
			IssuedAt:  jxon.Timestamp(timeProvider.Now()),
			ExpiresAt: jxon.Timestamp(timeProvider.Now().Add(time.Hour)),
			JWTID:     jti,
		},
	).SignedString([]byte(validSecretKey))
	require.NoError(t, err)

	expiredToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&enigma.UserJWTClaims{
			Subject:   userID,
			TokenType: "access",
			IssuedAt:  jxon.Timestamp(oneDayAgo),
			ExpiresAt: jxon.Timestamp(oneDayAgo.Add(time.Hour)),
			JWTID:     jti,
		},
	).SignedString([]byte(validSecretKey))
	require.NoError(t, err)

	tokenWithInvalidClaim, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&struct {
			InvalidClaim string `json:"invalid_claim"`
			jwt.RegisteredClaims
		}{},
	).SignedString([]byte(validSecretKey))
	require.NoError(t, err)

	testcases := map[string]struct {
		token     string
		secretKey string
		wantOk    bool
	}{
		"Valid access token": {
			token:     validAccessToken,
			secretKey: validSecretKey,
			wantOk:    true,
		},
		"Invalid secret key": {
			token:     validAccessToken,
			secretKey: "invalidsecretkey",
			wantOk:    false,
		},
		"Token of invalid type": {
			token:     tokenOfInvalidType,
			secretKey: validSecretKey,
			wantOk:    false,
		},
		"Invalid token": {
			token:     "ed0730889507fdb8549acfcd31548ee5",
			secretKey: validSecretKey,
			wantOk:    false,
		},
		"Expired token": {
			token:     expiredToken,
			secretKey: validSecretKey,
			wantOk:    false,
		},
		"Token with invalid claim": {
			token:     tokenWithInvalidClaim,
			secretKey: validSecretKey,
			wantOk:    false,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := enigma.NewEnigma(timeProvider, testcase.secretKey)

			claims, ok := e.ValidateUserJWT(testcase.token, "access")
			require.Equal(t, testcase.wantOk, ok)

			if testcase.wantOk {
				require.Equal(t, userID, claims.Subject)
				require.Equal(t, "access", claims.TokenType)
			}
		})
	}
}

func TestEnigmaCreateAPIKey(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	e := enigma.NewEnigma(timeProvider, secretKey)

	prefix, rawKey, hashedKey, err := e.CreateAPIKey()
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(rawKey))
	require.NoError(t, err)

	r := regexp.MustCompile(`^(\S+)\.(\S+)$`)
	matches := r.FindStringSubmatch(rawKey)

	require.Len(t, matches, 3)
	require.Equal(t, prefix, matches[1])
}

func TestEnigmaParseAPIKey(t *testing.T) {
	t.Parallel()

	timeProvider := timekeeper.NewFrozenProvider()
	secretKey := "deadbeef"

	testcases := map[string]struct {
		key        string
		wantPrefix string
		wantSecret string
		wantErr    bool
	}{
		"Valid API key": {
			key:        "TqxlYSSQ.Yj2j1jyAMC5407Nctsl51K7E8sOIPqYXn28SqT5Gnfg=",
			wantPrefix: "TqxlYSSQ",
			wantSecret: "Yj2j1jyAMC5407Nctsl51K7E8sOIPqYXn28SqT5Gnfg=",
			wantErr:    false,
		},
		"Invalid API key": {
			key:        "DeAdBeEf",
			wantPrefix: "",
			wantSecret: "",
			wantErr:    true,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := enigma.NewEnigma(timeProvider, secretKey)

			prefix, secret, err := e.ParseAPIKey(testcase.key)
			if testcase.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, testcase.wantPrefix, prefix)
			require.Equal(t, testcase.wantSecret, secret)
		})
	}
}
