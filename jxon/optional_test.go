package jxon_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alvii147/lobo/jxon"
	"github.com/stretchr/testify/require"
)

func TestOptionalMarshalJSONSuccess(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Int    jxon.Optional[int]       `json:"int"`
		String jxon.Optional[string]    `json:"string"`
		Time   jxon.Optional[time.Time] `json:"time"`
	}

	intValue := 42
	stringValue := "deadbeef"
	timeValue := time.Date(2010, 12, 16, 19, 12, 36, 0, time.UTC)

	testcases := map[string]struct {
		intValid    bool
		intValue    *int
		stringValid bool
		stringValue *string
		timeValid   bool
		timeValue   *time.Time
		wantRegexp  string
	}{
		"Valid values": {
			intValid:    true,
			intValue:    &intValue,
			stringValid: true,
			stringValue: &stringValue,
			timeValid:   true,
			timeValue:   &timeValue,
			wantRegexp:  `^\s*{\s*"int"\s*:\s*42\s*,\s*"string"\s*:\s*"deadbeef"\s*,\s*"time"\s*:\s*"2010-12-16T19:12:36Z"}\s*$`,
		},
		"Nil values": {
			intValid:    true,
			intValue:    nil,
			stringValid: true,
			stringValue: nil,
			timeValid:   true,
			timeValue:   nil,
			wantRegexp:  `^\s*{\s*"int"\s*:\s*null\s*,\s*"string"\s*:\s*null\s*,\s*"time"\s*:\s*null}\s*$`,
		},
		"Missing values": {
			intValid:    false,
			intValue:    &intValue,
			stringValid: false,
			stringValue: &stringValue,
			timeValid:   false,
			timeValue:   &timeValue,
			wantRegexp:  `^\s*{\s*"int"\s*:\s*null\s*,\s*"string"\s*:\s*null\s*,\s*"time"\s*:\s*null}\s*$`,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			p, err := json.Marshal(jsonStruct{
				Int: jxon.Optional[int]{
					Valid: testcase.intValid,
					Value: testcase.intValue,
				},
				String: jxon.Optional[string]{
					Valid: testcase.stringValid,
					Value: testcase.stringValue,
				},
				Time: jxon.Optional[time.Time]{
					Valid: testcase.timeValid,
					Value: testcase.timeValue,
				},
			})
			require.NoError(t, err)
			require.Regexp(t, testcase.wantRegexp, string(p))
		})
	}
}

func TestOptionalMarshalJSONError(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Func jxon.Optional[func()] `json:"func"`
	}

	f := func() {}
	_, err := json.Marshal(jsonStruct{
		Func: jxon.Optional[func()]{
			Valid: true,
			Value: &f,
		},
	})
	require.Error(t, err)
}

func TestOptionalUnmarshalJSON(t *testing.T) {
	t.Parallel()

	type jsonStruct struct {
		Int    jxon.Optional[int]       `json:"int"`
		String jxon.Optional[string]    `json:"string"`
		Time   jxon.Optional[time.Time] `json:"time"`
	}

	intValue := 42
	stringValue := "deadbeef"
	timeValue := time.Date(2010, 12, 16, 19, 12, 36, 0, time.UTC)

	testcases := map[string]struct {
		data            []byte
		wantErr         bool
		wantIntValid    bool
		wantInt         *int
		wantStringValid bool
		wantString      *string
		wantTimeValid   bool
		wantTime        *time.Time
	}{
		"Valid values": {
			data: []byte(`{
				"int": 42,
				"string": "deadbeef",
				"time": "2010-12-16T19:12:36Z"
			}`),
			wantErr:         false,
			wantIntValid:    true,
			wantInt:         &intValue,
			wantStringValid: true,
			wantString:      &stringValue,
			wantTimeValid:   true,
			wantTime:        &timeValue,
		},
		"Null values": {
			data: []byte(`{
				"int": null,
				"string": null,
				"time": null
			}`),
			wantErr:         false,
			wantIntValid:    true,
			wantInt:         nil,
			wantStringValid: true,
			wantString:      nil,
			wantTimeValid:   true,
			wantTime:        nil,
		},
		"Missing values": {
			data:            []byte(`{}`),
			wantErr:         false,
			wantIntValid:    false,
			wantInt:         nil,
			wantStringValid: false,
			wantString:      nil,
			wantTimeValid:   false,
			wantTime:        nil,
		},
		"Invalid values": {
			data: []byte(`{
				"int": "dead",
				"string": 42,
				"time": "beef"
			}`),
			wantErr:         true,
			wantIntValid:    false,
			wantInt:         nil,
			wantStringValid: false,
			wantString:      nil,
			wantTimeValid:   false,
			wantTime:        nil,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := jsonStruct{}
			err := json.Unmarshal(testcase.data, &s)
			if testcase.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			require.Equal(t, testcase.wantIntValid, s.Int.Valid)
			if s.Int.Valid {
				require.Equal(t, testcase.wantInt, s.Int.Value)
			}

			require.Equal(t, testcase.wantStringValid, s.String.Valid)
			if s.String.Valid {
				require.Equal(t, testcase.wantString, s.String.Value)
			}

			require.Equal(t, testcase.wantTimeValid, s.Time.Valid)
			if s.Time.Valid {
				require.Equal(t, testcase.wantTime, s.Time.Value)
			}
		})
	}
}
