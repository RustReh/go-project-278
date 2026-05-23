package apperr

import (
	"errors"
	"fmt"
)

// Code — машинный код для API и ветвления в handler.
type Code string

const (
	CodeNotFound   Code = "NOT_FOUND"
	CodeValidation Code = "VALIDATION"
	CodeConflict   Code = "CONFLICT"
	CodeInternal   Code = "INTERNAL"
)

// AppError — единый формат ошибки приложения.
type AppError struct {
	Code    Code
	Message string
	Payload any
	Err     error // внутренняя причина (sql, fmt, …)
}

// Error реализует стандартный интерфейс error.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap позволяет errors.Is / errors.As по цепочке.
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func WithPayload(code Code, message string, payload any, err error) *AppError {
	return &AppError{Code: code, Message: message, Payload: payload, Err: err}
}

func NotFound(message string) *AppError {
	return &AppError{Code: CodeNotFound, Message: message}
}

func Validation(message string, payload any) *AppError {
	return &AppError{Code: CodeValidation, Message: message, Payload: payload}
}

func Conflict(message string) *AppError {
	return &AppError{Code: CodeConflict, Message: message}
}

func Internal(message string, err error) *AppError {
	return &AppError{Code: CodeInternal, Message: message, Err: err}
}

// AsAppError извлекает *AppError из цепочки error.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
