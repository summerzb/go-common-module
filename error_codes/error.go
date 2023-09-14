package error_codes

import (
	"fmt"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	err     error  `json:"-"`
}

func (c *CustomError) Error() string {
	errStr := fmt.Sprintf("ErrorCode:%v ,Message:%v", c.Code, c.Message)
	if c.err != nil {
		errStr = errStr + " " + c.err.Error()
	}
	return errStr
}

func New(code int, message string) error {
	e := &CustomError{
		Code:    code,
		Message: message,
	}

	return e
}

func NewWithCode(code int) error {
	return &CustomError{
		Code:    code,
		Message: codeTextDict[code],
	}
}

func NewWithError(code int, message string, err error) error {
	e := &CustomError{
		Code:    code,
		Message: message,
		err:     err,
	}
	if message == "" {
		e.Message = codeTextDict[code]
	}

	return nil
}
