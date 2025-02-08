// pkg/websocket/manager.go
// websocket 管理器
// 用来对 WebSocket 连接进行管理，包括建立连接、关闭连接、管理连接池等
// use gorilla/websocket: github.com/gorilla/websocket

package websocket

import (
	"encoding/json"
	"gotoraft/pkg/logger"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // 允许所有来源
}

// NewUpgrader 创建一个新的 Upgrader
func NewUpgrader() *websocket.Upgrader {
	return upgrader
}

// Manager websocket 管理器
// 管理所有的 WebSocket 连接client
type Manager struct {
	clientsMu sync.RWMutex
	clients   map[*websocket.Conn]bool
}

// NewManager 创建一个新的 Manager
func NewManager() *Manager {
	return &Manager{
		clients: make(map[*websocket.Conn]bool),
	}
}

// Register 注册新的WebSocket连接
func (m *Manager) Register(conn *websocket.Conn) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	m.clients[conn] = true
	logger.Info("新的WebSocket连接注册成功",
		"remoteAddr", conn.RemoteAddr().String(),
	)
}

// Unregister 注销WebSocket连接
func (m *Manager) Unregister(conn *websocket.Conn) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	delete(m.clients, conn)
	logger.Info("WebSocket连接注销成功",
		"remoteAddr", conn.RemoteAddr().String(),
	)
}

// GetConnectionStats 获取连接统计信息
func (m *Manager) GetConnectionStats() ConnectionStats {
	m.clientsMu.RLock()
	count := len(m.clients)
	m.clientsMu.RUnlock()

	return ConnectionStats{
		ActiveConnections: count,
		Status:            "healthy",
	}
}

// Broadcast 广播消息给所有连接的客户端
func (m *Manager) Broadcast(message []byte) {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	for client := range m.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Errorf("websocket 广播消息失败: %v", err)
			continue
		}
	}
}

// BroadcastJSON 广播JSON消息给所有连接的客户端
func (m *Manager) BroadcastJSON(data interface{}) {
	_, err := json.Marshal(data) // 验证数据可以被序列化
	if err != nil {
		logger.Errorf("websocket JSON序列化失败: %v", err)
		return
	}

	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	for client := range m.clients {
		if err := client.WriteJSON(data); err != nil {
			logger.Errorf("websocket 广播JSON消息失败: %v", err)
			continue
		}
	}
}

// ConnectionStats WebSocket连接统计信息
type ConnectionStats struct {
	ActiveConnections int    `json:"activeConnections"`
	Status            string `json:"status"`
}
