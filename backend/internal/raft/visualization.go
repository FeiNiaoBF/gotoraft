// 该文件用来配置Raft的可视化
package raft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 事件类型
type EventType string

const (
	EventStateChange  EventType = "state_change"  // 状态变更事件
	EventVoteRequest EventType = "vote_request"  // 请求投票事件
	EventVoteGranted EventType = "vote_granted"  // 投票授予事件
	EventVoteDenied  EventType = "vote_denied"   // 投票拒绝事件
	EventLogAppend   EventType = "log_append"    // 日志追加事件
	EventLogCommit   EventType = "log_commit"    // 日志提交事件
	EventLeaderSync  EventType = "leader_sync"   // 领导者同步事件
)

// Event 表示一个事件
type Event struct {
	Type    EventType   `json:"type"`    // 事件类型
	NodeID  string      `json:"node_id"` // 节点ID
	Payload interface{} `json:"payload"` // 事件数据
}

// StateChangePayload 状态变更事件的数据
type StateChangePayload struct {
	OldState RaftState `json:"old_state"` // 旧状态
	NewState RaftState `json:"new_state"` // 新状态
	Term     uint64    `json:"term"`      // 当前任期
}

// VoteRequestPayload 请求投票事件的数据
type VoteRequestPayload struct {
	Term         uint64 `json:"term"`          // 候选人的任期号
	CandidateID  string `json:"candidate_id"`  // 请求选票的候选人的 ID
	LastLogIndex uint64 `json:"last_log_index"` // 候选人的最后日志条目的索引值
	LastLogTerm  uint64 `json:"last_log_term"`  // 候选人最后日志条目的任期号
}

// VoteResponsePayload 投票响应事件的数据
type VoteResponsePayload struct {
	Term        uint64 `json:"term"`         // 当前任期号
	VoteGranted bool   `json:"vote_granted"` // 是否投票
}

// LogAppendPayload 日志追加事件的数据
type LogAppendPayload struct {
	Term    uint64    `json:"term"`     // 任期号
	Index   uint64    `json:"index"`    // 日志索引
	Command interface{} `json:"command"` // 命令内容
}

// LogCommitPayload 日志提交事件的数据
type LogCommitPayload struct {
	Index uint64 `json:"index"` // 提交的日志索引
}

// LeaderSyncPayload 领导者同步事件的数据
type LeaderSyncPayload struct {
	Term      uint64            `json:"term"`       // 任期号
	LeaderID  string            `json:"leader_id"`  // 领导者ID
	NextIndex map[string]uint64 `json:"next_index"` // 下一个要发送的日志索引
}

// Visualizer 负责可视化 Raft 节点的状态
type Visualizer struct {
	// 基本配置
	addr      string
	port      int
	endpoint  string

	// WebSocket相关
	upgrader  websocket.Upgrader
	clients   sync.Map
	clientsMu sync.RWMutex

	// 事件通道
	eventCh   chan Event
	stopCh    chan struct{}

	// 状态缓存
	stateMu   sync.RWMutex
	nodeState map[string]*NodeState
}

// NodeState 表示节点状态
type NodeState struct {
	ID        string    `json:"id"`
	State     RaftState `json:"state"`
	Term      uint64    `json:"term"`
	Leader    string    `json:"leader"`
	LastSeen  time.Time `json:"last_seen"`
	LogLength uint64    `json:"log_length"`
	Peers     []string  `json:"peers"`
}

// NewVisualizer 创建一个新的可视化器
func NewVisualizer(addr string, port int) *Visualizer {
	v := &Visualizer{
		addr:     addr,
		port:     port,
		endpoint: fmt.Sprintf("%s:%d", addr, port),
		clients:  sync.Map{},
		eventCh:  make(chan Event, 1024),
		stopCh:   make(chan struct{}),
		nodeState: make(map[string]*NodeState),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return v
}

// Start 启动可视化服务
func (v *Visualizer) Start() error {
	// 处理 WebSocket 连接
	http.HandleFunc("/ws", v.handleWebSocket)

	// 设置API路由
	http.HandleFunc("/api/nodes", v.handleNodes)
	http.HandleFunc("/api/state", v.handleState)
	http.HandleFunc("/api/events", v.handleEvents)

	// 启动事件处理
	go v.processEvents()

	// 启动 HTTP 服务器
	return http.ListenAndServe(v.endpoint, nil)
}

// Stop 停止可视化服务
func (v *Visualizer) Stop() {
	close(v.stopCh)
	v.clients.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*websocket.Conn); ok {
			conn.Close()
		}
		return true
	})
}

// handleWebSocket 处理 WebSocket 连接
func (v *Visualizer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := v.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 存储连接
	v.clients.Store(conn.RemoteAddr().String(), conn)

	// 清理连接
	defer func() {
		conn.Close()
		v.clients.Delete(conn.RemoteAddr().String())
	}()

	// 保持连接
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// sendCurrentState 发送当前状态给新客户端
func (v *Visualizer) sendCurrentState(conn *websocket.Conn) {
	v.stateMu.RLock()
	state := make(map[string]*NodeState)
	for id, nodeState := range v.nodeState {
		state[id] = nodeState
	}
	v.stateMu.RUnlock()

	data, err := json.Marshal(state)
	if err != nil {
		return
	}

	conn.WriteMessage(websocket.TextMessage, data)
}

// handleNodes 处理节点列表请求
func (v *Visualizer) handleNodes(w http.ResponseWriter, r *http.Request) {
	v.stateMu.RLock()
	nodes := make([]string, 0, len(v.nodeState))
	for id := range v.nodeState {
		nodes = append(nodes, id)
	}
	v.stateMu.RUnlock()

	json.NewEncoder(w).Encode(nodes)
}

// handleState 处理状态请求
func (v *Visualizer) handleState(w http.ResponseWriter, r *http.Request) {
	v.stateMu.RLock()
	state := make(map[string]*NodeState)
	for id, nodeState := range v.nodeState {
		state[id] = nodeState
	}
	v.stateMu.RUnlock()

	json.NewEncoder(w).Encode(state)
}

// handleEvents 处理事件历史请求
func (v *Visualizer) handleEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现事件历史存储和查询
	w.WriteHeader(http.StatusNotImplemented)
}

// processEvents 处理事件
func (v *Visualizer) processEvents() {
	for {
		select {
		case <-v.stopCh:
			return
		case event := <-v.eventCh:
			// 更新状态
			v.updateState(event)
			// 广播给所有客户端
			v.broadcast(event)
		}
	}
}

// updateState 更新节点状态
func (v *Visualizer) updateState(event Event) {
	v.stateMu.Lock()
	defer v.stateMu.Unlock()

	state, exists := v.nodeState[event.NodeID]
	if !exists {
		state = &NodeState{
			ID: event.NodeID,
		}
		v.nodeState[event.NodeID] = state
	}

	state.LastSeen = time.Now()
	state.State = RaftState(event.Payload.(StateChangePayload).NewState)
	state.Term = event.Payload.(StateChangePayload).Term

	if data, ok := event.Payload.(map[string]interface{}); ok {
		if logLength, exists := data["log_length"]; exists {
			if length, ok := logLength.(uint64); ok {
				state.LogLength = length
			}
		}
	}
}

// broadcast 广播事件给所有客户端
func (v *Visualizer) broadcast(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	v.clients.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*websocket.Conn); ok {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				conn.Close()
				v.clients.Delete(key)
			}
		}
		return true
	})
}

// OnStateChange 记录状态变更
func (v *Visualizer) OnStateChange(nodeID string, oldState, newState RaftState, term uint64) {
	v.eventCh <- Event{
		Type: EventStateChange,
		NodeID: nodeID,
		Payload: StateChangePayload{
			OldState: oldState,
			NewState: newState,
			Term:     term,
		},
	}
}

// OnVoteRequest 处理请求投票事件
func (v *Visualizer) OnVoteRequest(nodeID string, args *RequestVoteArgs) {
	v.eventCh <- Event{
		Type:   EventVoteRequest,
		NodeID: nodeID,
		Payload: VoteRequestPayload{
			Term:         args.Term,
			CandidateID:  args.CandidateID,
			LastLogIndex: args.LastLogIndex,
			LastLogTerm:  args.LastLogTerm,
		},
	}
}

// OnVoteResponse 处理投票响应事件
func (v *Visualizer) OnVoteResponse(nodeID string, granted bool, term uint64) {
	eventType := EventVoteGranted
	if !granted {
		eventType = EventVoteDenied
	}
	v.eventCh <- Event{
		Type:   eventType,
		NodeID: nodeID,
		Payload: VoteResponsePayload{
			Term:        term,
			VoteGranted: granted,
		},
	}
}

// OnLogAppend 处理日志追加事件
func (v *Visualizer) OnLogAppend(entry LogEntry) {
	v.eventCh <- Event{
		Type:   EventLogAppend,
		NodeID: "", // 在实际应用中设置正确的节点ID
		Payload: LogAppendPayload{
			Term:    entry.Term,
			Index:   entry.Index,
			Command: entry.Command,
		},
	}
}

// OnLogCommit 处理日志提交事件
func (v *Visualizer) OnLogCommit(index uint64) {
	v.eventCh <- Event{
		Type:   EventLogCommit,
		NodeID: "", // 在实际应用中设置正确的节点ID
		Payload: LogCommitPayload{
			Index: index,
		},
	}
}

// OnLeaderSync 处理领导者同步事件
func (v *Visualizer) OnLeaderSync(term uint64, leaderID string, nextIndex map[string]uint64) {
	v.eventCh <- Event{
		Type:   EventLeaderSync,
		NodeID: leaderID,
		Payload: LeaderSyncPayload{
			Term:      term,
			LeaderID:  leaderID,
			NextIndex: nextIndex,
		},
	}
}

type StateChangeEvent struct {
	Type     string    `json:"type"`
	NodeID   string    `json:"node_id"`
	OldState RaftState `json:"old_state"`
	NewState RaftState `json:"new_state"`
	Term     uint64    `json:"term"`
}
