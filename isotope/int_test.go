package isotope_test

import (
	"testing"

	"github.com/alvii147/lobo/isotope"
	"github.com/stretchr/testify/require"
)

func TestInt64(t *testing.T) {
	t.Parallel()

	n := isotope.Int64(42)
	require.GreaterOrEqual(t, n, int64(0))
	require.Less(t, n, int64(42))
}

func TestSecureInt64(t *testing.T) {
	t.Parallel()

	n, err := isotope.SecureInt64(42)
	require.NoError(t, err)
	require.GreaterOrEqual(t, n, int64(0))
	require.Less(t, n, int64(42))
}
