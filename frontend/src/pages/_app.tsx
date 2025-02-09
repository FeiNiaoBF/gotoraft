import '@/styles/globals.css';
import type { AppProps } from 'next/app';
import { LanguageProvider } from '@/contexts/language-context';
import { inter, outfit, robotoMono, notoSansSC } from '@/lib/fonts';

export default function App({ Component, pageProps }: AppProps) {
  return (
    <LanguageProvider>
      <div
        className={`${inter.variable} ${outfit.variable} ${robotoMono.variable} ${notoSansSC.variable} font-sans`}>
        <Component {...pageProps} />
      </div>
    </LanguageProvider>
  );
}
