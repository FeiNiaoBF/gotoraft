import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { RaftNode, createInitialNode } from '@/lib/raft-types';
import { useInterval } from '@/hooks/use-interval';
import { RaftLayout } from '@/layout/raft-bg-layout';

// 节点配置
const NODE_CONFIG = {
  SIZE: 60,
  ORBIT_SPEED_RANGE: [0.0003, 0.0006], // 增加速度差异
  SELF_ROTATION_SPEED_RANGE: [0.0005, 0.001],
  ORBIT_RADIUS: [150, 300], // 轨道半径范围
  ANIMATION: {
    HEARTBEAT_INTERVAL: 1000, // 心跳间隔，单位：毫秒
  },
};

// 颜色主题
const THEME = {
  node: {
    // leader是绿色
    leader: {
      gradient: 'from-green-400 to-green-600',
      glow: 'rgba(100, 116, 139, 0.3)',
    },
    // candidate 是蓝色
    candidate: {
      gradient: 'from-amber-400 to-amber-600',
      glow: 'rgba(245, 158, 11, 0.5)',
    },
    // follower 是灰色
    follower: {
      gradient: 'from-slate-600 to-slate-800',
      glow: 'rgba(100, 116, 139, 0.3)',
    },
  },
  // 连接线颜色
  connection: {
    active: '#10B981',
    inactive: '#475569',
  },
};

export default function RaftVisualization() {
  return (
    <RaftLayout>

    </RaftLayout>
  );
}
