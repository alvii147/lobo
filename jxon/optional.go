package jxon

import (
	"encoding/json"

	"github.com/alvii147/lobo/errfmt"
)

// Optional represents a generic JSON field that can be null or unspecified.
// If specified, Valid is true and Value holds a pointer to the specified value.
// If null, Valid is true and Value is nil.
// If unspecified, Valid is false.
type Optional[T any] struct {
	Valid bool
	Value *T
}

// MarshalJSON converts Value into bytes.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte(`null`), nil
	}

	b, err := json.Marshal(o.Value)
	if err != nil {
		return nil, errfmt.Errorf(err, "json.Marshal failed")
	}

	return b, nil
}

// UnmarshalJSON converts bytes into Optional struct.
func (o *Optional[T]) UnmarshalJSON(p []byte) error {
	o.Valid = true

	err := json.Unmarshal(p, &o.Value)
	if err != nil {
		return errfmt.Errorf(err, "json.Unmarshal failed")
	}

	return nil
}
