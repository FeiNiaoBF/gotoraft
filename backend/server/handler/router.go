package handler

import (
	"gotoraft/server/controller"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", HomeHandler)
	r.GET("/view", RaftVisualizationHandler)
	r.GET("/log", RaftLogHandler)
	r.GET("/health", ServerHealthHandler)
	r.GET("/ping", PingHandler)

	return r
}

// HomeHandler 处理主页请求
func HomeHandler(c *gin.Context) {
	controller.HomeController(c)
}

// RaftVisualizationHandler 处理 Raft 可视化请求
func RaftVisualizationHandler(c *gin.Context) {
	controller.RaftVisualizationController(c)
}

// RaftLogHandler 处理 Raft 日志请求
func RaftLogHandler(c *gin.Context) {
	controller.RaftLogController(c)
}

// ServerHealthHandler 处理服务器健康检查请求
func ServerHealthHandler(c *gin.Context) {
	controller.ServerHealthController(c)
}

// PingHandler 处理 ping 请求
func PingHandler(c *gin.Context) {
	controller.PingController(c)
}
