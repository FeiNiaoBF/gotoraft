// 该文件用来实现Raft节点之间的RPC通信
package raft

import (
	"fmt"
)

// FooRPCTransport 是基于 foorpc 的传输层实现
type FooRPCTransport struct {
	id      string          // 节点 ID
	addr    string          // 节点地址
	peers   map[string]bool // 对等节点列表
	rpcCh   chan RPC        // RPC 通道
	closeCh chan struct{}   // 关闭通道
}

// NewFooRPCTransport 创建一个新的 FooRPC 传输层
func NewFooRPCTransport(id, addr string, peers map[string]bool) Transport {
	return &FooRPCTransport{
		id:      id,
		addr:    addr,
		peers:   peers,
		rpcCh:   make(chan RPC, 64),
		closeCh: make(chan struct{}),
	}
}

// Send 发送 RPC 请求
func (t *FooRPCTransport) Send(target string, rpc RPC) error {
	if t.IsShutdown() {
		return ErrTransportShutdown
	}

	select {
	case <-t.closeCh:
		return ErrTransportShutdown
	case t.rpcCh <- rpc:
		return nil
	default:
		return fmt.Errorf("transport channel is full")
	}
}

// Consumer 返回 RPC 通道
func (t *FooRPCTransport) Consumer() <-chan RPC {
	return t.rpcCh
}

// Close 关闭传输层
func (t *FooRPCTransport) Close() error {
	select {
	case <-t.closeCh:
		return nil
	default:
		close(t.closeCh)
	}
	return nil
}

// LocalAddr 返回本地地址
func (t *FooRPCTransport) LocalAddr() string {
	return t.addr
}

// IsShutdown 检查传输层是否已关闭
func (t *FooRPCTransport) IsShutdown() bool {
	select {
	case <-t.closeCh:
		return true
	default:
		return false
	}
}

// ErrTransportShutdown 表示传输层已关闭的错误
var ErrTransportShutdown = fmt.Errorf("transport shutdown")
