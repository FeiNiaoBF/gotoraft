// 主页面的背景布局
import React from 'react';

export default function PageBgLayout() {
  return (
    //   {/* Enhanced background gradients */}
    <div className='fixed inset-0 overflow-hidden pointer-events-none'>
      <div
        className='absolute top-[-10%] right-[-10%] w-[800px] h-[800px]
          bg-gradient-radial from-blue-500/30 via-blue-500/10 to-transparent
          rounded-full blur-3xl animate-pulse-slow'
      />
      <div
        className='absolute bottom-[-20%] left-[-10%] w-[1000px] h-[1000px]
          bg-gradient-radial from-blue-600/30 via-blue-600/10 to-transparent
          rounded-full blur-3xl'
      />
      <div
        className='absolute top-[30%] left-[45%] w-[500px] h-[500px]
          bg-gradient-radial from-blue-500/20 via-blue-500/10 to-transparent
          rounded-full blur-3xl'
      />
    </div>
  );
}
