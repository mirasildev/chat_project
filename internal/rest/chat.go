package rest

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/internal/rest/middleware"
	"github.com/mirasildev/chat_task/usecase"
)

// ChatHandler  represent the httphandler for article
type ChatHandler struct {
	Service usecase.ChatService
}

func NewChatService(e *echo.Echo, svc usecase.ChatService) {
	handler := &ChatHandler{
		Service: svc,
	}

	e.POST("/api/groups/create", handler.CreateChat)
	e.GET("/api/groups", handler.GetUserChats)
	e.PUT("/api/groups/join", handler.JoinChat)
	e.DELETE("/api/groups/leave", handler.LeaveChat)
}

func (h *ChatHandler) CreateChat(c echo.Context) error {
	var chat domain.Chat
	err := c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	payload, err := middleware.GetAuthPayload(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseError{Message: err.Error()})
	}

	chat.CreatedBy = payload.UserID
	err = h.Service.CreateChat(&chat)
	if err != nil {
		return c.JSON(getStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, chat)
}

func (h *ChatHandler) GetUserChats(c echo.Context) error {
	payload, err := middleware.GetAuthPayload(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseError{Message: err.Error()})
	}

	groups, err := h.Service.GetUserChats(payload.UserID)
	if err != nil {
		return c.JSON(getStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, groups)
}

func validateUUID(id string) (string, error) {
	respID, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return respID.String(), nil
}

func (h *ChatHandler) JoinChat(c echo.Context) error {
	chatID, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}
	payload, err := middleware.GetAuthPayload(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseError{Message: err.Error()})
	}

	if err := h.Service.JoinChat(int64(chatID), payload.UserID); err != nil {
		return c.JSON(getStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *ChatHandler) LeaveChat(c echo.Context) error {
	chatID, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}
	payload, err := middleware.GetAuthPayload(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseError{Message: err.Error()})
	}

	if err := h.Service.LeaveChat(chatID, payload.UserID); err != nil {
		return c.JSON(getStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
