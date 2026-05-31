package timekeeper_test

import (
	"testing"
	"time"

	"github.com/alvii147/lobo/timekeeper"
	"github.com/stretchr/testify/require"
)

func TestNewSystemProvider(t *testing.T) {
	t.Parallel()

	provider := timekeeper.NewSystemProvider()
	require.WithinDuration(t, time.Now().UTC(), provider.Now(), time.Minute)
}
