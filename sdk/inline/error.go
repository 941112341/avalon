package inline

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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
}

type CodeError struct {
	err     error
	message string
	code    AvalonErrorCode
	stack   string
}

func (c CodeError) Error() string {
	msg := c.message
	if c.code != Unknown {
		msg = fmt.Sprintf("(%d):%s", c.code, msg)
	}
	if c.stack != "" {
		msg = fmt.Sprintf("%s:%s", c.stack, c.message)
	}
	if c.err == nil {
		return msg
	}
	return fmt.Sprintf("%s\n%s", msg, c.err.Error())
}

func (c CodeError) RawError() error {
	aErr, ok := c.err.(AvalonError)
	if ok {
		return aErr.RawError()
	}
	return c.err
}

func (c CodeError) Code() AvalonErrorCode {
	if c.code == 0 {
		aErr, ok := c.err.(AvalonError)
		if ok {
			return aErr.Code()
		}
	}
	return c.code
}

func NewError(code AvalonErrorCode, f string, args ...interface{}) AvalonError {
	return &CodeError{
		message: fmt.Sprintf(f, args...),
		code:    code,
		stack:   RecordStack(2),
	}
}

func Error(f string, args ...interface{}) AvalonError {
	return &CodeError{
		err:     nil,
		message: fmt.Sprintf(f, args...),
		code:    0,
		stack:   RecordStack(2),
	}
}

func PrependErrorFmt(err error, f string, args ...interface{}) error {
	return &CodeError{
		err:     err,
		message: fmt.Sprintf(f, args...),
		code:    0,
		stack:   RecordStack(2),
	}
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

var GOPATH string

func init() {
	GOPATH = os.Getenv("GOPATH")
}

type Stack struct {
	File     string
	FuncName string
	Line     int
	Ok       bool
	Index    int
}

func (s Stack) ToString() string {
	if !s.Ok {
		return fmt.Sprintf("[index %d not ok]", s.Index)
	}
	return fmt.Sprintf("[%s:%d]", s.File, s.Line)
}

// 0代表当前line
func RecordStack(index int) string {

	// 忽略此层方法
	_, file, line, ok := runtime.Caller(index + 1)
	file = strings.TrimPrefix(file, GOPATH)
	stack := &Stack{
		File: file,
		//FuncName: runtime.FuncForPC(pc).Name(),
		Line:  line,
		Ok:    ok,
		Index: index,
	}
	return stack.ToString()
}
