// Package handler 提供HTTP和WebSocket处理器
package handler

import (
	"gotoraft/internal/websocket"
	"gotoraft/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WSHandler WebSocket处理器
type WSHandler struct {
	wsManager *websocket.Manager
}

// NewWSHandler 创建一个新的WebSocket处理器
func NewWSHandler(wsManager *websocket.Manager) *WSHandler {
	return &WSHandler{
		wsManager: wsManager,
	}
}

// HandleConnection 处理WebSocket连接请求
func (h *WSHandler) HandleConnection(c *gin.Context) {
	// 升级HTTP连接为WebSocket连接
	conn, err := websocket.NewUpgrader().Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	// 注册新的WebSocket连接
	h.wsManager.Register(conn)

	// 记录连接信息
	logger.Infof("New WebSocket connection established")
}

// HandleStats 处理WebSocket统计信息请求
func (h *WSHandler) HandleStats(c *gin.Context) {
	stats := h.wsManager.GetConnectionStats()
	c.JSON(http.StatusOK, stats)
}

// RegisterRoutes 注册WebSocket相关路由
func (h *WSHandler) RegisterRoutes(r *gin.Engine) {
	ws := r.Group("/ws")
	{
		ws.GET("/connect", h.HandleConnection) // WebSocket连接端点
		ws.GET("/stats", h.HandleStats)        // 获取WebSocket统计信息
	}
}
