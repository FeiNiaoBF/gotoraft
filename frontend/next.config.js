/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*' // 假设后端运行在 8080 端口
      }
    ]
  }
}

module.exports = nextConfig
