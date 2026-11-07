package utils

import "fmt"

// NotFoundError indicates that a requested resource could not be found.
type NotFoundError struct {
	Msg string
	Err error // optional wrapped error
}

func (e *NotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("not found: %s: %v", e.Msg, e.Err)
	}
	return fmt.Sprintf("not found: %s", e.Msg)
}

func (e *NotFoundError) Unwrap() error { return e.Err }

// ValidationError indicates that input validation failed.
type ValidationError struct {
	Msg string
	Err error // optional wrapped error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error: %s: %v", e.Msg, e.Err)
	}
	return fmt.Sprintf("validation error: %s", e.Msg)
}

func (e *ValidationError) Unwrap() error { return e.Err }
