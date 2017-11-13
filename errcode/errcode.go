package errcode

import "fmt"

// ErrorCode .
type ErrorCode struct {
	Code    int
	Message string
}

func (code ErrorCode) Error() string {
	return code.Message
}

// New create new error code
func New(code int, message string) *ErrorCode {
	return &ErrorCode{
		Code:    code,
		Message: fmt.Sprintf("code :%d %s", code, message),
	}
}
