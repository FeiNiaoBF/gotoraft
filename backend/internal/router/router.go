// Package router 提供HTTP路由管理功能
package router

import (
	"time"

	"gotoraft/internal/handler"
	"gotoraft/pkg/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 封装gin路由器
type Router struct {
	engine    *gin.Engine
	wsHandler *handler.WSHandler
}

// NewRouter 创建一个新的路由器实例
func NewRouter(wsManager *websocket.Manager) *Router {
	engine := gin.New() // 使用gin.New()而不是gin.Default()以自定义中间件

	// 添加中间件
	engine.Use(gin.Logger())   // 日志中间件
	engine.Use(gin.Recovery()) // 恢复中间件

	// CORS中间件配置
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Router{
		engine:    engine,
		wsHandler: handler.NewWSHandler(wsManager),
	}
}

// RegisterRoutes 注册所有路由
func (r *Router) RegisterRoutes() {
	// 基础健康检查路由
	systemHandler := handler.NewSystemHandler()
	api := r.engine.Group("/api")
	{
		systemGroup := api.Group("/system")
		{
			systemGroup.GET("/ping", systemHandler.HandlePing)
			systemGroup.GET("/health", systemHandler.HandleHealth)
			systemGroup.GET("/info", systemHandler.HandleSystemInfo)
			systemGroup.GET("/status", systemHandler.HandleSystemStatus)
		}
	}

	// WebSocket路由
	r.engine.GET("/ws", r.wsHandler.HandleConnection)

	// 注册WebSocket路由
	r.wsHandler.RegisterRoutes(r.engine)
}

// Run 启动HTTP服务器
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

// Engine 获取底层的gin引擎实例
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
