// Package util provides shared utility functions used across NetSentry packages.
package util

import (
	"errors"
	"fmt"
)

// Sentinel errors used across the application.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalidInput is returned when caller-supplied input fails validation.
	ErrInvalidInput = errors.New("invalid input")
	// ErrTimeout is returned when an operation exceeds its allowed duration.
	ErrTimeout = errors.New("timeout")
	// ErrPermission is returned when an operation is denied due to permissions.
	ErrPermission = errors.New("permission denied")
	// ErrUnsupported is returned when a feature or format is not supported.
	ErrUnsupported = errors.New("unsupported")
)

// Wrapf wraps an error with a formatted context message. If err is nil,
// Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// IsNotFound reports whether the error chain contains ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidInput reports whether the error chain contains ErrInvalidInput.
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsTimeout reports whether the error chain contains ErrTimeout.
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// MultiError aggregates multiple errors into a single error value.
type MultiError struct {
	Errors []error
}

// Add appends a non-nil error to the collection.
func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

// Err returns nil if no errors have been collected, otherwise returns the MultiError.
func (m *MultiError) Err() error {
	if len(m.Errors) == 0 {
		return nil
	}
	return m
}

// Error implements the error interface, formatting all collected errors.
func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	msg := fmt.Sprintf("%d error(s) occurred:", len(m.Errors))
	for i, e := range m.Errors {
		msg += fmt.Sprintf("\n  [%d] %s", i+1, e.Error())
	}
	return msg
}
