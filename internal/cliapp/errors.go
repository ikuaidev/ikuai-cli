package cliapp

import "fmt"

// AuthError indicates missing or invalid authentication.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string { return e.Message }

// ValidationError indicates invalid user input (bad JSON, invalid flag values).
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// NetworkError wraps connection-level failures (timeout, refused, DNS).
type NetworkError struct {
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *NetworkError) Unwrap() error { return e.Cause }
