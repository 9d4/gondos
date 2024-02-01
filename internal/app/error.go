package app

import "fmt"

type baseError struct {
	code    string
	message string
	cause   error
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

func newUserError(code string, message string, cause ...error) UserError {
	ue := UserError{baseError: baseError{
		code:    code,
		message: message,
	}}

	for _, err := range cause {
		ue.cause = fmt.Errorf("%w :%w", ue.cause, err)
	}

	return ue
}
