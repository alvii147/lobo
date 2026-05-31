package isotope

import (
	cryptorand "crypto/rand"
	"math/big"
	mathrand "math/rand/v2"

	"github.com/alvii147/lobo/errfmt"
)

// Int64 generates a random integer from the interval [0, max).
// This is not cryptographically secure.
func Int64(maxInt int64) int64 {
	return mathrand.Int64N(maxInt)
}

// SecureInt64 generates a cryptographically secure random integer from the interval [0, max).
func SecureInt64(maxInt int64) (int64, error) {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(maxInt))
	if err != nil {
		return 0, errfmt.Error(err, "rand.Int failed")
	}

	return n.Int64(), nil
}
