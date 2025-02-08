import { Navbar } from '@/components/navbar';
import { ScrollSection } from '@/components/scroll-section';
import RaftVisualization from '@/components/raft-visualization';
import { ScrollProgress } from '@/components/scroll-progress';
import { LanguageSwitch } from '@/components/language-switch';
import { useLanguage } from '@/contexts/language-context';

function Home() {
  const { t } = useLanguage();

  return (
    <main className='bg-[#1a1f2c] relative overflow-hidden'>
      <ScrollProgress />
      <LanguageSwitch />

      {/* Enhanced background gradients */}
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
        <div
          className='absolute top-[30%] left-[45%] w-[500px] h-[500px]
          bg-gradient-radial from-purple-900/20 via-purple-900/10 to-transparent
          rounded-full blur-3xl'
        />
      </div>

      <Navbar />

      {/* Fixed Raft Visualization */}
      <div className='fixed top-0 right-0 w-1/2 h-screen flex flex-col items-center justify-center'>
        <RaftVisualization />
        <div className='mt-4 text-center text-white'>
          <p className='text-sm'>{t('nodeInstruction')}</p>
          <p className='text-xs text-gray-400'>{t('nodeStates')}</p>
        </div>
      </div>

      {/* Scrollable content */}
      <div className='w-1/2 min-h-screen'>
        <ScrollSection index={0}>
          <div className='min-h-screen flex items-center'>
            <div className='space-y-10 px-12'>
              <div className='space-y-6'>
                <div className='inline-block'>
                  <p className='text-emerald-500 font-semibold tracking-wide px-6 py-3 text-lg rounded-full bg-emerald-500/10 border border-emerald-500/20'>
                    {t('title')}
                  </p>
                </div>
                <h1 className='text-7xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-white via-white to-slate-400 leading-tight'>
                  {t('subtitle')}
                </h1>
                <p className='text-slate-300 text-xl leading-relaxed max-w-xl'>
                  {t('description')}
                </p>
              </div>

              <div className='flex items-center space-x-6'>
                <div className='h-[3px] w-20 bg-gradient-to-r from-emerald-500 to-emerald-500/50' />
                <span className='text-emerald-500 font-mono text-2xl'>01</span>
                <div className='h-[3px] flex-grow bg-slate-800' />
                <span className='text-slate-500 font-mono text-2xl'>03</span>
              </div>

              <div className='flex space-x-6'>
                <button className='group relative px-8 py-4 text-lg bg-gradient-to-r from-emerald-500 to-emerald-600 rounded-full overflow-hidden transition-all duration-300 hover:shadow-[0_0_30px_rgba(16,185,129,0.5)]'>
                  <div className='absolute inset-0 bg-gradient-to-r from-emerald-400 to-emerald-500 translate-y-full group-hover:translate-y-0 transition-transform duration-300' />
                  <span className='relative text-white font-semibold'>
                    {t('startExploration')}
                  </span>
                </button>
                <button className='px-8 py-4 text-lg bg-slate-800 text-white rounded-full font-semibold hover:bg-slate-700 transition-colors'>
                  {t('learnMore')}
                </button>
              </div>
            </div>
          </div>
        </ScrollSection>

        <ScrollSection index={1}>
          <div className='min-h-screen flex items-center'>
            <div className='space-y-8 px-12 max-w-xl'>
              <h2 className='text-5xl font-bold text-white'>
                Understanding Raft
              </h2>
              <p className='text-xl text-slate-300'>
                Raft is designed to be easy to understand. It&apos;s equivalent
                to Paxos in fault-tolerance and performance. The difference is
                that it&apos;s decomposed into relatively independent
                subproblems, and it cleanly addresses all major pieces needed
                for practical systems.
              </p>
            </div>
          </div>
        </ScrollSection>

        <ScrollSection index={2}>
          <div className='min-h-screen flex items-center'>
            <div className='space-y-12 px-12'>
              {[
                {
                  title: 'Leader Election',
                  description:
                    'Raft uses a heartbeat mechanism to trigger leader election. When servers start up, they begin as followers. A server remains in follower state as long as it receives valid RPCs from a leader or candidate.',
                },
                {
                  title: 'Log Replication',
                  description:
                    'Once a leader has been elected, it begins servicing client requests. Each client request contains a command to be executed by the replicated state machines.',
                },
                {
                  title: 'Safety',
                  description:
                    'Raft guarantees that each of these properties is true at all times. The safety properties are upheld without regard to timing, network delays, partitions, and other factors.',
                },
              ].map((feature, index) => (
                <div
                  key={index}
                  className='bg-slate-800/50 backdrop-blur-sm rounded-xl p-8 border border-slate-700/50 transition-all duration-300 hover:border-emerald-500/30 hover:shadow-[0_0_30px_rgba(16,185,129,0.2)]'>
                  <h3 className='text-2xl font-semibold text-white mb-4'>
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
