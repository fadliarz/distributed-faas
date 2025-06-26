package core

import (
	"fmt"
	"strings"
)

type ErrorType string

const (
	ValidationError   ErrorType = "VALIDATION_ERROR"
	UnauthorizedError ErrorType = "UNAUTHORIZED_ERROR"
	NotFoundError     ErrorType = "NOT_FOUND_ERROR"
	DatabaseError     ErrorType = "DATABASE_ERROR"
	InternalError     ErrorType = "INTERNAL_ERROR"
)

type ServiceError struct {
	ErrorType ErrorType
	Message   string
	Err       error
}

func (e *ServiceError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", e.ErrorType, e.Message))
	if e.Err != nil {
		sb.WriteString(fmt.Sprintf(": %v", e.Err))
	}
	return sb.String()
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// Constructors for different error types

func NewValidationError(message string, err error) error {
	return &ServiceError{
		ErrorType: ValidationError,
		Message:   message,
		Err:       err,
	}
}

func NewUnauthorizedError(message string, err error) error {
	return &ServiceError{
		ErrorType: UnauthorizedError,
		Message:   message,
		Err:       err,
	}
}

func NewNotFoundError(message string, err error) error {
	return &ServiceError{
		ErrorType: NotFoundError,
		Message:   message,
		Err:       err,
	}
}

func NewDatabaseError(message string, err error) error {
	return &ServiceError{
		ErrorType: DatabaseError,
		Message:   message,
		Err:       err,
	}
}

func NewInternalError(message string, err error) error {
	return &ServiceError{
		ErrorType: InternalError,
		Message:   message,
		Err:       err,
	}
}

func IsErrorType(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ServiceError); ok {
		return e.ErrorType == errorType
	}

	return false
}
