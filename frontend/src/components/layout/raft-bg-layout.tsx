// 可以通用在首页和Raft可视化页面的背景布局
import { type ReactNode, useEffect, useRef, useCallback } from 'react';
import { motion } from 'framer-motion';

interface RaftLayoutProps {
  children: ReactNode;
  starDensity?: number;
  withGalaxyEffect?: boolean;
}

// 背景粒子系统参数类型
type StarParticle = {
  x: number;
  y: number;
  size: number;
  speed: number;
  alpha: number;
  direction: number;
};

export function RaftLayout({
  children,
  starDensity = 0.5,
  withGalaxyEffect = true,
}: RaftLayoutProps) {
  // Background canvas for star animation
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationFrameId = useRef<number>(null);
  const particles = useRef<StarParticle[]>([]);

  // 初始化粒子系统
  const initParticles = useCallback(
    (canvas: HTMLCanvasElement) => {
      const baseCount = 200;
      const count = Math.floor(baseCount * starDensity);

      // 初始化粒子位置和参数
      particles.current = Array.from({ length: count }, () => ({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        size: Math.random() * 2 + 0.3, // 调整星星大小
        speed: Math.random() * 0.5 + 0.1, // 降低移动速度
        alpha: Math.random() * 0.4 + 0.2, // 调整透明度
        direction: Math.random() * Math.PI * 2,
      }));
    },
    [starDensity]
  );
  // 动画循环
  const animate = useCallback(() => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (!canvas || !ctx) return;

    ctx.fillStyle = '#0a0e17';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // 绘制银河系核心（条件渲染）
    if (withGalaxyEffect) {
      ctx.beginPath();
      const gradient = ctx.createRadialGradient(
        canvas.width / 2,
        canvas.height / 2,
        0,
        canvas.width / 2,
        canvas.height / 2,
        canvas.width / 2
      );
      gradient.addColorStop(0, 'rgba(16, 185, 129, 0.05)');
      gradient.addColorStop(1, 'rgba(16, 185, 129, 0)');
      ctx.fillStyle = gradient;
      ctx.fillRect(0, 0, canvas.width, canvas.height);
    }
    // 更新并绘制粒子
    particles.current.forEach((particle) => {
      particle.x += Math.cos(particle.direction) * particle.speed;
      particle.y += Math.sin(particle.direction) * particle.speed;

      // 边界重置
      if (particle.x > canvas.width + 50) particle.x = -50;
      if (particle.x < -50) particle.x = canvas.width + 50;
      if (particle.y > canvas.height + 50) particle.y = -50;
      if (particle.y < -50) particle.y = canvas.height + 50;

      // 动态透明度
      particle.alpha = 0.3 + Math.sin(Date.now() * 0.002) * 0.2;

      // 绘制粒子
      ctx.beginPath();
      ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(255, 255, 255, ${particle.alpha})`;
      ctx.fill();
    });

    animationFrameId.current = requestAnimationFrame(animate);
  }, [withGalaxyEffect]);

  // 初始化和动画启动
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    // 初始化画布
    // 响应式
    const handleResize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
      initParticles(canvas);
    };
    // 初始化和动画
    handleResize();
    animationFrameId.current = requestAnimationFrame(animate);
    // 清理
    window.addEventListener('resize', handleResize);
    // 返回一个清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
      if (animationFrameId.current) {
        cancelAnimationFrame(animationFrameId.current);
      }
    };
  }, [animate, initParticles]);

  return (
    <div className="relative min-h-screen w-full overflow-hidden bg-[#0A0B14]">
      <canvas
        ref={canvasRef}
        className="absolute inset-0 h-full w-full"
        style={{ background: 'linear-gradient(to bottom, #0A0B14, #141625)' }}
      />
      {/* 增强型渐变层 */}
      <div className='absolute inset-0 pointer-events-none'>
        <div className='absolute inset-0 bg-gradient-radial from-emerald-500/10 via-transparent to-transparent mix-blend-soft-light' />
        <div className='absolute inset-0 bg-gradient-conic from-blue-500/5 via-transparent to-transparent mix-blend-color-dodge' />
        <div className='absolute bottom-0 h-1/3 w-full bg-gradient-to-t from-slate-900 via-transparent to-transparent' />
      </div>
      {/* 交互式星云效果 */}
      <motion.div
        className='absolute w-[800px] h-[800px] rounded-full blur-3xl bg-emerald-500/10'
        animate={{
          x: [0, 100, -50, 0],
          y: [0, -80, 60, 0],
          scale: [1, 1.2, 0.8, 1],
        }}
        transition={{
          duration: 30,
          repeat: Infinity,
          ease: 'easeInOut',
        }}
      />

      {/* 内容容器 */}
      <div className='relative z-10 w-full h-full backdrop-blur-[1px]'>
        {children}
      </div>
    </div>
  );
}
