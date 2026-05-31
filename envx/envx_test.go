package envx_test

import (
	"testing"

	"github.com/alvii147/lobo/envx"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigSuccess(t *testing.T) {
	t.Setenv("BOOLTRUE", "true")
	t.Setenv("BOOLFALSE", "false")
	t.Setenv("INT", "42")
	t.Setenv("INT64", "-42")
	t.Setenv("UINT64", "42")
	t.Setenv("FLOAT32", "3.2")
	t.Setenv("FLOAT64", "-6.4")
	t.Setenv("STRING", "Hello")

	config := struct {
		BoolTrue  bool    `env:"BOOLTRUE"`
		BoolFalse bool    `env:"BOOLFALSE"`
		Int       int     `env:"INT"`
		Int64     int64   `env:"INT64"`
		Uint64    uint64  `env:"UINT64"`
		Float32   float32 `env:"FLOAT32"`
		Float64   float64 `env:"FLOAT64"`
		String    string  `env:"STRING"`
	}{}

	err := envx.LoadConfig(&config)
	require.NoError(t, err)

	require.True(t, config.BoolTrue)
	require.False(t, config.BoolFalse)
	require.Equal(t, 42, config.Int)
	require.EqualValues(t, -42, config.Int64)
	require.EqualValues(t, 42, config.Uint64)
	require.InDelta(t, 3.2, config.Float32, 0.1)
	require.InDelta(t, -6.4, config.Float64, 0.1)
	require.Equal(t, "Hello", config.String)
}

func TestLoadConfigInvalidBoolError(t *testing.T) {
	t.Setenv("BOOL", "GR3Y")

	config := struct {
		Bool bool `env:"BOOL"`
	}{}

	err := envx.LoadConfig(&config)
	require.Error(t, err)
}

func TestLoadConfigInvalidIntError(t *testing.T) {
	t.Setenv("INT", "1o2")

	config := struct {
		Int int `env:"INT"`
	}{}

	err := envx.LoadConfig(&config)
	require.Error(t, err)
}

func TestLoadConfigInvalidUintError(t *testing.T) {
	t.Setenv("UINT", "-42")

	config := struct {
		Uint uint `env:"UINT"`
	}{}

	err := envx.LoadConfig(&config)
	require.Error(t, err)
}

func TestLoadConfigInvalidFloatError(t *testing.T) {
	t.Setenv("FLOAT", "B33F")

	config := struct {
		Float float64 `env:"FLOAT"`
	}{}

	err := envx.LoadConfig(&config)
	require.Error(t, err)
}

func TestLoadConfigNoTag(t *testing.T) {
	t.Setenv("STRING", "D34D")

	config := struct {
		Tag   string `env:"STRING"`
		NoTag string
	}{
		NoTag: "B33F",
	}

	err := envx.LoadConfig(&config)
	require.NoError(t, err)

	require.Equal(t, "D34D", config.Tag)
	require.Equal(t, "B33F", config.NoTag)
}

func TestLoadConfigInvalidTypeError(t *testing.T) {
	t.Setenv("STRING", "D3AD")

	config := struct {
		String struct{} `env:"STRING"`
	}{}

	err := envx.LoadConfig(&config)
	require.Error(t, err)
}
