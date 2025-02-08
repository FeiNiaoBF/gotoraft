'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Globe } from 'lucide-react';
import { useLanguage } from '@/contexts/language-context';

const navItems = [
  { href: '/', label: 'Home' },
  { href: '/process', labelKey: 'processTitle' },
  { href: '/configure', labelKey: 'configureTitle' },
  { href: '/monitor', labelKey: 'monitorTitle' },
  { href: '/performance', labelKey: 'performanceTitle' },
  { href: '/help', labelKey: 'helpTitle' },
  { href: '/about', labelKey: 'aboutTitle' },
  { href: '/settings', labelKey: 'settingsTitle' },
];

export function Navbar() {
  const pathname = usePathname();
  const { language, setLanguage, t } = useLanguage();

  return (
    <nav className='sticky top-0 z-50 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-sm'>
      <div className='container mx-auto px-4'>
        <div className='flex h-16 items-center justify-between'>
          <Link
            href='/'
            className='text-2xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent'>
            Raft
          </Link>
          <div className='hidden md:flex items-center space-x-4'>
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`text-sm font-medium transition-colors hover:text-white
                  ${pathname === item.href ? 'text-white' : 'text-slate-400'}`}>
                {item.labelKey ? t(item.labelKey) : item.label}
              </Link>
            ))}
            <button
              onClick={() => setLanguage(language === 'en' ? 'zh' : 'en')}
              className='ml-4 p-2 rounded-full hover:bg-slate-800/50 transition-colors'
              aria-label={t('switchLanguage')}>
              <Globe className='w-5 h-5 text-slate-400 hover:text-white' />
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}
