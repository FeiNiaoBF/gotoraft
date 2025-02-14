import { useState, useCallback } from 'react';
import { motion } from 'framer-motion';
import { RaftNode, NodeState, createInitialNode } from '@/lib/raft-types';
import { useInterval } from '@/hooks/use-interval';
import { RaftLayout } from '@/layout/raft-bg-layout';
import { useLanguage } from '@/contexts/language-context';

// Raft可视化配置
const RAFT_CONFIG = {
  OFFSET_X: 380, // 右侧偏移量
  OFFSET_Y: 400, // 上侧偏移量
  BASE_RADIUS: 250, // 基础轨道半径
  NODE_SIZE: 64, // 节点直径
  ORBIT_SPEED_RANGE: [0.0003, 0.0006], // 轨道速度范围
  SELF_ROTATION: 0.0004, // 自转速度
  HEARTBEAT_INTERVAL: 800, // 心跳间隔
  STAR_POINTS: 5, // 五角星的点数
  STAR_INNER_RATIO: 0.382, // 五角星内圈比例（黄金分割比）
};
// 节点可视化参数类型
interface VisualNode extends RaftNode {
  x: number;
  y: number;
  orbitRadius: number; // 轨道半径
  orbitAngle: number; // 轨道角度
  rotation: number; // 旋转角度
  orbitSpeed: number; // 轨道速度
}

export default function RaftVisualization() {
  const { t } = useLanguage();

  // 初始化节点
  const [nodes, setNodes] = useState<VisualNode[]>(() =>
    [1, 2, 3, 4, 5].map((id, index) => {
      const isLeader = id === 1;
      const angle = (index * 72 * Math.PI) / 180; // 均匀分布在圆上，每个节点间隔72度
      return {
        ...createInitialNode(id),
        state: isLeader ? 'leader' : 'follower',
        x: RAFT_CONFIG.OFFSET_X + Math.cos(angle) * RAFT_CONFIG.BASE_RADIUS,
        y: RAFT_CONFIG.OFFSET_Y + Math.sin(angle) * RAFT_CONFIG.BASE_RADIUS,
        orbitAngle: angle,
        orbitRadius: RAFT_CONFIG.BASE_RADIUS,
        rotation: Math.random() * Math.PI * 2,
        orbitSpeed:
          RAFT_CONFIG.ORBIT_SPEED_RANGE[0] +
          Math.random() *
            (RAFT_CONFIG.ORBIT_SPEED_RANGE[1] -
              RAFT_CONFIG.ORBIT_SPEED_RANGE[0]),
      };
    })
  );

  // 更新节点位置
  const updateOrbitPositions = useCallback(() => {
    setNodes((prev) =>
      prev.map((node) => {
        const newAngle = node.orbitAngle + node.orbitSpeed;
        return {
          ...node,
          x: RAFT_CONFIG.OFFSET_X + Math.cos(newAngle) * node.orbitRadius,
          y: RAFT_CONFIG.OFFSET_Y + Math.sin(newAngle) * node.orbitRadius,
          orbitAngle: newAngle,
          rotation: node.rotation + RAFT_CONFIG.SELF_ROTATION,
        };
      })
    );
  }, []);

  // 心跳间隔
  const [heartbeats, setHeartbeats] = useState<
    Array<{
      id: string;
      from: number;
      to: number;
      progress: number;
    }>
  >([]);

  // 生成心跳信号
  const generateHeartbeats = useCallback(() => {
    const leader = nodes.find((n) => n.state === 'leader');
    if (!leader) return;

    const followers = nodes.filter((n) => n.state === 'follower');
    const newParticles = followers.map((follower) => ({
      id: `${performance.now()}-${leader.id}-${follower.id}`,
      from: leader.id,
      to: follower.id,
      progress: 0,
    }));

    setHeartbeats((prev) => [...prev.slice(-30), ...newParticles]);
  }, [nodes]);

  // 心跳粒子动画
  useInterval(() => {
    setHeartbeats((prev) =>
      prev
        .map((p) => ({ ...p, progress: p.progress + 0.02 }))
        .filter((p) => p.progress <= 1)
    );
  }, 50);
  useInterval(updateOrbitPositions, 16);
  useInterval(generateHeartbeats, RAFT_CONFIG.HEARTBEAT_INTERVAL);

  // 处理节点点击
  const handleNodeClick = (id: number) => {
    setNodes((prev) =>
      prev.map((node) => {
        if (node.id === id) {
          const newState: NodeState =
            node.state === 'follower' ? 'leader' : 'follower';
          return { ...node, state: newState };
        }
        // 如果其他节点是leader，则降级为follower
        return node.state === 'leader' ? { ...node, state: 'follower' } : node;
      })
    );
  };

  return (
    <RaftLayout
      starDensity={0.8}
      withGalaxyEffect={false}>
      <div className='relative w-full h-full'>
        <svg
          className='absolute inset-0 w-full h-full pointer-events-none'
          style={{ zIndex: 1 }}>
          {/* 渲染圆形轨道 */}
          <circle
            cx={RAFT_CONFIG.OFFSET_X + RAFT_CONFIG.NODE_SIZE / 2}
            cy={RAFT_CONFIG.OFFSET_Y + RAFT_CONFIG.NODE_SIZE / 2}
            r={RAFT_CONFIG.BASE_RADIUS}
            stroke='rgba(100, 149, 237, 0.1)'
            strokeWidth='1'
            fill='none'
          />
          {/* 渲染心跳粒子 */}
          {heartbeats.map((particle) => {
            const fromNode = nodes.find((n) => n.id === particle.from);
            const toNode = nodes.find((n) => n.id === particle.to);
            if (!fromNode || !toNode) return null;

            const x1 = fromNode.x;
            const y1 = fromNode.y;
            const x2 = toNode.x;
            const y2 = toNode.y;

            const dx = x2 - x1;
            const dy = y2 - y1;
            const x = x1 + dx * particle.progress;
            const y = y1 + dy * particle.progress;

            return (
              <circle
                key={particle.id}
                cx={x}
                cy={y}
                r={3}
                fill='#64B5F6'
                opacity={1 - particle.progress}
              />
            );
          })}
        </svg>

        {/* 渲染节点 */}
        {nodes.map((node) => {
          const isLeader = node.state === 'leader';
          const leaderGradient =
            'bg-gradient-to-br from-emerald-400 to-emerald-600';
          const followerGradient =
            'bg-gradient-to-br from-slate-500 to-slate-700';

          return (
            <motion.div
              key={node.id}
              className={`absolute flex items-center justify-center rounded-full cursor-pointer
              shadow-xl transition-colors duration-300 ${
                isLeader ? leaderGradient : followerGradient
              }`}
              style={{
                width: RAFT_CONFIG.NODE_SIZE,
                height: RAFT_CONFIG.NODE_SIZE,
                x: node.x - RAFT_CONFIG.NODE_SIZE / 2,
                y: node.y - RAFT_CONFIG.NODE_SIZE / 2,
                rotate: node.rotation,
              }}
              whileHover={{ scale: 1.15 }}
              whileTap={{ scale: 0.95 }}
              transition={{ type: 'spring', stiffness: 300 }}
              onClick={() => handleNodeClick(node.id)}>
              <span className='text-white font-medium text-sm z-10'>
                {node.id}
              </span>

              {/* 领导者光环 */}
              {isLeader && (
                <motion.div
                  className='absolute inset-0 rounded-full border-2 border-emerald-400/30'
                  animate={{
                    scale: [1, 1.4],
                    opacity: [0.3, 0],
                  }}
                  transition={{
                    duration: 2,
                    repeat: Infinity,
                    ease: 'easeInOut',
                  }}
                />
              )}
            </motion.div>
          );
        })}
      </div>
      {/* 说明文字 */}
      <div className='absolute bottom-8 left-1/2 -translate-x-1/2 text-center'>
        <p className='text-sm text-emerald-100/90'>{t('nodeInstruction')}</p>
        <p className='text-xs text-slate-400/80 mt-1'>{t('nodeStates')}</p>
      </div>
    </RaftLayout>
  );
}
