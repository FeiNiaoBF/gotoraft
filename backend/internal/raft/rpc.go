package raft

import (
	"context"
	"gotoraft/internal/foorpc"
)

// RPCClient 封装RPC客户端
type RPCClient struct {
	*foorpc.Client
}

// Call 封装RPC调用
func (c *RPCClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	ctx := context.Background()
	return c.Client.Call(ctx, serviceMethod, args, reply)
}

type RPCEnd struct {
	client *RPCClient
}

func (e *RPCEnd) Call(serviceMethod string, args interface{}, reply interface{}) bool {
	if e.client == nil {
		return false
	}
	err := e.client.Call(serviceMethod, args, reply)
	return err == nil
}
