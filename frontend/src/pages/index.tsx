import React from 'react';
import { ScrollProgress } from '@/components/layout/scroll-progress';
import { Navbar } from '@/ui/navbar';
import ScrollSection from '@/components/layout/scroll-section';
import { useLanguage } from '@/contexts/language-context';
import dynamic from 'next/dynamic';
import { Suspense } from 'react';
import ErrorBoundary from '@/components/error-boundary';
import LoadingSpinner from '@/components/loading-spinner';
import { motion } from 'framer-motion';

// 动画变体配置
const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.3,
      delayChildren: 0.2,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
      ease: 'easeOut',
    },
  },
};

const featureCardVariants = {
  hidden: { opacity: 0, scale: 0.9 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: {
      duration: 0.5,
      ease: [0.43, 0.13, 0.23, 0.96],
    },
  },
  hover: {
    scale: 1.05,
    boxShadow: '0 0 30px rgba(59,130,246,0.3)',
    transition: {
      duration: 0.3,
      ease: 'easeInOut',
    },
  },
};

// 动态导入Raft的可视化组件
const RaftVisualization = dynamic(
  () => import('@/components/raft-demo/raft-visualization'),
  {
    loading: () => <LoadingSpinner />,
    ssr: false, // 由于可视化组件可能依赖浏览器API，禁用SSR
  }
);

import PageBgLayout from '@/layout/page-bg-layout';

function Home() {
  const { t } = useLanguage();

  return (
    <main className='bg-[#0f172a] relative min-h-screen'>
      {/* 顶部进度条 */}
      <ScrollProgress
        orientation='horizontal'
        color='#3b82f6'
        size={4}
        className='fixed top-0 left-0 z-10'
      />

      {/* Navbar */}
      <Navbar />

      {/* Enhanced background gradients */}
      <PageBgLayout />

      {/* Raft Visualization Demo*/}
      <div className='fixed top-0 right-0 w-1/2  h-screen flex flex-col items-center justify-center'>
        <ErrorBoundary
          fallback={
            <div className='text-red-500'>
              Something went wrong with the visualization
            </div>
          }>
          <Suspense fallback={<LoadingSpinner />}>
            <RaftVisualization />
          </Suspense>
        </ErrorBoundary>
        <div className='mt-4 text-center text-white'>
          <p className='text-sm'>{t('nodeInstruction')}</p>
          <p className='text-xs text-gray-400'>{t('nodeStates')}</p>
        </div>
      </div>

      {/* Scrollable content */}
      <div className='w-1/2 pb-32'>
        <ScrollSection>
          <motion.div
            variants={containerVariants}
            initial='hidden'
            animate='visible'
            className='min-h-screen flex items-center'>
            <div className='space-y-10 px-12'>
              <motion.div
                variants={itemVariants}
                className='space-y-6'>
                <motion.div
                  className='inline-block'
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}>
                  <motion.p
                    className='text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-teal-400 shadow-lg'
                    initial={{ opacity: 0, x: -50 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.8, ease: 'easeOut' }}>
                    {t('title')}
                  </motion.p>
                </motion.div>
                <motion.h2
                  variants={itemVariants}
                  className={`font-display text-5xl md:text-6xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white to-blue-200 leading-tight antialiased subpixel-antialiased`}>
                  {t('subtitle')}
                </motion.h2>
                <motion.p
                  variants={itemVariants}
                  className={`text-slate-300 text-xl md:text-2xl leading-relaxed max-w-xl antialiased ${
                    t('currentLang') === 'zh' ? 'font-cn' : 'font-sans'
                  }`}>
                  {t('description')}
                </motion.p>
              </motion.div>

              <div className='flex space-x-6'>
                <button className='group relative px-8 py-4 text-lg bg-gradient-to-r from-blue-500 to-blue-600 rounded-full overflow-hidden transition-all duration-300 hover:shadow-[0_0_30px_rgba(59,130,246,0.5)]'>
                  <div className='absolute inset-0 bg-gradient-to-r from-blue-400 to-blue-500 translate-y-full group-hover:translate-y-0 transition-transform duration-300' />
                  <span className='relative text-white font-semibold'>
                    {t('startExploration')}
                  </span>
                </button>
                <button className='px-8 py-4 text-lg bg-slate-800/60 backdrop-blur-sm text-white rounded-full font-semibold hover:bg-slate-700/80 transition-colors border border-slate-700/50'>
                  {t('learnMore')}
                </button>
              </div>
            </div>
          </motion.div>
        </ScrollSection>

        <ScrollSection>
          <motion.div
            variants={containerVariants}
            initial='hidden'
            whileInView='visible'
            viewport={{ once: true, amount: 0.3 }}
            className='min-h-screen flex items-center'>
            <div className='space-y-8 px-12 max-w-xl'>
              <motion.h2
                variants={itemVariants}
                className={`text-5xl md:text-6xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white to-blue-200 antialiased subpixel-antialiased ${
                  t('currentLang') === 'zh'
                    ? 'font-cn tracking-normal'
                    : 'font-display tracking-wide'
                }`}>
                {t('processTitle')}
              </motion.h2>
              <motion.p
                variants={itemVariants}
                className={`text-slate-300 text-xl md:text-2xl leading-relaxed antialiased ${
                  t('currentLang') === 'zh' ? 'font-cn' : 'font-sans'
                }`}>
                {t('leaderElectionDescription')}
              </motion.p>
            </div>
          </motion.div>
        </ScrollSection>

        <ScrollSection>
          <motion.div
            variants={containerVariants}
            initial='hidden'
            whileInView='visible'
            viewport={{ once: true, amount: 0.3 }}
            className='min-h-screen flex items-center'>
            <div className='space-y-12 px-12'>
              {[
                {
                  title: t('leaderElection'),
                  description: t('leaderElectionDescription'),
                },
                {
                  title: t('logReplication'),
                  description: t('logReplicationDescription'),
                },
                {
                  title: t('safety'),
                  description: t('safetyDescription'),
                },
              ].map((feature, index) => (
                <motion.div
                  key={index}
                  variants={featureCardVariants}
                  whileHover='hover'
                  className='bg-slate-800/40 backdrop-blur-sm rounded-xl p-8 border border-slate-700/50'>
                  <motion.h3
                    variants={itemVariants}
                    className={`text-2xl md:text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white to-blue-200 mb-4 antialiased subpixel-antialiased ${
                      t('currentLang') === 'zh'
                        ? 'font-cn tracking-normal'
                        : 'font-display tracking-wide'
                    }`}>
                    {feature.title}
                  </motion.h3>
                  <motion.p
                    variants={itemVariants}
                    className={`text-slate-300 text-lg leading-relaxed antialiased ${
                      t('currentLang') === 'zh' ? 'font-cn' : 'font-sans'
                    }`}>
                    {feature.description}
                  </motion.p>
                </motion.div>
              ))}
            </div>
          </motion.div>
        </ScrollSection>
      </div>
    </main>
  );
}

export default function ClientHome() {
  return <Home />;
}
