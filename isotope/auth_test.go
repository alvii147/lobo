package isotope_test

import (
	"net/mail"
	"testing"

	"github.com/alvii147/lobo/isotope"
	"github.com/stretchr/testify/require"
)

func TestEmail(t *testing.T) {
	t.Parallel()

	for range 10 {
		email := isotope.Email()
		_, err := mail.ParseAddress(email)
		require.NoError(t, err)
	}
}

func TestPassword(t *testing.T) {
	t.Parallel()

	for range 10 {
		password := isotope.Password()
		require.NotEmpty(t, password)
	}
}
