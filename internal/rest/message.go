package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/usecase"
)

// ChatHandler  represent the httphandler for article
type MessageHandler struct {
	Service usecase.MessageService
	Hub     *websocket.Conn
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewMessageService(e *echo.Echo, svc usecase.MessageService, hub *websocket.Conn) {
	handler := &MessageHandler{
		Service: svc,
		Hub:     hub,
	}

	e.GET("/api/messages", handler.GetChatMessages)
}

func (h *MessageHandler) GetChatMessages(c echo.Context) error {
	chatID, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	messages, err := h.Service.GetChatMessages(chatID, limit, offset)
	if err != nil {
		return c.JSON(getStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}
