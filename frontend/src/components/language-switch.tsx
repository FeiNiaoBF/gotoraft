'use client';

import { useLanguage } from '@/contexts/language-context';

export function LanguageSwitch() {
  const { language, setLanguage } = useLanguage();

  return (
    <button
      onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
      className='fixed top-4 right-4 px-4 py-2 rounded-full bg-slate-800/50 text-white hover:bg-slate-700/50 transition-colors'>
      {language === 'en' ? '中文' : 'English'}
    </button>
  );
}
