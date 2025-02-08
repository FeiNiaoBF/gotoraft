import React, { useEffect, useState } from 'react';
import Head from 'next/head';

interface SystemOverview {
  system: {
    status: string;
    startTime: string;
    version: string;
  };
  stats: {
    totalRequests: number;
    uptime: string;
    memoryUsage: string;
  };
  features: string[];
}

const Home: React.FC = () => {
  const [overview, setOverview] = useState<SystemOverview | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOverview = async () => {
      try {
        const response = await fetch('/api/system/info');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setOverview(data);
      } catch (error) {
        console.error('获取系统概览失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchOverview();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-gray-600">加载中...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Head>
        <title>系统概览 - GoToRaft</title>
        <meta name="description" content="系统概览页面" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="container mx-auto px-4 py-8">
        <h1 className="text-4xl font-bold text-center text-gray-900 mb-8">
          系统概览
        </h1>

        {overview && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="card">
              <h2 className="card-title">系统信息</h2>
              <div className="space-y-3">
                <p className="stat-label">系统状态：
                  <span className="stat-value text-green-600 ml-2">{overview.system.status}</span>
                </p>
                <p className="stat-label">启动时间：
                  <span className="stat-value ml-2">{overview.system.startTime}</span>
                </p>
                <p className="stat-label">版本：
                  <span className="stat-value ml-2">{overview.system.version}</span>
                </p>
              </div>
            </div>

            <div className="card">
              <h2 className="card-title">运行统计</h2>
              <div className="space-y-3">
                <p className="stat-label">总请求数：
                  <span className="stat-value ml-2">{overview.stats.totalRequests}</span>
                </p>
                <p className="stat-label">运行时长：
                  <span className="stat-value ml-2">{overview.stats.uptime}</span>
                </p>
                <p className="stat-label">内存使用：
                  <span className="stat-value ml-2">{overview.stats.memoryUsage}</span>
                </p>
              </div>
            </div>
          </div>
        )}

        {overview && (
          <div className="card mt-8">
            <h2 className="card-title">系统特性</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {overview.features.map((feature, index) => (
                <div key={index} className="bg-blue-50 rounded-lg p-4 text-center">
                  <span className="text-blue-600 font-medium">{feature}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default Home;
