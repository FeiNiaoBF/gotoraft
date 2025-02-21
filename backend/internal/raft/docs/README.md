# Raft 开发文档

关于该Raft的核心实现，我是借鉴[hashicorp/raft](https://github.com/hashicorp/raft)的实现。来完成我的Raft核心，个人无法做好Raft的所有功能，但是我想实现Raft的主要功能来展示可视化，比如**领导选举**、**心跳，日志同步**、**安全性**和**持久性**。

## 核心想法

我想通过一个键值存储（Key/Value store）来模拟使用了Raft共识算法的节点（peers）之间的通信，并可视化该服务的集群，以展示Raft的核心功能。

## Key/Value store

### Write

在对于写端有两种情况：

1. 未登录情况下。使用内存来存储写的数据
2. 登录情况下。使用第三方来存储来为写入的数据提供持久化

**Set**，用来向Key/Value store中写入数据

### Read

**Get**，用来从Key/Value store中读取数据

## Raft 核心

[Raft 核心](core.md)

```mermaid
graph TD
    %% 颜色定义
    classDef blue fill:#b3d9ff,stroke:#333;
    classDef green fill:#ccffcc,stroke:#333;
    classDef purple fill:#e6ccff,stroke:#333;
    classDef gold fill:#ffd700,stroke:#333;
    classDef red fill:#ffb3b3,stroke:#333;

    %% 客户端层
    Client("客户端
    (HTTP/gRPC请求)"):::blue
    -->|SET/GET请求| KVService

    %% KV服务层
    subgraph K/V服务节点
        KVService["KV服务代理层
        接收客户端请求
        转发Raft提案"]:::green
    end

    %% Raft共识层
    subgraph Raft集群节点
        %% Leader节点结构
        subgraph Leader节点
            RaftCore_Leader["Raft核心
            角色: Leader
            Term: 3
            日志索引: 3
            提交索引: 2"]
            StateMachine_Leader[/"状态机
            a → 1
            b → 2
            x → 5"/]:::green
            Persistence_Leader["持久化存储
            最新快照
            索引: 3"]:::green
        end

        %% Follower节点结构
        subgraph Follower1节点
            RaftCore_F1["Raft核心
            角色: Follower
            Term: 3
            日志索引: 2"]
            StateMachine_F1[/"状态机
            a → 1
            b → 2"/]:::green
            Persistence_F1["持久化存储
            最新快照
            索引: 1"]:::green
        end

        subgraph Follower2节点
            RaftCore_F2["Raft核心
            角色: Follower
            Term: 3
            日志索引: 2"]
            StateMachine_F2[/"状态机
            a → 1
            b → 2"/]:::green
            Persistence_F2["持久化存储
            最新快照
            索引: 1"]:::green
        end
    end

    %% 数据流关系
    KVService -->|1.提案提交| RaftCore_Leader
    RaftCore_Leader -->|2.日志复制| RaftCore_F1
    RaftCore_Leader -->|2.日志复制| RaftCore_F2
    RaftCore_Leader -->|3.提交日志| StateMachine_Leader
    RaftCore_F1 -->|4.应用日志| StateMachine_F1
    RaftCore_F2 -->|4.应用日志| StateMachine_F2

    %% 持久化关系
    RaftCore_Leader -.->|5.持久化日志| Persistence_Leader
    RaftCore_F1 -.->|5.持久化日志| Persistence_F1
    RaftCore_F2 -.->|5.持久化日志| Persistence_F2

    %% 快照机制
    Persistence_Leader -->|6.生成快照| Snapshot[/"压缩快照
    LastIndex: 3
    大小: 128KB"/]:::red

    Snapshot -->|7.定期备份| disk[("硬盘存储")]:::blue
    Snapshot -.->|8.安装快照| Persistence_F1
    Snapshot -.->|8.安装快照| Persistence_F2

    %% 崩溃恢复
    Persistence_Leader -->|9.启动恢复| CrashRecovery["恢复流程
    1. 加载快照@index3
    2. 重放后续日志"]:::red

    %% 可视化关键点
    class RaftCore_Leader,RaftCore_F1,RaftCore_F2 purple
    class StateMachine_Leader,StateMachine_F1,StateMachine_F2 green
    class Persistence_Leader,Persistence_F1,Persistence_F2 green
    class Snapshot red

```

## 开发流程

```mermaid
gantt
    title Raft可视化开发计划
    section 核心实现
    Raft状态机          :active,  des1, 2024-02-16, 4d
    选举/心跳逻辑       :active,  des2, 2024-02-20, 1d
    section 集成开发
    KV服务绑定          :         des3, 2024-02-21, 2d
    可视化数据管道      :         des4, 2024-02-23, 2d
    section 测试
    选举异常测试        :         des5, 2024-02-25, 2d
    日志恢复测试        :         des6, 2024-02-27, 1d
```

## 使用Raft核心可视化流程

```mermaid
graph TD
    subgraph 前端
        A[控制面板] -->|WebSocket| B(可视化引擎)
        B --> C[节点状态图]
        B --> D[日志流展示]
    end

    subgraph 后端
        E[WebSocket服务] --> F[Raft控制器]
        F --> G[注册中心]
        F --> H[Raft节点1]
        F --> I[Raft节点2]
        H <-->|FooRPC| I
        H --> J[存储切换器]
        I --> J
    end

    J -->|未登录| K[内存存储]
    J -->|已登录| L[持久化存储]
```

这个是对Raft核心结构的描述，下面我们来看看Raft的可视化流程

### 对于websocket来说

```mermaid
graph LR
    A[WebSocket] --> B[ConfigManager]
    B --> C[修改Raft参数]
    B --> D[调整网络模拟]
    C --> E[实时生效]
    D --> E
```

### 可视化接口开发

```mermaid
graph TD
    A[WebSocket事件流] --> B[节点状态]
    A --> C[日志流]
    A --> D[网络拓扑]
    E[REST API] --> F[历史记录]
    E --> G[手动操作]
```

### KV服务集成

```mermaid
graph TD
    A[客户端] --> B[kv服务端] --> C[Raft核心]-->D[状态机]
    D[状态机] --> B[kv服务端]
    B[kv服务端] --> A[客户端]
```
