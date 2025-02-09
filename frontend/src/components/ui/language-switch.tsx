import { useLanguage } from '@/contexts/language-context';
import { motion } from 'framer-motion';

export function LanguageSwitch({ className = '' }: { className?: string }) {
  const { language, setLanguage } = useLanguage();

  return (
    <motion.div
      className={`flex items-center space-x-2 bg-slate-800/60 backdrop-blur-sm rounded-full p-2 ${className}`}
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}>
      <button
        onClick={() => setLanguage('zh')}
        className={`px-3 py-1 rounded-full transition-all duration-200 ${
          language === 'zh'
            ? 'bg-blue-500 text-white shadow-[0_0_15px_rgba(59,130,246,0.5)]'
            : 'text-gray-400 hover:text-white hover:bg-blue-500/20'
        }`}>
        中文
      </button>
      <button
        onClick={() => setLanguage('en')}
        className={`px-3 py-1 rounded-full transition-all duration-200 ${
          language === 'en'
            ? 'bg-blue-500 text-white shadow-[0_0_15px_rgba(59,130,246,0.5)]'
            : 'text-gray-400 hover:text-white hover:bg-blue-500/20'
        }`}>
        EN
      </button>
    </motion.div>
  );
}
