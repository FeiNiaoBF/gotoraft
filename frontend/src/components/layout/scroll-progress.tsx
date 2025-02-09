import { useEffect, useState, memo, useMemo, useRef } from 'react';
import { motion, useScroll, useSpring, useTransform } from 'framer-motion';

interface ScrollProgressProps {
  /** 进度条颜色 (默认: #10b981) */
  color?: string;
  /** 进度条宽度/高度 (默认: 4px) */
  size?: number;
  /** 显示方向 (默认: vertical) */
  orientation?: 'vertical' | 'horizontal';
  /** 触发可见性的滚动阈值 (默认: 100px) */
  visibilityThreshold?: number;
  /** 自定义类名 */
  className?: string;
  /** 自定义样式 */
  style?: React.CSSProperties;
}

export const ScrollProgress: React.FC<ScrollProgressProps> = memo(
  function ScrollProgress({
    color = '#10b981',
    size = 4,
    orientation = 'vertical',
    visibilityThreshold = 100,
    className = '',
    style = {},
  }: ScrollProgressProps) {
    const [isVisible, setIsVisible] = useState(false);
    const [isScrollable, setIsScrollable] = useState(false);
    const { scrollYProgress } = useScroll();
    const { scrollY } = useScroll();
    const containerRef = useRef<HTMLDivElement>(null);

    // 检查页面是否可滚动
    useEffect(() => {
      const checkScrollable = () => {
        const { clientHeight, scrollHeight } = document.documentElement;
        setIsScrollable(scrollHeight > clientHeight);
      };

      checkScrollable();
      window.addEventListener('resize', checkScrollable);
      return () => window.removeEventListener('resize', checkScrollable);
    }, []);

    // 可见性阈值处理
    useEffect(() => {
      return scrollY.onChange((latest) => {
        const shouldShow = latest > visibilityThreshold;
        setIsVisible(shouldShow);
      });
    }, [scrollY, visibilityThreshold]);

    // 动画配置
    const springConfig = {
      stiffness: 300,
      damping: 30,
      mass: 0.5,
      restDelta: 0.001,
    };

    const smoothProgress = useSpring(scrollYProgress, springConfig);
    const scale = useTransform(smoothProgress, [0, 1], [0.98, 1]);
    const opacity = useTransform(
      smoothProgress,
      [0, 0.1, 0.9, 1],
      [0, 1, 1, 0.4]
    );

    // 方向样式
    const orientationStyles = useMemo(() => {
      const base = {
        background: color,
        borderRadius: `${size / 2}px`,
        position: 'fixed',
        boxShadow: `0 0 ${size * 2}px ${color}40`,
      };

      return orientation === 'horizontal'
        ? {
            ...base,
            width: '100%',
            height: size,
            left: 0,
            top: 0,
            transformOrigin: '0% 50%',
          }
        : {
            ...base,
            width: size,
            height: '100vh',
            right: 0,
            top: 0,
            transformOrigin: '50% 0%',
          };
    }, [orientation, size, color]);

    // 动态变换
    const transform = useTransform(smoothProgress, (value) =>
      orientation === 'horizontal' ? `scaleX(${value})` : `scaleY(${value})`
    );

    if (!isScrollable) return null;

    return (
      <motion.div
        ref={containerRef}
        className={`scroll-progress ${className}`}
        style={{
          ...orientationStyles,
          ...style,
          opacity,
          scale,
          transform,
        }}
        initial={{ opacity: 0 }}
        animate={{
          opacity: isVisible ? 1 : 0,
          transition: { duration: 0.3 },
        }}
        aria-hidden='false'
        role='progressbar'
        aria-valuenow={Math.round(smoothProgress.get() * 100)}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label='页面滚动进度'
      />
    );
  }
);
