package domain

import "fmt"

// ErrorCode represents different error types
type ErrorCode string

const (
	ErrorCodeInvalidArgument ErrorCode = "invalid_argument"
	ErrorCodeInternal        ErrorCode = "internal"
	ErrorCodeUpstreamTimeout ErrorCode = "upstream_timeout"
	ErrorCodeUpstreamError   ErrorCode = "upstream_error"
	ErrorCodeNotFound        ErrorCode = "not_found"
)

// DomainError represents a domain-level error
type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

