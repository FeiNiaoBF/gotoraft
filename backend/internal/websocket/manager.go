// pkg/websocket/manager.go
// websocket 管理器
// 用来对 WebSocket 连接进行管理，包括建立连接、关闭连接、管理连接池等
// use gorilla/websocket: github.com/gorilla/websocket

package websocket

import (
	"encoding/json"
	"errors"
	"gotoraft/pkg/logger"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client WebSocket 客户端
type Client struct {
	ID         string // 客户端ID
	Conn       *websocket.Conn
	SendChan   chan []byte
	CloseChan  chan struct{}
	LastActive time.Time
}

// Manager websocket 管理器
// 管理所有的 WebSocket 连接client
type Manager struct {
	clientsMu      sync.RWMutex
	clients        map[string]*Client
	maxConnections int
	config         Config
}

// Config WebSocket 配置
type Config struct {
	MaxConnections   int           // 最大连接数
	HeartbeatTimeout time.Duration // 心跳超时时间
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // 允许所有来源
}

// NewUpgrader 创建一个新的 Upgrader
func NewUpgrader() *websocket.Upgrader {
	return upgrader
}

// NewManager 创建一个新的 Manager
func NewManager(config Config) *Manager {
	return &Manager{
		clients:        make(map[string]*Client),
		maxConnections: config.MaxConnections,
		config:         config,
	}
}

// Register 注册新的WebSocket连接
func (m *Manager) Register(conn *websocket.Conn) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	m.clients[conn.RemoteAddr().String()] = true // 使用地址字符串作为键
	logger.Info("新的WebSocket连接注册成功",
		"remoteAddr", conn.RemoteAddr().String(),
	)
}

// RegisterClient 注册客户端时生成唯一ID
func (m *Manager) RegisterClient(conn *websocket.Conn) (string, error) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()

	if len(m.clients) >= m.maxConnections {
		return "", errors.New("达到最大连接数限制")
	}

	clientID := uuid.New().String()
	client := &Client{
		ID:         clientID,
		Conn:       conn,
		SendChan:   make(chan []byte, 256),
		CloseChan:  make(chan struct{}),
		LastActive: time.Now(),
	}

	m.clients[clientID] = client
	go m.handleClient(client)

	return clientID, nil
}

// Unregister 注销WebSocket连接
func (m *Manager) Unregister(conn *websocket.Conn) {
	m.clientsMu.Lock()
	defer m.clientsMu.Unlock()
	delete(m.clients, conn.RemoteAddr().String()) // 使用地址字符串作为键
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

func (m *Manager) StartHeartbeat() {
	ticker := time.NewTicker(m.config.HeartbeatTimeout / 2)
	defer ticker.Stop()

	for range ticker.C {
		m.clientsMu.RLock()
		for _, client := range m.clients {
			if time.Since(client.LastActive) > m.config.HeartbeatTimeout {
				m.UnregisterClient(client.ID)
				continue
			}

			// 发送Ping消息
			if err := client.Conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(time.Second),
			); err != nil {
				m.UnregisterClient(client.ID)
			}
		}
		m.clientsMu.RUnlock()
	}
}
