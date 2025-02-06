package controller

import (
	"github.com/gin-gonic/gin"
)

// RaftVisualizationController 处理 Raft 可视化请求
func RaftVisualizationController(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Raft Visualization Page",
		"status":  "ok",
	})
}
