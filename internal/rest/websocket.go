package rest

import (
	"log"
	"net/http"

	"github.com/mirasildev/chat_task/config"
	"github.com/mirasildev/chat_task/internal/websocket"
	"github.com/mirasildev/chat_task/usecase"
)

type WebSocketHandler struct {
	hub *websocket.Hub
	cfg config.Config
}

func NewWebSocketHandler(h *websocket.Hub, cfg config.Config) *WebSocketHandler {
	return &WebSocketHandler{
		hub: h,
		cfg: cfg,
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "templates/index.html")
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request, chatService usecase.ChatService, messageService usecase.MessageService) {
	hub := websocket.NewHub(chatService, messageService)
	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	log.Println("Websocket server started in port ", h.cfg.WsPort)
	err := http.ListenAndServe(h.cfg.WsPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
