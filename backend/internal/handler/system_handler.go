// Package handler 提供HTTP和WebSocket处理器
// 一些系统级的处理
package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	version = "1.0.0"
)

type status string

const (
	running status = "running"
	stopped status = "stopped"
	healthy status = "healthy"
	pong    status = "pong"
	error   status = "error"
)

type Response struct {
	Status  status      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type SystemHandler struct {
	// 系统启动时间
	startTime time.Time
	// 系统结束时间
	endTime      time.Time
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HandleSystemInfo 处理系统信息请求
func (h *SystemHandler) HandleSystemInfo(c *gin.Context) {
	// TODO: 实现获取系统信息的逻辑
	h.endTime = time.Now()
	// 创建 MemStats 实例并读取内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sysInfo := gin.H{
		"status":     running,
		"startTime":  h.startTime,
		"version":    version,
		"uptime":     h.endTime.Sub(h.startTime).String(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"goroutines": runtime.NumGoroutine(),
		"memory": gin.H{
			"alloc":      memStats.Alloc,      // 当前分配的内存
			"totalAlloc": memStats.TotalAlloc, // 累计分配的内存
			"sys":        memStats.Sys,        // 从系统获取的内存
			"numGC":      memStats.NumGC,      // GC 次数
			"heapAlloc":  memStats.HeapAlloc,  // 堆内存分配量
			"heapSys":    memStats.HeapSys,    // 堆内存从系统获取的量
		},
	}
	c.JSON(http.StatusOK, Response{
		Status:  running,
		Data:    sysInfo,
		Message: "ok",
	})
}

// handlePing 处理ping请求
func (h *SystemHandler) HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  pong,
		Message: "pong",
		Data:    time.Now().Format(time.RFC3339),
	})
}

// handleHealth 处理健康检查请求
func (h *SystemHandler) HandleHealth(c *gin.Context) {
	healthStatus := healthy
	// TODO: 实现健康检查逻辑
	services := map[string]status{
		"api":     healthy,
		"storage": healthy,
	}

	// 检查所有服务状态
	for _, status := range services {
		if status != healthy {
			healthStatus = stopped
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    healthStatus,
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   version,
		"services":  services,
	})
}

// handleSystemStatus 处理系统状态请求
func (h *SystemHandler) HandleSystemStatus(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  running,
		Message: "ok",
		Data: gin.H{
			"uptime": h.endTime.Sub(h.startTime).String(),
			"services": gin.H{
				"api":     healthy,
				"storage": healthy,
			},
		},
	})
}

// 添加错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  error,
				"message": c.Errors.String(),
			})
		}
	}
}
