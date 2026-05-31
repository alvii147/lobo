package sanity

import (
	"fmt"
	"regexp"
)

// reSlug is a compiled regular expression for slug string validation.
var reSlug = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// validator validates given values and accumulates validation errors.
type validator struct {
	failures map[string][]string
}

// NewValidator creates and returns a new validator.
func NewValidator() *validator {
	return &validator{
		failures: make(map[string][]string),
	}
}

// addFailuref records a validation failure.
func (v *validator) addFailuref(field string, format string, args ...any) {
	v.failures[field] = append(v.failures[field], fmt.Sprintf(format, args...))
}

// Failures returns recorded validation failures.
func (v *validator) Failures() map[string][]string {
	return v.failures
}

// Passed returns whether or not all validations have passed.
func (v *validator) Passed() bool {
	return len(v.failures) == 0
}
