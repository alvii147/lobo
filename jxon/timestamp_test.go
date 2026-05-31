package jxon_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alvii147/lobo/jxon"
	"github.com/stretchr/testify/require"
)

func TestTimestampMarshalJSON(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Timestamp jxon.Timestamp `json:"timestamp"`
	}

	p, err := json.Marshal(jsonStruct{
		Timestamp: jxon.Timestamp(time.Date(2010, 12, 16, 19, 12, 36, 0, time.UTC)),
	})
	require.NoError(t, err)
	require.Regexp(t, `^\s*{\s*"timestamp"\s*:\s*1292526756\s*}\s*$`, string(p))
}

func TestTimestampUnmarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Timestamp jxon.Timestamp `json:"timestamp"`
	}

	s := jsonStruct{}
	err := json.Unmarshal([]byte(`{"timestamp":1292526756}`), &s)
	require.NoError(t, err)
	require.Equal(t, time.Date(2010, 12, 16, 19, 12, 36, 0, time.UTC), time.Time(s.Timestamp))
}

func TestTimestampUnmarshalJSONError(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Timestamp jxon.Timestamp `json:"timestamp"`
	}

	s := jsonStruct{}
	err := json.Unmarshal([]byte(`{"timestamp":"string value"}`), &s)
	require.Error(t, err)
}
