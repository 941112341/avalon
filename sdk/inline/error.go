package inline

import (
	"fmt"
)

type AvalonErrorCode int32

const (
	_ AvalonErrorCode = iota
	Unknown
	ErrArgs
	ErrDBCas
)

func (e AvalonErrorCode) I32() int32 {
	return int32(e)
}

type AvalonError interface {
	Error() string
	RawError() error
	Code() AvalonErrorCode
	WrapErr(message string, codes ...AvalonErrorCode) AvalonError
	Unwrap() error
}

type CodeError struct {
	err     error
	message string
	code    AvalonErrorCode
}

func (c CodeError) Error() string {
	msg := c.message
	if c.code != Unknown {
		msg = fmt.Sprintf("%s:%s", c.code, msg)
	}
	if c.err == nil {
		return msg
	}
	return msg + fmt.Sprintf("[%s]", c.err.Error())
}

func (c CodeError) RawError() error {
	aErr, ok := c.err.(AvalonError)
	if ok {
		return aErr.RawError()
	}
	return c.err
}

func (c CodeError) Code() AvalonErrorCode {
	return c.code
}

func (c CodeError) WrapErr(message string, codes ...AvalonErrorCode) AvalonError {
	parentCode := c.code
	if len(codes) != 0 {
		parentCode = codes[0]
	} else {
		c.code = Unknown // 避免重复打印无意义code
	}
	return &CodeError{
		err:     c,
		message: message,
		code:    parentCode,
	}
}

func (c CodeError) Unwrap() error {
	return c.err
}

func NewError(code AvalonErrorCode, f string, args ...interface{}) AvalonError {
	return &CodeError{
		message: fmt.Sprintf(f, args...),
		code:    code,
	}
}

func Error(f string, args ...interface{}) AvalonError {
	return NewError(Unknown, f, args...)
}

func PrependErrorWithCode(err error, code AvalonErrorCode, f string, args ...interface{}) error {
	message := fmt.Sprintf(f, args...)
	aErr, ok := err.(AvalonError)
	if ok {
		return aErr.WrapErr(message)
	}
	return &CodeError{
		err:     err,
		message: message,
		code:    code,
	}
}

func PrependError(err error, message string) error {
	if aErr, ok := err.(AvalonError); ok {
		return PrependErrorWithCode(err, aErr.Code(), message)
	}
	return PrependErrorWithCode(err, Unknown, message)
}

func PrependErrorFmt(err error, f string, args ...interface{}) error {
	return PrependError(err, fmt.Sprintf(f, args...))
}

// compare
func IsErr(err1, err2 error) bool {
	if err1 == err2 {
		return true
	}
	aErr1, ok := err1.(AvalonError)
	if ok {
		if aErr1.RawError() == err2 {
			return true
		}
	}
	aErr2, ok2 := err2.(AvalonError)
	if ok2 {
		if aErr2.RawError() == err1 {
			return true
		}
	}
	if ok && ok2 {
		if aErr1.RawError() == aErr2.RawError() {
			return true
		}
	}
	return false
}

func IsCode(err error, errCode AvalonErrorCode) bool {
	aErr, ok := err.(AvalonError)
	if !ok {
		return false
	}
	return aErr.Code() == errCode
}
