package raft

import (
	"fmt"
	"gotoraft/internal/foorpc"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RaftController struct {
	mu      sync.Mutex
	rafts   []*Raft
	n       int    // 节点数量
	wsConns map[string]*websocket.Conn  // WebSocket 连接池
}

// RaftState 用于传输到前端的状态信息
type RaftState struct {
	NodeID      int    `json:"nodeId"`
	State       State  `json:"state"`
	CurrentTerm int    `json:"currentTerm"`
	VotedFor    int    `json:"votedFor"`
	LogLength   int    `json:"logLength"`
	IsLeader    bool   `json:"isLeader"`
}

func NewController(n int) *RaftController {
	rc := &RaftController{
		n:       n,
		wsConns: make(map[string]*websocket.Conn),
	}
	rc.init()
	return rc
}

func (rc *RaftController) init() {
	// 初始化 Raft 节点
	for i := 0; i < rc.n; i++ {
		// 创建 Raft 节点
		ends := make([]*RPCEnd, rc.n)
		for j := 0; j < rc.n; j++ {
			client, _ := foorpc.Dial("tcp", fmt.Sprintf(":%d", basePort+j))
			ends[j] = &RPCEnd{client: &RPCClient{Client: client}}
		}
		applyCh := make(chan ApplyMsg)
		rf := Make(ends, i, MakePersister(), applyCh)
		rc.rafts = append(rc.rafts, rf)

		// 启动状态监控
		go rc.monitorState(i)
	}
}

func (rc *RaftController) monitorState(i int) {
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		rf := rc.rafts[i]
		if rf == nil {
			continue
		}

		state := RaftState{
			NodeID:      i,
			State:       rf.state,
			CurrentTerm: rf.currentTerm,
			VotedFor:    rf.votedFor,
			LogLength:   len(rf.logs),
			IsLeader:    rf.state == Leader,
		}

		// 广播状态到所有 WebSocket 连接
		rc.broadcastState(state)
	}
}
