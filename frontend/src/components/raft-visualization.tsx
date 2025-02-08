'use client';

import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useInterval } from '@/hooks/use-interval';
import { LeaderElection } from '@/components/leader-election';

interface Node {
  id: number;
  x: number;
  y: number;
  state: 'follower' | 'leader';
  isAlive: boolean;
  vx: number;
  vy: number;
}

interface HeartbeatParticle {
  id: string;
  fromNode: number;
  toNode: number;
  progress: number;
}

const STATES = {
  NORMAL: 'normal',
  ELECTION: 'election',
  FAILURE: 'failure',
} as const;

export default function RaftVisualization() {
  const [currentState, setCurrentState] =
    useState<keyof typeof STATES>('NORMAL');
  const [nodes, setNodes] = useState<Node[]>([]);
  const [hoveredNode, setHoveredNode] = useState<number | null>(null);
  const [heartbeatParticles, setHeartbeatParticles] = useState<
    HeartbeatParticle[]
  >([]);

  // Initialize nodes with random positions and velocities
  useEffect(() => {
    const totalNodes = 5;
    const newNodes = Array.from({ length: totalNodes }, (_, i) => ({
      id: i + 1,
      x: Math.random() * 400 + 50, // Random x between 50 and 450
      y: Math.random() * 400 + 50, // Random y between 50 and 450
      state: i === 0 ? 'leader' : 'follower',
      isAlive: true,
      vx: (Math.random() - 0.5) * 2, // Random velocity between -1 and 1
      vy: (Math.random() - 0.5) * 2,
    }));
    setNodes(newNodes);
  }, []);

  // Update node positions
  useInterval(() => {
    setNodes((prevNodes) =>
      prevNodes.map((node) => {
        let newX = node.x + node.vx;
        let newY = node.y + node.vy;
        let newVx = node.vx;
        let newVy = node.vy;

        // More strict boundary checking
        const padding = 30; // Ensure nodes stay within visible area
        if (newX < padding || newX > 500 - padding) {
          newVx = -newVx * 0.8; // Add some damping
          newX = newX < padding ? padding : 500 - padding;
        }
        if (newY < padding || newY > 500 - padding) {
          newVy = -newVy * 0.8; // Add some damping
          newY = newY < padding ? padding : 500 - padding;
        }

        return {
          ...node,
          x: newX,
          y: newY,
          vx: newVx,
          vy: newVy,
        };
      })
    );
  }, 50);

  // Heartbeat animation
  useInterval(() => {
    if (currentState === 'NORMAL') {
      const leader = nodes.find((n) => n.state === 'leader');
      if (leader) {
        const followers = nodes.filter(
          (n) => n.state === 'follower' && n.isAlive
        );
        followers.forEach((follower) => {
          const newParticle = {
            id: `${Date.now()}-${follower.id}`,
            fromNode: leader.id,
            toNode: follower.id,
            progress: 0,
          };
          setHeartbeatParticles((prev) => [...prev, newParticle]);
        });
      }
    }
  }, 1000);

  // Update heartbeat particles
  useInterval(() => {
    setHeartbeatParticles((prev) => {
      return prev
        .map((particle) => ({
          ...particle,
          progress: particle.progress + 0.1,
        }))
        .filter((particle) => particle.progress <= 1);
    });
  }, 50);

  // Cycle through states periodically
  useInterval(() => {
    setCurrentState((current) => {
      switch (current) {
        case 'NORMAL':
          return 'ELECTION';
        case 'ELECTION':
          return 'FAILURE';
        case 'FAILURE':
          return 'NORMAL';
        default:
          return 'NORMAL';
      }
    });
  }, 8000);

  const getNodeColor = (node: Node) => {
    if (!node.isAlive) return '#ef4444';
    switch (node.state) {
      case 'leader':
        return 'url(#leaderGradient)';
      default:
        return 'url(#followerGradient)';
    }
  };

  const handleNodeClick = (id: number) => {
    setNodes((prevNodes) => {
      const clickedNode = prevNodes.find((node) => node.id === id);
      if (!clickedNode) return prevNodes;

      const newState = clickedNode.state === 'follower' ? 'leader' : 'follower';

      return prevNodes.map((node) => {
        if (node.id === id) {
          return { ...node, state: newState };
        }
        // If this node becomes leader, all other nodes become followers
        if (newState === 'leader' && node.id !== id) {
          return { ...node, state: 'follower' };
        }
        return node;
      });
    });
  };

  return (
    <div className='relative w-full h-full flex items-center justify-center'>
      {/* Add gradient background */}
      <div className='absolute inset-0 bg-gradient-radial from-slate-900 via-emerald-900/20 to-slate-900 rounded-3xl' />
      <div className='absolute inset-0 backdrop-blur-[100px]' />

      <svg
        className='relative w-4/5 h-4/5'
        viewBox='0 0 500 500'>
        <defs>
          <linearGradient
            id='leaderGradient'
            x1='0%'
            y1='0%'
            x2='100%'
            y2='100%'>
            <stop
              offset='0%'
              stopColor='#10b981'
            />
            <stop
              offset='100%'
              stopColor='#059669'
            />
          </linearGradient>
          <linearGradient
            id='followerGradient'
            x1='0%'
            y1='0%'
            x2='100%'
            y2='100%'>
            <stop
              offset='0%'
              stopColor='#64748b'
            />
            <stop
              offset='100%'
              stopColor='#475569'
            />
          </linearGradient>

          <filter id='glow'>
            <feGaussianBlur
              stdDeviation='4'
              result='coloredBlur'
            />
            <feMerge>
              <feMergeNode in='coloredBlur' />
              <feMergeNode in='SourceGraphic' />
            </feMerge>
          </filter>

          <filter id='particleGlow'>
            <feGaussianBlur
              stdDeviation='2'
              result='coloredBlur'
            />
            <feMerge>
              <feMergeNode in='coloredBlur' />
              <feMergeNode in='SourceGraphic' />
            </feMerge>
          </filter>
        </defs>

        {/* Connections */}
        {nodes.map((source, i) =>
          nodes.slice(i + 1).map((target, j) => (
            <motion.path
              key={`${source.id}-${target.id}`}
              d={`M${source.x},${source.y} L${target.x},${target.y}`}
              stroke={
                source.state === 'leader' || target.state === 'leader'
                  ? 'url(#leaderGradient)'
                  : '#475569'
              }
              strokeWidth={1.5}
              strokeDasharray={
                source.state === 'leader' || target.state === 'leader'
                  ? '0'
                  : '4,4'
              }
              initial={{ pathLength: 0, opacity: 0 }}
              animate={{
                pathLength: 1,
                opacity: 0.3,
                strokeDashoffset:
                  source.state === 'leader' || target.state === 'leader'
                    ? 0
                    : -20,
              }}
              transition={{
                pathLength: { duration: 0.5 },
                opacity: { duration: 0.3 },
                strokeDashoffset: {
                  repeat: Number.POSITIVE_INFINITY,
                  duration: 2,
                  ease: 'linear',
                },
              }}
            />
          ))
        )}

        {/* Heartbeat particles */}
        {heartbeatParticles.map((particle) => {
          const fromNode = nodes.find((n) => n.id === particle.fromNode);
          const toNode = nodes.find((n) => n.id === particle.toNode);
          if (!fromNode || !toNode) return null;

          const x = fromNode.x + (toNode.x - fromNode.x) * particle.progress;
          const y = fromNode.y + (toNode.y - fromNode.y) * particle.progress;

          return (
            <motion.circle
              key={particle.id}
              cx={x}
              cy={y}
              r={4}
              fill='#10b981'
              filter='url(#particleGlow)'
              initial={{ opacity: 0, scale: 0 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0 }}
              transition={{ duration: 0.3 }}
            />
          );
        })}

        {/* Nodes */}
        {nodes.map((node) => (
          <g
            key={node.id}
            transform={`translate(${node.x},${node.y})`}>
            <AnimatePresence>
              {node.state === 'leader' && (
                <motion.circle
                  cx={0}
                  cy={0}
                  r={25}
                  fill='none'
                  stroke='url(#leaderGradient)'
                  strokeWidth={2}
                  initial={{ scale: 1, opacity: 0 }}
                  animate={{ scale: 1.5, opacity: 0.2 }}
                  exit={{ scale: 1, opacity: 0 }}
                  transition={{
                    repeat: Number.POSITIVE_INFINITY,
                    duration: 2,
                    ease: 'easeInOut',
                  }}
                />
              )}
            </AnimatePresence>

            <motion.circle
              cx={0}
              cy={0}
              r={20}
              fill={getNodeColor(node)}
              stroke={node.state === 'leader' ? '#10b981' : '#475569'}
              strokeWidth={2}
              initial={{ scale: 0 }}
              animate={{
                scale: 1,
                filter: hoveredNode === node.id ? 'url(#glow)' : 'none',
              }}
              whileHover={{ scale: 1.1 }}
              onMouseEnter={() => setHoveredNode(node.id)}
              onMouseLeave={() => setHoveredNode(null)}
              onClick={() => handleNodeClick(node.id)}
              transition={{ duration: 0.3 }}
              style={{ cursor: 'pointer' }}
            />

            <text
              x={0}
              y={0}
              textAnchor='middle'
              dy='.3em'
              className='text-xs font-medium fill-white'
              pointerEvents='none'>
              {node.id}
            </text>
          </g>
        ))}
      </svg>
      <LeaderElection numNodes={5} />
    </div>
  );
}
