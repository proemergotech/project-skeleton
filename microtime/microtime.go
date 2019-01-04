package microtime

import (
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
)

const (
	Nanosecond  = Duration(time.Nanosecond)
	Microsecond = Duration(time.Microsecond)
	Millisecond = Duration(time.Millisecond)
	Second      = Duration(time.Second)
	Minute      = Duration(time.Minute)
	Hour        = Duration(time.Hour)
)

type Time struct {
	time.Time
}

type Duration time.Duration

func Now() Time {
	return Time{Time: time.Now().UTC()}
}

func (t Time) Sub(u Time) Duration {
	return Duration(t.Time.Sub(u.Time))
}

func (t Time) String() string {
	return t.Time.UTC().Round(time.Microsecond).Format(time.RFC3339Nano)
}

func FromString(str string) (Time, error) {
	tim, err := dateparse.ParseStrict(str)
	if err != nil {
		return Time{}, err
	}

	return Time{tim.UTC().Round(time.Microsecond)}, nil
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

func (t *Time) MarshalBinary() (data []byte, err error) {
	return []byte(t.String()), nil
}

func (t *Time) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var err error
	*t, err = FromString(string(data))

	return err
}

func (t Time) RedisArg() interface{} {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Time) RedisScan(src interface{}) error {
	if src == nil {
		return nil
	}

	var str string
	switch val := src.(type) {
	case []byte:
		str = string(val)
	case string:
		str = val
	default:
		return errors.Errorf("schema.RedisScan: invalid time: %v", src)
	}

	unixTime, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return errors.Errorf("schema.RedisScan: invalid time: %v", str)
	}

	*t = Time{time.Unix(unixTime, 0)}

	return nil
}

func (d Duration) Round(m Duration) Duration {
	return Duration(time.Duration(d).Round(time.Duration(m)))
}

func (d Duration) RedisArg() interface{} {
	return strconv.FormatInt(int64(d), 10)
}

func (d *Duration) RedisScan(src interface{}) error {
	if src == nil {
		return nil
	}

	var str string
	switch val := src.(type) {
	case []byte:
		str = string(val)
	case string:
		str = val
	default:
		return errors.Errorf("schema.RedisScan: invalid duration: %v", src)
	}

	dur, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return errors.Errorf("schema.RedisScan: invalid time: %v", str)
	}

	*d = Duration(dur)

	return nil
}
