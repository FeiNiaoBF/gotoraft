package main

import (
	"fmt"
	"log"

	"gotoraft/internal/config"
	"gotoraft/server/handler"

	"github.com/gin-gonic/gin"
)

// 启动服务器 main

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}

	r := gin.Default()

	// 设置路由
	r.GET("/", handler.HomeHandler)
	r.GET("/view", handler.RaftVisualizationHandler)
	r.GET("/log", handler.RaftLogHandler)
	r.GET("/health", handler.ServerHealthHandler)
	r.GET("/ping", handler.PingHandler)

	// 启动服务器
	address := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Printf("Starting server on %s\n", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
