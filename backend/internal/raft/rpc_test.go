package raft

import (
	"testing"
	"time"
)

func TestFooRPCTransport_Basic(t *testing.T) {
	peers := map[string]bool{
		"node2": true,
		"node3": true,
	}
	transport := NewFooRPCTransport("node1", "localhost:8000", peers)

	// 测试初始状态
	if transport.LocalAddr() != "localhost:8000" {
		t.Errorf("期望地址为 localhost:8000，实际为 %s", transport.LocalAddr())
	}
}

func TestFooRPCTransport_Send(t *testing.T) {
	transport := NewFooRPCTransport("node1", "localhost:8000", nil)

	// 创建测试 RPC 消息
	rpc := RPC{
		Type: VoteRequest,
		From: "node1",
		To:   "node2",
		Args: &RequestVoteArgs{
			Term:        1,
			CandidateID: "node1",
		},
	}

	// 测试发送消息
	if err := transport.Send("node2", rpc); err != nil {
		t.Fatalf("发送 RPC 消息失败: %v", err)
	}

	// 测试接收消息
	select {
	case received := <-transport.Consumer():
		if received.Type != rpc.Type {
			t.Errorf("RPC 类型不匹配，期望 %v，实际为 %v", rpc.Type, received.Type)
		}
		if received.From != rpc.From {
			t.Errorf("发送者不匹配，期望 %s，实际为 %s", rpc.From, received.From)
		}
		if received.To != rpc.To {
			t.Errorf("接收者不匹配，期望 %s，实际为 %s", rpc.To, received.To)
		}
	case <-time.After(time.Second):
		t.Fatal("接收 RPC 消息超时")
	}
}

func TestFooRPCTransport_Close(t *testing.T) {
	transport := NewFooRPCTransport("node1", "localhost:8000", nil)

	// 测试关闭传输层
	if err := transport.Close(); err != nil {
		t.Fatalf("关闭传输层失败: %v", err)
	}

	// 确认已关闭
	if !transport.IsShutdown() {
		t.Error("传输层应该已关闭")
	}

	// 测试向已关闭的传输层发送消息
	rpc := RPC{
		Type: VoteRequest,
		From: "node1",
		To:   "node2",
	}
	if err := transport.Send("node2", rpc); err != ErrTransportShutdown {
		t.Errorf("向已关闭的传输层发送消息应该返回 ErrTransportShutdown，实际为 %v", err)
	}
}

func TestFooRPCTransport_Consumer(t *testing.T) {
	transport := NewFooRPCTransport("node1", "localhost:8000", nil)

	// 获取消费者通道
	consumer := transport.Consumer()
	if consumer == nil {
		t.Fatal("消费者通道不应为 nil")
	}

	// 测试通道缓冲区
	for i := 0; i < 64; i++ {
		rpc := RPC{
			Type: VoteRequest,
			From: "node1",
			To:   "node2",
			Args: &RequestVoteArgs{Term: uint64(i)},
		}
		if err := transport.Send("node2", rpc); err != nil {
			t.Fatalf("发送第 %d 个消息失败: %v", i, err)
		}
	}

	// 验证所有消息都能收到
	for i := 0; i < 64; i++ {
		select {
		case rpc := <-consumer:
			args := rpc.Args.(*RequestVoteArgs)
			if args.Term != uint64(i) {
				t.Errorf("消息顺序错误，期望任期为 %d，实际为 %d", i, args.Term)
			}
		case <-time.After(time.Second):
			t.Fatalf("接收第 %d 个消息超时", i)
		}
	}
}
