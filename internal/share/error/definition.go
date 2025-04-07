package errors

import (
	"errors"
	"strings"
)

type Error struct {
	message string
	source  string
	err     error
}

func New(s ...interface{}) *Error {
	err := new(Error)
	for i := 0; i < len(s); i++ {
		switch s[i].(type) {
		case string:
			if i == 0 {
				err.message = s[i].(string)
			} else if i == 1 {
				err.source = s[i].(string)
			}
		case error:
			err.err = errors.Join(s[i].(error))
		}
	}

	return err
}

// WARN: Concurrent modification of "constant" fields.
// WARN: Maybe we need to copy the error.
// WARN: But then it will cause memory problems. Probably GC will deal with them
func (e *Error) Src(s string) *Error {
	ne := new(Error)
	ne.err = e.err
	ne.message = e.message
	ne.source = s

	return ne
}

func (e *Error) Msg(s string) *Error {
	ne := new(Error)
	ne.err = e.err
	ne.message = s
	ne.source = e.source

	return ne
}

func (e *Error) Err(s error) *Error {
	ne := new(Error)
	ne.err = s
	ne.message = e.message
	ne.source = e.source

	return ne
}

func (e *Error) Error() string {
	format := strings.Builder{}
	format.WriteString("error: ")
	format.WriteString(e.err.Error())

	if e.message != "" {
		format.WriteString(e.message)
	}

	if e.source != "" {
		format.WriteString(" | src: ")
		format.WriteString(e.source)
	}

	if e.err != nil {
		format.WriteString(" | err: ")
		format.WriteString(e.err.Error())
	}

	return format.String()
}
