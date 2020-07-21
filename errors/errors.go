package errors

import (
	"encoding/json"

	"fmt"
	"runtime"
)

var (
	prefix = "Error Code: "
	period = ". "
)

type Error struct {
	Act     string `json:"-"`
	Code    int    `json:"c"`
	Message string `json:"m"`
}

func New(act string, code int, errMessage string) *Error {
	return &Error{Act: act, Code: code, Message: errMessage}
}

func (e Error) String() string {
	return prefix + fmt.Sprintf("%v", e.Code) + period + e.Act + period + e.Message
}

func (e Error) Error() string {
	return e.String()
}

func (e Error) toJsonString() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Error) Format(args ...interface{}) Error {
	e.Message = fmt.Sprintf(e.Message, args...)
	return e
}

func (e Error) Append(format string, a ...interface{}) Error {
	e.Message += fmt.Sprintf(format, a...)
	return e
}

func (e Error) AppendErr(err error) Error {
	return e.Append(err.Error())
}

func (e Error) Panic() {
	_, fn, line, _ := runtime.Caller(1)
	errMsg := e.Message
	errMsg += "\nCaller was: " + fmt.Sprintf("%s:%d", fn, line)
	panic(errMsg)
}

func (e Error) Panicf(args ...interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	errMsg := e.Format(args...).Error()
	errMsg += "\nCaller was: " + fmt.Sprintf("%s:%d", fn, line)
	panic(errMsg)
}


func BadRequest(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    400,
		Message: msg,
	}
}

func Unauthorized(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    401,
		Message: msg,
	}
}

func PaymentRequired(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    402,
		Message: msg,
	}
}

func Forbidden(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    403,
		Message: msg,
	}
}

func NotFound(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    404,
		Message: msg,
	}
}

func MethodNotAllowed(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    405,
		Message: msg,
	}
}

func InternalServerError(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    500,
		Message: msg,
	}
}

func NotImplemented(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    501,
		Message: msg,
	}
}

func BadGateway(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    502,
		Message: msg,
	}
}

func ServiceUnavailable(act, msg string) error {
	return &Error{
		Act:     act,
		Code:    503,
		Message: msg,
	}
}