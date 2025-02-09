import React, { createContext, useContext, useState, useCallback } from 'react';

type Language = 'en' | 'zh';

interface Translations {
  [key: string]: {
    en: string;
    zh: string;
  };
}

const translations: Translations = {
  // 第一页的内容
  title: {
    en: 'Raft Interactive Learning',
    zh: 'Raft 交互式学习',
  },
  subtitle: {
    en: 'Raft Visualization',
    zh: 'Raft共识算法的可视化',
  },
  raftVisualizationPlaceholder: {
    en: 'Interactive Raft Simulation',
    zh: '交互式 Raft 模拟',
  },
  switchLanguage: {
    en: 'Switch to 中文',
    zh: 'Switch to English',
  },
  processTitle: {
    en: 'Raft Process',
    zh: 'Raft 进程',
  },
  configureTitle: {
    en: 'Node Configuration',
    zh: '节点配置',
  },
  monitorTitle: {
    en: 'Monitoring',
    zh: '监控服务',
  },
  performanceTitle: {
    en: 'Performance Analysis',
    zh: '性能分析',
  },
  helpTitle: {
    en: 'Help ',
    zh: '帮助',
  },
  aboutTitle: {
    en: 'About Raft',
    zh: '关于 Raft',
  },
  settingsTitle: {
    en: 'Settings',
    zh: '设置',
  },
  description: {
    en: 'Discover the core principles of the Raft distributed consensus algorithm through an interactive guide. Learn how leader election, log replication, and safety mechanisms ensure reliable distributed consensus.',
    zh: '通过Raft简洁的交互式指南，探索 Raft 的核心原理。了解领导者选举、日志复制和安全机制如何确保分布式共识的可靠性。',
  },
  startExploration: {
    en: 'Begin Learning',
    zh: '开始学习',
  },
  learnMore: {
    en: 'Learn More',
    zh: '了解更多',
  },
  // 第二页的内容

  leaderElectionDescription: {
    en: 'Learn how Raft uses a simple heartbeat and election process to select a leader when needed.',
    zh: '了解 Raft 如何通过简单的心跳和选举过程在需要时选出领导者。',
  },
  leaderElection: {
    en: 'Leader Election',
    zh: '领导者选举',
  },
  logReplication: {
    en: 'Log Replication',
    zh: '日志复制',
  },
  logReplicationDescription: {
    en: 'Discover how the leader replicates log entries to ensure consistency across all nodes.',
    zh: '了解领导者如何复制日志条目以确保所有节点间的数据一致性。',
  },
  safety: {
    en: 'Safety Mechanisms',
    zh: '安全机制',
  },
  safetyDescription: {
    en: 'Explore the safety features that maintain reliable consensus even in adverse conditions.',
    zh: '探索在恶劣条件下依然保持共识可靠性的安全特性。',
  },
  nodeInstruction: {
    en: 'Click on nodes to interact with the simulation',
    zh: '点击节点以与模拟互动',
  },
  nodeStates: {
    en: 'Nodes can be Followers, Candidates, or Leaders',
    zh: '节点可能为跟随者、候选人或领导者',
  },
  backToHome: {
    en: 'Back to Home',
    zh: '返回首页',
  },
};

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

// 语言上下文
const LanguageContext = createContext<LanguageContextType | undefined>(
  undefined
);

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [language, setLanguage] = useState<Language>('en');

  const t = useCallback(
    (key: string) => {
      return translations[key]?.[language] || key;
    },
    [language]
  );

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
}

export function useLanguage() {
  const context = useContext(LanguageContext);
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
}
