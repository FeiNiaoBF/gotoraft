import { useState, useCallback } from 'react';

interface Node {
  id: number;
  x: number;
  y: number;
  state: 'follower' | 'candidate' | 'leader';
  isAlive: boolean;
}

const INITIAL_NODES: Node[] = [
  { id: 1, x: 400, y: 150, state: 'leader', isAlive: true },
  { id: 2, x: 200, y: 300, state: 'follower', isAlive: true },
  { id: 3, x: 600, y: 300, state: 'follower', isAlive: true },
  { id: 4, x: 300, y: 450, state: 'follower', isAlive: true },
  { id: 5, x: 500, y: 450, state: 'follower', isAlive: true },
];

export function useRaftNodes() {
  const [nodes, setNodes] = useState<Node[]>(INITIAL_NODES);

  const updateNodeState = useCallback((id: number, state: Node['state']) => {
    setNodes((prev) =>
      prev.map((node) => (node.id === id ? { ...node, state } : node))
    );
  }, []);

  const toggleNodeStatus = useCallback((id: number) => {
    setNodes((prev) =>
      prev.map((node) =>
        node.id === id ? { ...node, isAlive: !node.isAlive } : node
      )
    );
  }, []);

  return {
    nodes,
    updateNodeState,
    toggleNodeStatus,
  };
}
