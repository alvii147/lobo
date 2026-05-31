package isotope

import (
	"github.com/alvii147/lobo/errfmt"
)

const (
	// CharsLowerAlpha includes all lowercase alphabets.
	CharsLowerAlpha = "abcdefghijklmnopqrstuvwxyz"
	// CharsUpperAlpha includes all uppercase alphabets.
	CharsUpperAlpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// CharsUpperAlpha includes all numeric characters.
	CharsNumeric = "0123456789"
)

// String generates random string of given length.
// allowLowerAlpha, allowUpperAlpha, and allowNumeric can be used to include/exclude
// lowercase, uppercase, and numeric characters respectively.
// This is not cryptographically secure.
func String(
	n int,
	allowLowerAlpha bool,
	allowUpperAlpha bool,
	allowNumeric bool,
) string {
	lowerAlpha := ""
	if allowLowerAlpha {
		lowerAlpha = CharsLowerAlpha
	}

	upperAlpha := ""
	if allowUpperAlpha {
		upperAlpha = CharsUpperAlpha
	}

	numeric := ""
	if allowNumeric {
		numeric = CharsNumeric
	}

	allowed := []rune(lowerAlpha + upperAlpha + numeric)
	randRunes := make([]rune, n)

	for i := range randRunes {
		randRunes[i] = allowed[Int64(int64(len(allowed)))]
	}

	randString := string(randRunes)

	return randString
}

// SecureString generates a cryptographically secure random string of given length.
// allowLowerAlpha, allowUpperAlpha, and allowNumeric can be used to include/exclude
// lowercase, uppercase, and numeric characters respectively.
func SecureString(
	n int,
	allowLowerAlpha bool,
	allowUpperAlpha bool,
	allowNumeric bool,
) (string, error) {
	lowerAlpha := ""
	if allowLowerAlpha {
		lowerAlpha = CharsLowerAlpha
	}

	upperAlpha := ""
	if allowUpperAlpha {
		upperAlpha = CharsUpperAlpha
	}

	numeric := ""
	if allowNumeric {
		numeric = CharsNumeric
	}

	allowed := []rune(lowerAlpha + upperAlpha + numeric)
	randRunes := make([]rune, n)

	for i := range randRunes {
		n, err := SecureInt64(int64(len(allowed)))
		if err != nil {
			return "", errfmt.Error(err)
		}

		randRunes[i] = allowed[n]
	}

	randString := string(randRunes)

	return randString, nil
}
