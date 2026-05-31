package sanity

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

// ValidateStringNotBlank validates that a given string is not blank.
func (v *validator) ValidateStringNotBlank(field string, value string) {
	if utf8.RuneCountInString(strings.TrimSpace(value)) < 1 {
		v.addFailuref(field, "\"%s\" cannot be blank", field)
	}
}

// ValidateStringMaxLength validates that a given string is at most of a given length.
func (v *validator) ValidateStringMaxLength(field string, value string, maxLen int) {
	if utf8.RuneCountInString(value) > maxLen {
		v.addFailuref(field, "\"%s\" cannot be more than %d characters long", field, maxLen)
	}
}

// ValidateStringMinLength validates that a given string is at least of a given length.
func (v *validator) ValidateStringMinLength(field string, value string, minLen int) {
	if utf8.RuneCountInString(value) < minLen {
		v.addFailuref(field, "\"%s\" must be at least %d characters long", field, minLen)
	}
}

// ValidateStringEmail validates the format of a given email address.
func (v *validator) ValidateStringEmail(field string, email string) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		v.addFailuref(field, "\"%s\" must be a valid email address", field)
	}
}

// ValidateStringSlug validates that a given string is a valid slug.
func (v *validator) ValidateStringSlug(field string, value string) {
	if !reSlug.MatchString(value) {
		v.addFailuref(field, "\"%s\" must be a slug", field)
	}
}

// ValidateStringOptions validates that a given string belongs to one of the given options.
func (v *validator) ValidateStringOptions(field string, value string, options []string, caseSensitive bool) {
	if !caseSensitive {
		value = strings.ToLower(value)
	}

	for _, option := range options {
		if !caseSensitive {
			option = strings.ToLower(option)
		}

		if value == option {
			return
		}
	}

	v.addFailuref(field, "\"%s\" must be one of the following options: %v", field, options)
}
