package websocket

import (
	"gotoraft/internal/raft"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

type Handler struct {
	controller *raft.RaftNode
}

func NewHandler(controller *raft.RaftNode) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientID := r.URL.Query().Get("clientId")
	h.controller.RegisterClient(clientID, conn)

	// 处理客户端消息
	go h.handleMessages(clientID, conn)
}

func (h *Handler) handleMessages(clientID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.controller.UnregisterClient(clientID)
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		// 处理来自前端的命令
		h.handleCommand(messageType, p)
	}
}
