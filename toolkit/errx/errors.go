package errx

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type Error struct {
	Code    int32
	Message string
	cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.cause
}

// New 创建一个新的自定义错误
func New(code int32, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Newf 格式化创建错误
func Newf(code int32, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

// WithCause 返回带有底层 cause 的新错误
func (e *Error) WithCause(cause error) *Error {
	newErr := *e
	newErr.cause = cause
	return &newErr
}

// FromError 转换任意 error 为 *Error 类型
func FromError(err error) *Error {
	if err == nil {
		return New(http.StatusOK, "请求成功")
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	// 非自定义错误，封装成未知错误
	return New(CodeUnknown, err.Error())
}

// GORMErr 根据 gorm 错误类型转换为业务错误
func GORMErr(err error) *Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(CodeNoData, "数据不存在")
	}
	return New(CodeOrmInvalid, err.Error())
}

// GORMErrMsg 支持自定义找不到数据时的错误消息
func GORMErrMsg(err error, msg string) *Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if msg == "" {
			return New(CodeNoData, "数据不存在")
		}
		return New(CodeNoData, msg)
	}
	return New(CodeOrmInvalid, err.Error())
}

// BizErr 业务逻辑错误快捷方法
func BizErr(msg string) *Error {
	return New(CodeBizErr, msg)
}

// AuthErr 认证相关错误快捷方法
func AuthErr(msg string) *Error {
	return New(CodeLoginErr, msg)
}
