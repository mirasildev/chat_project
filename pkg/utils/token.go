package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mirasildev/chat_task/config"
)

type TokenParams struct {
	UserID   string
	Username string
	Email    string
	Duration time.Duration
}

// CreateToken creates a new token
func CreateToken(cfg *config.Config, params *TokenParams) (string, *Payload, error) {
	payload, err := NewPayload(params)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(cfg.TokenSymmetricKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func VerifyToken(cfg *config.Config, token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.TokenSymmetricKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
