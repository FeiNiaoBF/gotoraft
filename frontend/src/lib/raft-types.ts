// 首页模拟Raft的数据结构
// Raft节点状态类型
export type NodeState = 'follower' | 'candidate' | 'leader'

// 日志条目类型
export interface LogEntry {
  term: number
  command: string
}

// Raft节点类型
export interface RaftNode {
  id: number
  state: NodeState
  term: number
  votedFor: number | null
  log: LogEntry[]
  commitIndex: number;
  lastApplied: number;
}

// 创建初始节点的工厂函数
export function createInitialNode(id: number): RaftNode {
  return {
    id,
    state: 'follower',
    term: 0,
    votedFor: null,
    log: [],
    commitIndex: 0,
    lastApplied: 0,
  };
}
