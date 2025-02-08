'use client';

import { useRef, useEffect, useState } from 'react';
import { motion, useScroll, useTransform } from 'framer-motion';
import type React from 'react'; // Added import for React

interface ScrollSectionProps {
  children: React.ReactNode;
  index: number;
}

export function ScrollSection({ children, index }: ScrollSectionProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [elementTop, setElementTop] = useState(0);
  const [clientHeight, setClientHeight] = useState(0);

  const { scrollY } = useScroll();

  useEffect(() => {
    const element = ref.current;
    const onResize = () => {
      if (element) {
        setElementTop(element.offsetTop);
        setClientHeight(window.innerHeight);
      }
    };
    onResize();
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);

  const progress = useTransform(
    scrollY,
    [elementTop - clientHeight, elementTop + clientHeight],
    [0, 1]
  );

  const opacity = useTransform(progress, [0, 0.5, 1], [0.3, 1, 0.3]);
  const scale = useTransform(progress, [0, 0.5, 1], [0.8, 1, 0.8]);

  return (
    <motion.div
      ref={ref}
      style={{
        opacity,
        scale,
      }}
      className='w-full'>
      {children}
    </motion.div>
  );
}
