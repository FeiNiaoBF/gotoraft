// internal/handler/cluster_handler.go
package handler

import (
	"gotoraft/internal/kvstore/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ClusterHandler 处理集群管理的请求
type ClusterHandler struct {
	store *store.Store
}

// NewClusterHandler 创建一个新的集群管理处理器
func NewClusterHandler(store *store.Store) *ClusterHandler {
	return &ClusterHandler{
		store: store,
	}
}

// JoinRequest 加入集群的请求
type JoinRequest struct {
	NodeID   string `json:"nodeId" binding:"required"`
	RaftAddr string `json:"raftAddr" binding:"required"`
}

// HandleJoin 处理加入集群的请求
func (h *ClusterHandler) HandleJoin(c *gin.Context) {
	var req JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.store.Join(req.NodeID, req.RaftAddr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to join cluster: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Node joined the cluster successfully",
		"data": gin.H{
			"nodeId":   req.NodeID,
			"raftAddr": req.RaftAddr,
		},
	})
}

// LeaveRequest 离开集群的请求
type LeaveRequest struct {
	NodeID string `json:"nodeId" binding:"required"`
}

// HandleLeave 处理离开集群的请求
func (h *ClusterHandler) HandleLeave(c *gin.Context) {
	var req LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.store.Leave(req.NodeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to leave cluster: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Node left the cluster successfully",
		"data": gin.H{
			"nodeId": req.NodeID,
		},
	})
}

// HandleClusterStatus 处理获取集群状态的请求
func (h *ClusterHandler) HandleClusterStatus(c *gin.Context) {
	status, err := h.store.GetClusterStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get cluster status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   status,
	})
}
