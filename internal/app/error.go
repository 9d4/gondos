package app

import (
	"fmt"
)

var (
	ErrUserRegistered UserError = newUserError(DuplicateErrorKind, "user.registered", "User already registered")
)

type ErrorKind string

const (
	DuplicateErrorKind  ErrorKind = "DUPLICATE"
	ValidationErrorKind ErrorKind = "VALIDATE"
)

type baseError struct {
	kind    ErrorKind
	code    string
	message string
	cause   error
}

func (e baseError) Kind() ErrorKind {
	return e.kind
}

func (e baseError) Code() string {
	return e.code
}

func (e baseError) Message() string {
	return e.message
}

func (e baseError) Error() string {
	return "app: " + e.code + ": " + e.message
}

func (e baseError) Cause() error {
	return fmt.Errorf("app: "+e.code+": "+e.message+": %w", e.cause)
}

type InternalError struct{ baseError }

type UserError struct{ baseError }

func newUserError(kind ErrorKind, code string, message string, cause ...error) UserError {
	ue := UserError{baseError: baseError{
		kind:    kind,
		code:    code,
		message: message,
	}}

	for _, err := range cause {
		ue.cause = fmt.Errorf("%w :%w", ue.cause, err)
	}

	return ue
}
