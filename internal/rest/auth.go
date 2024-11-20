package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/usecase"
)

type AuthHandler struct {
	Service     usecase.AuthService
	UserService usecase.UserService
}

func NewAuthService(e *echo.Echo, svc usecase.AuthService, userService usecase.UserService) {
	handler := &AuthHandler{
		Service:     svc,
		UserService: userService,
	}

	e.POST("/auth/register", handler.Register)
	e.POST("/auth/verify", handler.Verify)
	e.POST("/auth/login", handler.Login)
}

func (a *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if !validatePassword(req.Password) {
		return c.JSON(http.StatusBadRequest, errorResponse(ErrWeakPassword))
	}

	_, err = a.UserService.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
	}

	err = a.Service.Register(&domain.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"Message": "success",
	})
}

func validatePassword(password string) bool {
	var capitalLetter, smallLetter, number, symbol bool

	for i := 0; i < len(password); i++ {
		if password[i] >= 65 && password[i] <= 90 {
			capitalLetter = true
		} else if password[i] >= 97 && password[i] <= 122 {
			smallLetter = true
		} else if password[i] >= 48 && password[i] <= 57 {
			number = true
		} else {
			symbol = true
		}
	}
	return capitalLetter && smallLetter && number && symbol
}
