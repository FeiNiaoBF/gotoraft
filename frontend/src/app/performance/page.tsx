'use client';

import { PageLayout } from '@/components/page-layout';
import { BackToHome } from '@/components/back-to-home';
import { useLanguage } from '@/contexts/language-context';

export default function PerformancePage() {
  const { t } = useLanguage();

  return (
    <PageLayout
      title={t('performanceTitle')}
      description={t(
        'Analyze the performance of your Raft cluster under different conditions.'
      )}>
      <div className='mb-6'>
        <BackToHome />
      </div>
      <div className='bg-slate-800/50 p-6 rounded-lg'>
        <p className='text-slate-300'>
          {t('Performance charts and analysis will be implemented here.')}
        </p>
      </div>
    </PageLayout>
  );
}
