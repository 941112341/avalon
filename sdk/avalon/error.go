package avalon

import "github.com/pkg/errors"

const UnknownErr = 500

// Code must unique
type AErr struct {
	Message string
	Code    int32
	error
}

func (A *AErr) Equal(err error) bool {
	if err == nil {
		return A == nil || A.error == err
	}
	if A.Error() == err.Error() {
		return true
	}
	a2, ok := err.(*AErr)
	if ok && a2.Code == A.Code {
		return true
	}
	return false
}

func IsError(err1, err2 error) bool {
	if err1 == err2 {
		return true
	}

	a1, ok := err1.(*AErr)
	if ok {
		return a1.Equal(err2)
	}
	return false
}

func NewError(code int32, message string) error {
	return &AErr{
		Message: message,
		Code:    code,
		error:   errors.New(message),
	}
}

func Wrap(err error) error {
	aErr, ok := err.(*AErr)
	if ok {
		return WrapWithCode(aErr, aErr.Code)
	}
	return WrapWithCode(err, UnknownErr)
}

func WrapWithCode(err error, code int32) error {
	return WrapWithMessage(err, code, "unKnown err")
}

func WrapWithMessage(err error, code int32, message string) error {
	if err == nil {
		return nil
	}
	return &AErr{
		Message: err.Error(),
		Code:    code,
		error:   errors.Wrap(err, message),
	}
}
