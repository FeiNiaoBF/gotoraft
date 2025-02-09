import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useLanguage } from '@/contexts/language-context';
import { LanguageSwitch } from '@/ui/language-switch';

const navItems = [
  //   { href: '/', labelKey: 'Home' },
  { href: '/process', labelKey: 'processTitle' },
  { href: '/configure', labelKey: 'configureTitle' },
  { href: '/monitor', labelKey: 'monitorTitle' },
  //   { href: '/performance', labelKey: 'performanceTitle' },
  { href: '/help', labelKey: 'helpTitle' },
  { href: '/about', labelKey: 'aboutTitle' },
  { href: '/settings', labelKey: 'settingsTitle' },
];

export function Navbar() {
  const pathname = usePathname();
  const { t } = useLanguage();
  return (
    <div className='w-full'>
      <nav className='sticky top-0 z-50 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-sm'>
        <div className='mx-auto max-w-7xl px-4'>
          <div className='flex h-16 items-center justify-between'>
            <Link
              href='/'
              className='text-2xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent'>
              GoToRaft
            </Link>
            <div className='hidden md:flex items-center space-x-4'>
              {navItems.map((item) => (
                <div
                  key={item.href}
                  className='flex items-center gap-4'>
                  <Link
                    href={item.href}
                    className={`text-sm font-medium transition-colors hover:text-white
                    ${
                      pathname === item.href ? 'text-white' : 'text-slate-400'
                    }`}>
                    {t(item.labelKey)}
                    {/* {item.labelKey ? t(item.labelKey) : item.label} */}
                  </Link>
                  <div>
                    {/* 语言切换按钮 */}
                    <LanguageSwitch className='fixed top-8 right-8 z-50' />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </nav>
    </div>
  );
}
