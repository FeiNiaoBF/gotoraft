'use client';

import type { ReactNode } from 'react';

interface PageLayoutProps {
  children: ReactNode;
  title: string;
  description: string;
}

export function PageLayout({ children, title, description }: PageLayoutProps) {
  return (
    <div className='min-h-screen bg-[#1a1f2c] text-white'>
      {/* Background gradients */}
      <div className='fixed inset-0 overflow-hidden pointer-events-none'>
        <div
          className='absolute top-[-10%] right-[-10%] w-[800px] h-[800px]
          bg-gradient-radial from-emerald-900/40 via-emerald-900/20 to-transparent
          rounded-full blur-3xl animate-pulse-slow'
        />
        <div
          className='absolute bottom-[-20%] left-[-10%] w-[1000px] h-[1000px]
          bg-gradient-radial from-blue-900/30 via-blue-900/20 to-transparent
          rounded-full blur-3xl'
        />
      </div>

      <div className='relative container mx-auto px-4 py-12'>
        <div className='mb-12'>
          <h1 className='text-4xl font-bold mb-4 bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent'>
            {title}
          </h1>
          <p className='text-xl text-slate-300'>{description}</p>
        </div>
        {children}
      </div>
    </div>
  );
}
