# GoToRaft - Raft算法可视化平台

基于Golang和React实现的Raft分布式算法可视化系统，帮助开发者直观地理解Raft协议的运行机制。

## 功能特性

- 实时可视化Raft节点状态和通信过程
- 动态展示选举过程和日志复制
- 支持节点数量、网络延迟等参数配置
- 详细的节点信息和日志查看
- 用户登录和权限管理
- 系统运行日志记录

## 技术栈

### 后端

- 语言：Go 1.22+
- Web框架：Gin
- WebSocket：gorilla/websocket
- 数据库：InfluxDB（用于时序数据存储）(暂定)
- 认证：Github OAuth  （暂定）

### 前端

- node: 20+
- 框架：Next.js 15
- UI库：React 18
- 状态管理：Redux Toolkit
- 可视化：D3.js
- 动画：Framer Motion
- UI组件：Ant Design
- WebSocket：Socket.io-client

## 开发路线图

### 第一阶段：基础架构搭建

1. 后端Raft核心实现
   - Raft节点状态管理
   - 选举机制实现
   - 日志复制功能
   - WebSocket服务

2. 前端项目初始化
   - Next.js项目搭建
   - 基础页面布局
     - 首页
     - Raft 算法流程图
     - 节点设置页面(不确定可以实现吗)
     - 日志和状态监控页面（打算用key/value做一个模拟的Raft使用场景来看日志）
     - 性能分析页面(不确定可以实现吗)
     - 帮助与文档页面(Raft的简介和算法的核心)
     - 关于页面（我打算做一个可以发布，正常工作的项目）
     - 设置页面（用一些按钮实现， 主题切换、源码位置、语言切换等）
   - WebSocket连接管理
   - 状态管理配置

### 第二阶段：可视化实现

1. 节点状态可视化
   - 节点关系图绘制
   - 状态切换动画
   - 消息流动效果

2. 交互功能开发
   - 节点配置界面
   - 参数调节控制
   - 状态查看面板

### 第三阶段：功能完善

1. 用户系统
   - 登录注册
   - 权限管理
   - 个人设置

2. 监控与日志
   - 系统运行日志
   - 性能监控
   - 错误追踪

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- InfluxDB

### 后端启动

```bash
cd backend
go mod tidy
go run main.go
```

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

## 项目结构

```shell
gotoraft/
├── backend/                # Go后端代码
│   ├── api/               # API接口
│   ├── core/              # Raft核心实现
│   ├── models/            # 数据模型
│   └── websocket/         # WebSocket服务
├── frontend/              # Next.js前端代码
│   ├── components/        # React组件
│   ├── pages/            # 页面文件
│   └── store/            # Redux状态管理
└── docs/                 # 项目文档
