// Package errors provides custom error handling functionality for the application.
// It includes error types and utilities for creating and managing errors with
// additional context like source and message information.
package errors

import (
	"errors"
	"strings"
)

// Error represents a custom error type that includes additional context
// such as the error message, source of the error, and underlying error.
type Error struct {
	message string
	source  string
	err     error
}

// New creates a new Error instance with the provided arguments.
// Arguments can be strings (for message and source) or errors.
func New(s ...interface{}) *Error {
	err := new(Error)
	for i, v := range s {
		switch v := v.(type) {
		case string:
			if i == 0 {
				err.message = v
			} else if i == 1 {
				err.source = v
			}
		case error:
			err.err = errors.Join(v)
		}
	}
	return err
}

// Src sets the source of the error and returns a new Error instance.
// This is useful for tracking where an error originated from.
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

// Msg sets the error message and returns a new Error instance.
// This allows for customizing the error message while preserving other fields.
func (e *Error) Msg(s string) *Error {
	ne := new(Error)
	ne.err = e.err
	ne.message = s
	ne.source = e.source

	return ne
}

// Err sets the underlying error and returns a new Error instance.
// This is used to wrap another error while adding context.
func (e *Error) Err(s error) *Error {
	ne := new(Error)
	ne.err = s
	ne.message = e.message
	ne.source = e.source

	return ne
}

// Error implements the error interface.
func (e *Error) Error() string {
	format := strings.Builder{}
	format.WriteString("error: ")
	format.WriteString(e.message)

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
