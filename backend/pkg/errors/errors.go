package errors

import (
	"fmt"
	"net/http"
)

// Error 自定义错误类型
type Error struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Err     error  `json:"-"`      // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

// NewError 创建新的错误
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// WrapError 包装已有错误
func WrapError(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// StatusCode 获取HTTP状态码
func (e *Error) StatusCode() int {
	switch {
	case e.Code >= 10000 && e.Code < 20000:
		return http.StatusBadRequest
	case e.Code >= 20000 && e.Code < 30000:
		return http.StatusUnauthorized
	case e.Code >= 30000 && e.Code < 40000:
		return http.StatusForbidden
	case e.Code >= 40000 && e.Code < 50000:
		return http.StatusNotFound
	case e.Code >= 50000 && e.Code < 60000:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
