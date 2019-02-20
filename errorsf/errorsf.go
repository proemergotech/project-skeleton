package errorsf

import (
	"fmt"
	"io"
)

type withFields struct {
	fields []interface{}
	error
}

func WithFields(err error, keyValues ...interface{}) error {
	return &withFields{
		fields: keyValues,
		error:  err,
	}
}

func (w *withFields) Fields() []interface{} {
	return w.fields
}

func (w *withFields) Cause() error {
	return w.error
}

func (w *withFields) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.Error())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}

func Field(err error, key interface{}) interface{} {
	type causer interface {
		Cause() error
	}

	type fielder interface {
		Fields() []interface{}
	}

	for err != nil {
		if fErr, ok := err.(fielder); ok {
			fields := fErr.Fields()
			for i := 0; i < len(fields)-1; i += 2 {
				if fields[i] == key {
					return fields[i+1]
				}
			}
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return nil
}
