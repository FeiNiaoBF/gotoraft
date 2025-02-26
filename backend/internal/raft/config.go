package raft

import (
	"fmt"
	"time"
)

// Config 包含Raft节点的配置信息
type Config struct {
	// 节点标识
	NodeID   string   // 节点ID
	PeerIDs  []string // 其他节点的ID列表
	
	// 时间配置
	HeartbeatTimeout time.Duration // 心跳超时（默认100ms）
	ElectionTimeout  time.Duration // 选举超时（默认1000ms）
	
	// 可视化配置
	VisualizeEnabled bool   // 是否启用可视化
	VisualizeAddr    string // 可视化服务地址
	VisualizePort    int    // 可视化服务端口
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		HeartbeatTimeout: 100 * time.Millisecond,
		ElectionTimeout:  1000 * time.Millisecond,
		VisualizeEnabled: true,
		VisualizeAddr:   "0.0.0.0",
		VisualizePort:   8080,
	}
}

// ValidateConfig 验证配置是否有效
func ValidateConfig(c *Config) error {
	if c.NodeID == "" {
		return fmt.Errorf("NodeID cannot be empty")
	}
	if c.HeartbeatTimeout >= c.ElectionTimeout {
		return fmt.Errorf("HeartbeatTimeout must be less than ElectionTimeout")
	}
	return nil
}
