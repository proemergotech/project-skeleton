package microtime

import (
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	return t.UTC().Round(time.Microsecond).MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	t.Time, err = dateparse.ParseStrict(unquoted)
	if err != nil {
		return err
	}
	t.Time = t.Time.UTC()

	return nil
}
