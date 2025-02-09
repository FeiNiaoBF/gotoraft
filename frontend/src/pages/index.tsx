import React from 'react';
import { ScrollProgress } from '@/components/layout/scroll-progress';
import { Navbar } from '@/ui/navbar';
import ScrollSection from '@/components/layout/scroll-section';
import { useLanguage } from '@/contexts/language-context';
import dynamic from 'next/dynamic';
import { Suspense } from 'react';
import ErrorBoundary from '@/components/error-boundary';
import LoadingSpinner from '@/components/loading-spinner';

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
      <div className='fixed top-0 right-0 w-1/2 h-screen flex flex-col items-center justify-center'>
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
          <div className='min-h-screen flex items-center'>
            <div className='space-y-10 px-12'>
              <div className='space-y-6'>
                <div className='inline-block'>
                  <p className='text-blue-400 font-semibold tracking-wide px-6 py-3 text-lg rounded-full bg-blue-500/10 border border-blue-500/40'>
                    {t('title')}
                  </p>
                </div>
                <h1 className='text-7xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white via-blue-100 to-blue-200 leading-tight'>
                  {t('subtitle')}
                </h1>
                <p className='text-slate-300 text-xl leading-relaxed max-w-xl'>
                  {t('description')}
                </p>
              </div>

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
          </div>
        </ScrollSection>

        <ScrollSection>
          <div className='min-h-screen flex items-center'>
            <div className='space-y-8 px-12 max-w-xl'>
              <h2 className='text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white to-blue-200'>
                {t('processTitle')}
              </h2>
              <p className='text-xl text-slate-300'>
                {t('leaderElectionDescription')}
              </p>
            </div>
          </div>
        </ScrollSection>

        <ScrollSection>
          <div className='min-h-screen flex items-center'>
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
                <div
                  key={index}
                  className='bg-slate-800/40 backdrop-blur-sm rounded-xl p-8 border border-slate-700/50 transition-all duration-300 hover:border-blue-500/30 hover:shadow-[0_0_30px_rgba(59,130,246,0.2)]'>
                  <h3 className='text-2xl font-semibold text-transparent bg-clip-text bg-gradient-to-r from-white to-blue-200 mb-4'>
                    {feature.title}
                  </h3>
                  <p className='text-slate-300 text-lg leading-relaxed'>
                    {feature.description}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </ScrollSection>
      </div>
    </main>
  );
}

export default function ClientHome() {
  return <Home />;
}
