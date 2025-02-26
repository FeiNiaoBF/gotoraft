// Package router 提供HTTP路由管理功能
package router

import (
	"time"

	"gotoraft/internal/handler"
	"gotoraft/internal/kvstore/store"
	"gotoraft/internal/observer"
	"gotoraft/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 封装gin路由器
type Router struct {
	engine    *gin.Engine
	store     *store.Store
	wsManager *websocket.Manager
	observer  *observer.RaftStateObserver
}

// NewRouter 创建一个新的路由器实例
func NewRouter(wsManager *websocket.Manager, store *store.Store, observer *observer.RaftStateObserver) *Router {
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
		store:     store,
		wsManager: wsManager,
		observer:  observer,
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
	r.registerWebSocketRoutes()

	// KV存储路由
	r.registerKVStoreRoutes()

	// TODO: 实现配置管理路由
	r.engine.GET("/api/config", r.handleGetConfig)
	r.engine.PUT("/api/config", r.handleUpdateConfig)

	// TODO: 实现集群管理路由
	r.engine.POST("/api/cluster/join", r.handleJoinCluster)
	r.engine.POST("/api/cluster/leave", r.handleLeaveCluster)

}

// registerWebSocketRoutes 注册WebSocket相关路由
func (r *Router) registerWebSocketRoutes() {
	websocketHandler := handler.NewWebSocketHandler(r.wsManager)
	websocketGroup := r.engine.Group("/ws")
	{
		// WebSocket基础路由
		websocketGroup.GET("/", websocketHandler.Handle)
		// WebSocket连接端点
		websocketGroup.GET("/connect", websocketHandler.HandleConnection)
		// 获取WebSocket统计信息
		websocketGroup.GET("/stats", websocketHandler.HandleStats)
	}
}

// registerKVStoreRoutes 注册KV存储相关路由
func (r *Router) registerKVStoreRoutes() {
	kvStoreHandler := handler.NewKVStoreHandler(r.store)
	kvStoreGroup := r.engine.Group("/api/kv")
	{
		kvStoreGroup.GET("/:key", kvStoreHandler.HandleGet)
		kvStoreGroup.POST("", kvStoreHandler.HandleSet)
		kvStoreGroup.DELETE("/:key", kvStoreHandler.HandleDelete)
	}
}

// Run 启动HTTP服务器
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
