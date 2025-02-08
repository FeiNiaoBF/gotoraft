'use client';

import { useState, useEffect } from 'react';
import { type RaftNode, createInitialNode } from '@/lib/raft-types';

const ELECTION_TIMEOUT_MIN = 150;
const ELECTION_TIMEOUT_MAX = 300;

interface LeaderElectionProps {
  numNodes: number;
}

export function LeaderElection({ numNodes }: LeaderElectionProps) {
  const [nodes, setNodes] = useState<RaftNode[]>([]);

  useEffect(() => {
    // Initialize nodes
    setNodes(
      Array.from({ length: numNodes }, (_, i) => createInitialNode(i + 1))
    );
  }, [numNodes]);

  useEffect(() => {
    // Set up election timeouts for each node
    const timeouts = nodes.map((node) => {
      if (node.state !== 'leader') {
        return setTimeout(() => startElection(node.id), getRandomTimeout());
      }
      return null;
    });

    return () => {
      timeouts.forEach((timeout) => {
        if (timeout) clearTimeout(timeout);
      });
    };
  }, [nodes]);

  function getRandomTimeout() {
    return (
      Math.floor(
        Math.random() * (ELECTION_TIMEOUT_MAX - ELECTION_TIMEOUT_MIN + 1)
      ) + ELECTION_TIMEOUT_MIN
    );
  }

  function startElection(nodeId: number) {
    setNodes((prevNodes) => {
      const newNodes = [...prevNodes];
      const candidateNode = newNodes.find((n) => n.id === nodeId);
      if (candidateNode && candidateNode.state !== 'leader') {
        candidateNode.state = 'candidate';
        candidateNode.term += 1;
        candidateNode.votedFor = candidateNode.id;
        // Request votes from other nodes
        requestVotes(candidateNode, newNodes);
      }
      return newNodes;
    });
  }

  function requestVotes(candidate: RaftNode, allNodes: RaftNode[]) {
    const votes = allNodes.map((node) => {
      if (node.id === candidate.id) return true;
      if (node.term < candidate.term) {
        node.term = candidate.term;
        node.votedFor = candidate.id;
        return true;
      }
      return false;
    });

    const voteCount = votes.filter(Boolean).length;
    if (voteCount > allNodes.length / 2) {
      // Candidate wins the election
      setNodes((prevNodes) => {
        const newNodes = [...prevNodes];
        const newLeader = newNodes.find((n) => n.id === candidate.id);
        if (newLeader) {
          newLeader.state = 'leader';
          // Reset other nodes to followers
          newNodes.forEach((node) => {
            if (node.id !== newLeader.id) {
              node.state = 'follower';
              node.votedFor = null;
            }
          });
        }
        return newNodes;
      });
    }
  }

  return (
    <div className='grid grid-cols-3 gap-4'>
      {nodes.map((node) => (
        <div
          key={node.id}
          className={`p-4 rounded-lg ${
            node.state === 'leader'
              ? 'bg-emerald-600'
              : node.state === 'candidate'
              ? 'bg-yellow-600'
              : 'bg-slate-600'
          }`}>
          <h3 className='text-lg font-semibold'>Node {node.id}</h3>
          <p>State: {node.state}</p>
          <p>Term: {node.term}</p>
        </div>
      ))}
    </div>
  );
}
