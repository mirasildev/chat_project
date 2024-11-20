package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/usecase"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

type Payload struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IssuedAt  string `json:"issued_at"`
	ExpiredAt string `json:"expired_at"`
}

func AuthMiddleware(authService usecase.AuthService, resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get(authorizationHeaderKey)
			if len(accessToken) == 0 {
				err := errors.New("authorization header is not provided")
				return c.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			payload, err := authService.VerifyToken(context.Background(), &domain.VerifyTokenRequest{
				AccessToken: accessToken,
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, errorResponse(err))
			}

			c.Set(authorizationPayloadKey, Payload{
				ID:        payload.Id,
				UserID:    payload.UserID,
				Email:     payload.Email,
				Username:  payload.Username,
				IssuedAt:  payload.IssuedAt,
				ExpiredAt: payload.ExpiredAt,
			})

			// Continue to next handler
			return next(c)
		}
	}
}

func GetAuthPayload(c echo.Context) (*Payload, error) {
	value := c.Get(authorizationPayloadKey)
	if value == nil {
		return nil, errors.New("no authentication payload found")
	}

	payload, ok := value.(Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}

	return &payload, nil
}

func errorResponse(err error) *domain.ResponseError {
	return &domain.ResponseError{
		Message: err.Error(),
	}
}
