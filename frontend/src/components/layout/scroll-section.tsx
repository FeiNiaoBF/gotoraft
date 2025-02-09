import { useRef, useEffect, useState, useCallback } from 'react';
import { motion, useScroll, useTransform } from 'framer-motion';
import type { ReactNode } from 'react';

// 滑动动画配置
interface ScrollSectionProps {
  children: ReactNode;
  animationConfig?: {
    opacityRange?: [number, number, number];
    scaleRange?: [number, number, number];
    yShift?: number;
  };
}

/* ScrollSection
    @remarks 实现滚动动画
    @param children 需要动画的子元素
    @param animationConfig 动画配置
    @returns 滚动动画组件
*/
export default function ScrollSection({
  children,
  animationConfig = {
    opacityRange: [0.3, 1, 0.3],
    scaleRange: [0.8, 1, 0.8],
    yShift: 0,
  },
}: ScrollSectionProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [isClient, setIsClient] = useState(false);
  
  // elementMetrics是元素的位置信息
  // top是元素的顶部位置
  // clientHeight是元素的高度
  const [elementMetrics, setElementMetrics] = useState({
    top: 0,
    clientHeight: 0,
  });

  // 防抖处理更新逻辑
  const updateMetrics = useCallback(() => {
    if (!ref.current || typeof window === 'undefined') return;

    const newTop = ref.current.offsetTop;
    const newClientHeight = window.innerHeight;

    setElementMetrics((prev) =>
      prev.top === newTop && prev.clientHeight === newClientHeight
        ? prev
        : { top: newTop, clientHeight: newClientHeight }
    );
  }, []);

  // 客户端初始化
  useEffect(() => {
    setIsClient(true);
    updateMetrics();
  }, [updateMetrics]);

  useEffect(() => {
    if (!isClient) return;

    const element = ref.current;
    if (!element) return;

    // 初始化测量
    updateMetrics();

    // 双保险监测机制
    const resizeObserver = new ResizeObserver(updateMetrics);
    resizeObserver.observe(element);

    const resizeHandler = () => {
      requestAnimationFrame(updateMetrics);
    };

    window.addEventListener('resize', resizeHandler);
    return () => {
      resizeObserver.disconnect();
      window.removeEventListener('resize', resizeHandler);
    };
  }, [isClient, updateMetrics]);

  // 滚动相关
  const { scrollY } = useScroll();

  // 动画计算
  const scrollRange = [
    elementMetrics.top - elementMetrics.clientHeight,
    elementMetrics.top + elementMetrics.clientHeight,
  ];

  // progress是滚动进度
  const progress = useTransform(scrollY, scrollRange, [0, 1]);
  // opacity是透明度
  const opacity = useTransform(
    progress,
    [0, 0.5, 1],
    animationConfig.opacityRange || [0.3, 1, 0.3]
  );
  // scale是缩放
  const scale = useTransform(
    progress,
    [0, 0.5, 1],
    animationConfig.scaleRange || [0.8, 1, 0.8]
  );

  const y = useTransform(progress, [0, 1], [animationConfig.yShift || 0, 0]);

  // 服务器端渲染时返回无动画版本
  if (!isClient) {
    return <div ref={ref}>{children}</div>;
  }

  return (
    <motion.div
      ref={ref}
      style={{
        opacity,
        scale,
        y,
      }}
    >
      {children}
    </motion.div>
  );
}
