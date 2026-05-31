package jxon

import (
	"strconv"
	"time"

	"github.com/alvii147/lobo/errfmt"
)

// Timestamp represents time that is convert to unix timestamp in JSON.
type Timestamp time.Time

// MarshalJSON converts Timestamp into bytes.
func (ut Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(ut).UTC().Unix(), 10)), nil
}

// UnmarshalJSON converts bytes into Timestamp.
func (ut *Timestamp) UnmarshalJSON(p []byte) error {
	s := string(p)
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return errfmt.Errorf(err, "strconv.ParseInt failed for timestamp %s", s)
	}

	*(*time.Time)(ut) = time.Unix(ts, 0).UTC()

	return nil
}
