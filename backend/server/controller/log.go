package controller

import (
	"github.com/gin-gonic/gin"
)

// RaftLogController 处理 Raft 日志请求
func RaftLogController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Raft Log Page",
		"status":  "ok",
	})
}
