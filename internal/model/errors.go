package model

import "fmt"

type ErrorCode uint8

const (
	BadRequestCode ErrorCode = 1
	NotFoundCode   ErrorCode = 2
)

type Error struct {
	baseError error
	code      ErrorCode
	message   string
}

func (e Error) Error() string {
	switch e.code {
	case BadRequestCode:
		return e.message + ": bad request"
	case NotFoundCode:
		return e.message + ": not found"
	}

	return fmt.Sprintf("%s: unknown code %d", e.message, e.code)
}

func (e Error) Unwrap() error {
	return e.baseError
}

func (e Error) Code() ErrorCode {
	return e.code
}

func NewBadRequestErrf(format string, a ...any) Error {
	return Error{nil, BadRequestCode, fmt.Sprintf(format, a...)}
}

func WrapBadRequestErrf(e error, format string, a ...any) Error {
	return Error{e, BadRequestCode, fmt.Sprintf(format, a...)}
}

func NewNotFoundErrf(format string, a ...any) Error {
	return Error{nil, NotFoundCode, fmt.Sprintf(format, a...)}
}

func WrapNotFoundErrf(e error, format string, a ...any) Error {
	return Error{e, NotFoundCode, fmt.Sprintf(format, a...)}
}
