// Package router 提供HTTP路由管理功能
package router

import (
	"net/http"
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
	r.engine.GET("/ping", r.handlePing)
	r.engine.GET("/health", r.handleHealth)

	// API版本v1
	v1 := r.engine.Group("/api/v1")
	{
		// 系统信息路由组
		systemGroup := v1.Group("/system")
		{
			systemGroup.GET("/info", r.handleSystemInfo)
			systemGroup.GET("/status", r.handleSystemStatus)
		}
	}

	// 注册WebSocket路由
	r.wsHandler.RegisterRoutes(r.engine)
}

// handlePing 处理ping请求
func (r *Router) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// handleHealth 处理健康检查请求
func (r *Router) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0", // 可以从配置或环境变量中获取
	})
}

// handleSystemInfo 处理系统信息请求
func (r *Router) handleSystemInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "GoToRaft",
		"version":     "1.0.0",
		"environment": gin.Mode(),
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

// handleSystemStatus 处理系统状态请求
func (r *Router) handleSystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"uptime": time.Now().Format(time.RFC3339),
		"services": gin.H{
			"api":     "healthy",
			"storage": "healthy",
		},
	})
}

// Run 启动HTTP服务器
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

// Engine 获取底层的gin引擎实例
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
