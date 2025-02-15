import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useInterval } from '@/hooks/use-interval';
import { RaftLayout } from '@/layout/raft-bg-layout';


interface Node {
  id: number;
  x: number;
  y: number;
  state: NodeState;
  isAlive: boolean;
  vx: number;
  vy: number;
}

enum NodeState {
  Follower,
  Leader,
}

// 心跳粒子
interface HeartbeatParticle {
  id: string;
  fromNode: number;
  toNode: number;
  progress: number;
  scale: number;
  opacity: number;
}

// 动画配置
const ANIMATION_CONFIG = {
  HEARTBEAT_INTERVAL: 1000, // 心跳间隔
  PARTICLE_SPEED: 0.02, // 粒子移动速度
  PULSE_DURATION: 1.5, // 脉动周期
  NODE_GLOW_RADIUS: 4, // 节点发光半径
  PARTICLE_SIZE: 4, // 粒子大小
  PARTICLE_TRAIL_LENGTH: 2, // 粒子拖尾长度
} as const;

export default function RaftVisualization() {

  // 节点状态
  const [nodes, setNodes] = useState<Node[]>([]);
  // 鼠标悬停的节点
  const [hoveredNode, setHoveredNode] = useState<number | null>(null);
  // 心跳粒子列表
  const [heartbeatParticles, setHeartbeatParticles] = useState<
    HeartbeatParticle[]
  >([]);
  // 初始化节点
  useEffect(() => {
    const totalNodes = 5;
    const newNodes = Array.from({ length: totalNodes }, (_, i) => ({
      id: i + 1,
      x: Math.random() * 400 + 50, // Random x between 50 and 450
      y: Math.random() * 300 + 150, // Random y between 150 and 450
      state: i === 0 ? NodeState.Leader : NodeState.Follower,
      isAlive: true,
      vx: (Math.random() - 0.5) * 2, // Random velocity between -1 and 1
      vy: (Math.random() - 0.5) * 2,
    }));
    setNodes(newNodes);
  }, []);

  // 更新节点位置
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
    const leader = nodes.find((n) => n.state === NodeState.Leader);
    if (leader) {
      const followers = nodes.filter(
        (n) => n.state === NodeState.Follower && n.isAlive
      );
      followers.forEach((follower) => {
        const newParticle = {
          id: `${Date.now()}-${follower.id}`,
          fromNode: leader.id,
          toNode: follower.id,
          progress: 0,
          scale: 1,
          opacity: 1,
        };
        setHeartbeatParticles((prev) => [...prev, newParticle]);
      });
    }
  }, ANIMATION_CONFIG.HEARTBEAT_INTERVAL);

  // Update heartbeat particles
  useInterval(() => {
    setHeartbeatParticles((prev) => {
      return prev
        .map((particle) => {
          // 计算新的进度
          const newProgress =
            particle.progress + ANIMATION_CONFIG.PARTICLE_SPEED;

          // 计算脉动效果
          const pulsePhase =
            (newProgress * Math.PI * 2) / ANIMATION_CONFIG.PULSE_DURATION;
          const scale = 0.8 + 0.4 * Math.sin(pulsePhase);

          // 计算透明度
          const opacity =
            newProgress < 0.1
              ? newProgress * 10 // 淡入
              : newProgress > 0.9
                ? (1 - newProgress) * 10 // 淡出
                : 1; // 完全不透明

          return {
            ...particle,
            progress: newProgress,
            scale,
            opacity,
          };
        })
        .filter((particle) => particle.progress <= 1);
    });
  }, 50);

  const getNodeColor = (node: Node) => {
    if (!node.isAlive) return '#ef4444';
    switch (node.state) {
      case NodeState.Leader:
        return 'url(#leaderGradient)';
      default:
        return 'url(#followerGradient)';
    }
  };

  const handleNodeClick = (id: number) => {
    setNodes((prevNodes) => {
      const clickedNode = prevNodes.find((node) => node.id === id);
      if (!clickedNode) return prevNodes;

      const newState =
        clickedNode.state === NodeState.Follower
          ? NodeState.Leader
          : NodeState.Follower;

      return prevNodes.map((node) => {
        if (node.id === id) {
          return { ...node, state: newState };
        }
        // If this node becomes leader, all other nodes become followers
        if (newState === NodeState.Leader && node.id !== id) {
          return { ...node, state: NodeState.Follower };
        }
        return node;
      });
    });
  };

  return (
    <RaftLayout
      starDensity={0.5}
      withGalaxyEffect={false}>
      <div className='relative w-full h-full flex flex-col'>
        {/* Add gradient background */}
        <div className='absolute inset-0 bg-gradient-radial from-slate-900 via-emerald-900/20 to-slate-900 rounded-3xl' />
        <div className='absolute inset-0 backdrop-blur-[100px]' />

        <div className='relative flex-1 flex items-center justify-center'>
          <svg
            className='w-full h-full max-w-[800px] max-h-[600px]'
            viewBox='0 0 500 500'
            preserveAspectRatio='xMidYMid meet'
            style={{ minHeight: '500px' }}>
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

              {/* 添加粒子发光渐变 */}
              <radialGradient id='particleGradient'>
                <stop
                  offset='0%'
                  stopColor='#10b981'
                  stopOpacity='1'
                />
                <stop
                  offset='50%'
                  stopColor='#10b981'
                  stopOpacity='0.5'
                />
                <stop
                  offset='100%'
                  stopColor='#10b981'
                  stopOpacity='0'
                />
              </radialGradient>
            </defs>

            {/* Connections */}
            {nodes.map((source, i) =>
              nodes.slice(i + 1).map((target) => (
                <motion.path
                  key={`${source.id}-${target.id}`}
                  d={`M${source.x},${source.y} L${target.x},${target.y}`}
                  stroke={
                    source.state === NodeState.Leader ||
                      target.state === NodeState.Leader
                      ? 'url(#leaderGradient)'
                      : '#475569'
                  }
                  strokeWidth={1.5}
                  strokeDasharray={
                    source.state === NodeState.Leader ||
                      target.state === NodeState.Leader
                      ? '0'
                      : '4,4'
                  }
                  initial={{ pathLength: 0, opacity: 0 }}
                  animate={{
                    pathLength: 1,
                    opacity: 0.3,
                    strokeDashoffset:
                      source.state === NodeState.Leader ||
                        target.state === NodeState.Leader
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

              // 计算粒子位置
              const x = fromNode.x + (toNode.x - fromNode.x) * particle.progress;
              const y = fromNode.y + (toNode.y - fromNode.y) * particle.progress;

              // 创建拖尾效果
              const trail = Array.from(
                { length: ANIMATION_CONFIG.PARTICLE_TRAIL_LENGTH },
                (_, i) => {
                  const trailProgress = Math.max(0, particle.progress - i * 0.05);
                  const trailX =
                    fromNode.x + (toNode.x - fromNode.x) * trailProgress;
                  const trailY =
                    fromNode.y + (toNode.y - fromNode.y) * trailProgress;
                  const trailOpacity =
                    particle.opacity *
                    (1 - i / ANIMATION_CONFIG.PARTICLE_TRAIL_LENGTH);

                  return (
                    <motion.circle
                      key={`${particle.id}-trail-${i}`}
                      cx={trailX}
                      cy={trailY}
                      r={
                        ANIMATION_CONFIG.PARTICLE_SIZE *
                        (1 - i / ANIMATION_CONFIG.PARTICLE_TRAIL_LENGTH)
                      }
                      fill='#10b981'
                      opacity={trailOpacity * 0.3}
                      filter='url(#particleGlow)'
                    />
                  );
                }
              );

              return (
                <g key={particle.id}>
                  {trail}
                  <motion.circle
                    cx={x}
                    cy={y}
                    r={ANIMATION_CONFIG.PARTICLE_SIZE}
                    fill='#10b981'
                    filter='url(#particleGlow)'
                    style={{
                      opacity: particle.opacity,
                      scale: particle.scale,
                    }}
                  />
                </g>
              );
            })}

            {/* Nodes */}
            {nodes.map((node) => (
              <g
                key={node.id}
                transform={`translate(${node.x},${node.y})`}>
                <AnimatePresence>
                  {node.state === NodeState.Leader && (
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
                  stroke={node.state === NodeState.Leader ? '#10b981' : '#475569'}
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
        </div>
      </div>
    </RaftLayout>
  );
}
