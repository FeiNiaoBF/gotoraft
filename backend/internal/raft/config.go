// from MIT 6.5840
// config.go 是raft的配置文件，用来配置raft的参数和建立上下文
// 可以配置raft的节点数量，选举超时时间，心跳时间，以及raft的日志文件
// 初始化raft的配置文件
package raft

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"gotoraft/internal/foorpc"
	// 如果你有自己的 registry / xclient，需要相应 import
	// "gotoraft/internal/foorpc/registry"
	// "gotoraft/internal/foorpc/xclient"
)

// Config 用于存储 Raft 的关键配置
type Config struct {
	mu                sync.Mutex    // 互斥锁
	NodeCount         int           // 节点数量
	ElectionTimeout   time.Duration // 选举超时时间
	HeartbeatInterval time.Duration // 心跳间隔
	LogFile           string        // 日志文件路径，示例用途
	persister         *Persister    // 持久化器
	peers             []string      // 节点列表

}

// NewConfig 创建默认配置的工厂方法，可根据需要扩展
func NewConfig() *Config {
	return &Config{
		NodeCount:         3,
		ElectionTimeout:   1500 * time.Millisecond,
		HeartbeatInterval: 100 * time.Millisecond,
		LogFile:           "./raft.log",
	}
}

// SetupRaftEnv 根据配置来初始化整个 Raft "上下文" 环境
func (cfg *Config) SetupRaftEnv() error {
	// 1. 初始化每个节点并注册到 RPC 服务
	for i := 0; i < cfg.NodeCount; i++ {
		// 示例：创建一个监听端口
		addr := fmt.Sprintf(":%d", 8000+i)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %v", addr, err)
		}

		// 2. 启动 RPC Server 并注册服务
		go func(i int, ln net.Listener) {
			// 可以使用你定义的 foorpc.NewServer() 来创建 RPCServer
			rpcServer := foorpc.NewServer()

			// 如果你有 "RaftService" 之类的结构来承载节点逻辑，可以在这里注册
			// s := NewRaftService(node)
			// _ = rpcServer.Register(s)

			// 开始接受 RPC 连接
			rpcServer.Accept(ln)
			log.Printf("Node %d started RPC server on %s\n", i, addr)
		}(i, ln)
	}

	// 3. 设置选举超时时间、心跳间隔等参数
	// 你可以在 Raft 节点初始化时，把这些值赋进去
	// 示例：
	// for _, node := range allRaftNodes {
	//     node.SetElectionTimeout(cfg.ElectionTimeout)
	//     node.SetHeartbeatInterval(cfg.HeartbeatInterval)
	// }

	// 4. 如果需要可以在这里做日志文件初始化
	// f, err := os.OpenFile(cfg.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	//     return fmt.Errorf("failed to open log file %s: %v", cfg.LogFile, err)
	// }
	// log.SetOutput(f)

	// 5. 如果需要注册到 registry 或者启动心跳，可在这里实现
	// registry.HandleHTTP()
	// registry.Heartbeat("http://localhost:9999/_foorpc_/registry", addr, 0)

	return nil
}
