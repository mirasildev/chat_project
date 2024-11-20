package rest

import (
	"errors"

	"github.com/mirasildev/chat_task/domain"
)

func errorResponse(err error) *domain.ResponseError {
	return &domain.ResponseError{
		Message: err.Error(),
	}
}

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrEmailExists      = errors.New("email already exists")
	ErrUserNotVerified  = errors.New("user not verified")
	ErrIncorrectCode    = errors.New("incorrect verification code")
	ErrCodeExpired      = errors.New("verification code has been expired")
	ErrNotAllowed       = errors.New("method not allowed")
	ErrWeakPassword     = errors.New("password must contain at least one small letter, one capital letter, one number, one symbol")
)
