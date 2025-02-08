'use client';

import { createContext, useContext, useState, type ReactNode } from 'react';

type Language = 'en' | 'zh';

interface Translations {
  [key: string]: {
    en: string;
    zh: string;
  };
}

const translations: Translations = {
  title: {
    en: 'Raft Consensus Algorithm',
    zh: 'Raft 共识算法',
  },
  subtitle: {
    en: 'Distributed Consensus Simplified',
    zh: '分布式共识的简化实现',
  },
  description: {
    en: 'Explore the intricacies of the Raft algorithm, a cornerstone of modern distributed systems. Witness the dynamic interplay of nodes, leader election, and consensus formation through our interactive visualization.',
    zh: '探索 Raft 算法的精妙之处，这是现代分布式系统的基石。通过我们的交互式可视化，见证节点的动态交互、领导者选举和共识形成的过程。',
  },
  startExploration: {
    en: 'Start Exploration',
    zh: '开始探索',
  },
  learnMore: {
    en: 'Learn More',
    zh: '了解更多',
  },
  nodeInstruction: {
    en: 'Click nodes to toggle their states',
    zh: '点击节点切换状态',
  },
  nodeStates: {
    en: 'Follower ↔ Leader',
    zh: '跟随者 ↔ 领导者',
  },
  backToHome: {
    en: 'Back to Home',
    zh: '返回首页',
  },
  switchLanguage: {
    en: 'Switch to 中文',
    zh: 'Switch to English',
  },
  // Page titles
  processTitle: {
    en: 'Raft Algorithm Process',
    zh: 'Raft 算法流程',
  },
  configureTitle: {
    en: 'Node Configuration',
    zh: '节点配置',
  },
  monitorTitle: {
    en: 'Logs and Monitoring',
    zh: '日志和监控',
  },
  performanceTitle: {
    en: 'Performance Analysis',
    zh: '性能分析',
  },
  helpTitle: {
    en: 'Help and Documentation',
    zh: '帮助和文档',
  },
  aboutTitle: {
    en: 'About This Project',
    zh: '关于本项目',
  },
  settingsTitle: {
    en: 'Settings',
    zh: '设置',
  },
};

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const LanguageContext = createContext<LanguageContextType | null>(null);

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>('en');

  const t = (key: string): string => {
    return translations[key]?.[language] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
}

export function useLanguage() {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
}
