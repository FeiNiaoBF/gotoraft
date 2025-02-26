// internal/handler/kv_handler.go
package handler

import (
	"gotoraft/internal/kvstore/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// KVStoreHandler 处理KV存储的请求
type KVStoreHandler struct {
	store *store.Store
}

// NewKVStoreHandler 创建一个新的KV存储处理器
func NewKVStoreHandler(store *store.Store) *KVStoreHandler {
	return &KVStoreHandler{
		store: store,
	}
}

// HandleGet 处理获取键值的请求
func (h *KVStoreHandler) HandleGet(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
		})
		return
	}

	value, err := h.store.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"key":   key,
			"value": value,
		},
	})
}

// HandleSet 处理设置键值的请求
type SetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func (h *KVStoreHandler) HandleSet(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.store.Set(req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to set value: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"key":   req.Key,
			"value": req.Value,
		},
	})
}

// HandleDelete 处理删除键值的请求
func (h *KVStoreHandler) HandleDelete(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
		})
		return
	}

	if err := h.store.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete key: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Key deleted successfully",
		"data": gin.H{
			"key": key,
		},
	})
}
