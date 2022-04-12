package errr

import "errors"

// StatusCodeError implement error interface
type StatusCodeError struct {
	Err        error
	StatusCode int
}

func (s StatusCodeError) Error() string {
	return s.Err.Error()
}

// New return StatusCodeError with same message, use this if we know exactly
// what error it is and what status code to return
func New(message string, statusCode int) StatusCodeError {
	return StatusCodeError{
		Err:        errors.New(message),
		StatusCode: statusCode,
	}
}

// Wrap transform error to StatusCodeError with same message,
// use this if we know exactly what status code to return
func Wrap(err error, statusCode int) StatusCodeError {
	return StatusCodeError{
		Err:        err,
		StatusCode: statusCode,
	}
}
