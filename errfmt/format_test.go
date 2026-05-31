package errfmt_test

import (
	"errors"
	"testing"

	"github.com/alvii147/lobo/errfmt"
	"github.com/stretchr/testify/require"
)

func errorCaller(err error, msgs ...any) error {
	return errfmt.Error(err, msgs...)
}

func errorCallerWrapper(err error, msgs ...any) error {
	return errfmt.Error(errorCaller(err, msgs...), msgs...)
}

func errorfCaller(err error, format string, args ...any) error {
	return errfmt.Errorf(err, format, args...)
}

func errorfCallerWrapper(err error, format string, args ...any) error {
	return errfmt.Errorf(errorfCaller(err, format, args...), format, args...)
}

func TestError(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		err           error
		msgs          []any
		wantErrString string
	}{
		"No messages": {
			err:  errors.New("wrapped error"),
			msgs: nil,
			wantErrString: "errfmt_test.errorCallerWrapper " +
				"-> errfmt_test.errorCaller -> wrapped error",
		},
		"Including messages": {
			err:  errors.New("wrapped error"),
			msgs: []any{"included", "messages"},
			wantErrString: "errfmt_test.errorCallerWrapper: " +
				"included messages -> " +
				"errfmt_test.errorCaller: " +
				"included messages -> wrapped error",
		},
		"Nil error": {
			err:  nil,
			msgs: []any{"included", "messages"},
			wantErrString: "errfmt_test.errorCallerWrapper: " +
				"included messages -> " +
				"errfmt_test.errorCaller: included messages",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := errorCallerWrapper(testcase.err, testcase.msgs...)
			require.Equal(t, testcase.wantErrString, err.Error())
			if testcase.err != nil {
				require.ErrorIs(t, err, testcase.err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		err           error
		format        string
		args          []any
		wantErrString string
	}{
		"No formatted messages": {
			err:    errors.New("wrapped error"),
			format: "",
			args:   nil,
			wantErrString: "errfmt_test.errorfCallerWrapper " +
				"-> errfmt_test.errorfCaller -> wrapped error",
		},
		"Including formatted messages": {
			err:    errors.New("wrapped error"),
			format: "%s-%d",
			args:   []any{"deadbeef", 42},
			wantErrString: "errfmt_test.errorfCallerWrapper: " +
				"deadbeef-42 -> errfmt_test.errorfCaller: " +
				"deadbeef-42 -> wrapped error",
		},
		"Nil error": {
			err:    nil,
			format: "%s-%d",
			args:   []any{"deadbeef", 42},
			wantErrString: "errfmt_test.errorfCallerWrapper: " +
				"deadbeef-42 -> errfmt_test.errorfCaller: " +
				"deadbeef-42",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := errorfCallerWrapper(testcase.err, testcase.format, testcase.args...)
			require.Equal(t, testcase.wantErrString, err.Error())
			if testcase.err != nil {
				require.ErrorIs(t, err, testcase.err)
			}
		})
	}
}
