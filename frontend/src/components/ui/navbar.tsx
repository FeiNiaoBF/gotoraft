import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useLanguage } from '@/contexts/language-context';
import { LanguageSwitch } from '@/ui/language-switch';
import { motion } from 'framer-motion';

const navItems = [
  { href: '/process', labelKey: 'processTitle' },
  { href: '/configure', labelKey: 'configureTitle' },
  { href: '/monitor', labelKey: 'monitorTitle' },
  //   { href: '/performance', labelKey: 'performanceTitle' },
  { href: '/help', labelKey: 'helpTitle' },
  { href: '/about', labelKey: 'aboutTitle' },
  { href: '/settings', labelKey: 'settingsTitle' },
];

// 使用motion来做动画
// navVariants 是导航栏的动画
const navVariants = {
  hidden: {
    opacity: 0,
    y: -20,
  },
  visible: {
    opacity: 2,
    y: 0,
    transition: {
      duration: 0.5,
      ease: 'easeInOut',
    },
  },
};

// linkVariants 是用于链接的动画
const linkVariants = {
  hover: {
    scale: 1.1,
    transition: {
      duration: 0.2,
      ease: 'easeInOut',
    },
  },
};

// 登录按钮动画
const loginButtonVariants = {
  hover: {
    scale: 1.05,
    boxShadow: '0 0 20px rgba(59,130,246,0.5)',
  },
  tap: {
    scale: 0.95,
  },
  initial: {
    opacity: 0,
    x: 20,
  },
  animate: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.3,
      ease: 'easeOut',
    },
  },
};

export function Navbar() {
  const pathname = usePathname();
  const { t } = useLanguage();

  const handleLogin = () => {
    // TODO: 实现登录逻辑
    console.log('Login clicked');
  };

  return (
    <div className='w-full'>
      <motion.nav
        className='sticky top-0 z-50 border-b border-slate-800/50 bg-slate-900/50 backdrop-blur-sm'
        initial='hidden'
        animate='visible'
        variants={navVariants}>
        <div className='mx-auto max-w-7xl px-4'>
          <div className='flex h-16 items-center justify-between'>
            <Link
              href='/'
              className='text-2xl md:text-3xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent'>
              GoToRaft
            </Link>
            <div className='hidden md:flex items-center'>
              <div className='flex items-center justify-between w-[800px] mr-8'>
                {navItems.map((item) => (
                  <motion.div
                    key={item.href}
                    className='flex items-center justify-center'
                    variants={linkVariants}
                    whileHover='hover'>
                    <Link
                      href={item.href}
                      className={`text-sm md:text-base font-medium transition-colors hover:text-white px-2 fontsize-sm
                      ${
                        pathname === item.href ? 'text-white' : 'text-slate-400'
                      }`}>
                      {t(item.labelKey)}
                    </Link>
                  </motion.div>
                ))}
              </div>
              {/* 语言切换按钮 */}
              <div className='flex items-center space-x-8'>
                <LanguageSwitch className='z-50' />
                <motion.button
                  onClick={handleLogin}
                  variants={loginButtonVariants}
                  initial='initial'
                  animate='animate'
                  whileHover='hover'
                  whileTap='tap'
                  className='px-6 py-2 rounded-full bg-gradient-to-r from-blue-500 to-blue-600 text-white font-medium text-sm
                    transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900'>
                  {t('login')}
                </motion.button>
              </div>
            </div>
          </div>
        </div>
      </motion.nav>
    </div>
  );
}
