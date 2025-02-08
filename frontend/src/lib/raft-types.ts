export type NodeState = 'follower' | 'candidate' | 'leader';

export interface RaftNode {
  id: number;
  state: NodeState;
  term: number;
  votedFor: number | null;
  log: LogEntry[];
  commitIndex: number;
  lastApplied: number;
}

export interface LogEntry {
  term: number;
  command: string;
}

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
