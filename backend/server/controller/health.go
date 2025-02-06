package controller

import (
	"github.com/gin-gonic/gin"
)

// ServerHealthController 处理服务器健康检查请求
func ServerHealthController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server is healthy",
		"status":  "ok",
	})
}

// PingController 处理 ping 请求
func PingController(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
