# Raft的流程核心

我需要理清Raft的核心流程，以及各个节点之间的通信。下面是一个简单的流程图。

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client
    participant Leader as Leader
    participant Follower as Follower
    participant LogStore as LogStore (Persistent Storage)

    Client->>Leader: Send write request (command)
    Leader->>Leader: Store command in local log
    Leader->>Leader: Broadcast AppendEntries RPC to followers
    Follower->>Follower: Append command to log
    Follower-->>Leader: Reply success
    Leader->>Leader: Commit log once quorum is reached
    Leader->>Leader: Apply log to state machine
    Leader->>Client: Respond to client with success

    Follower->>Follower: Apply committed log to state machine
    Follower-->>Leader: Reply success after applying log
    Leader->>LogStore: Persist log
    Follower->>LogStore: Persist log

    %% Example of Raft Election
    alt If Leader crashes
        Follower->>Follower: Start election (RequestVote RPC)
        Follower-->>Leader: Vote for a new leader
        Leader->>Follower: Announce new leader
        Leader->>Follower: Start heartbeat (AppendEntries RPC)
    end
```

这个是对Raft核心结构的描述
```mermaid
graph TD
    A[Raft核心结构] --> B[选举机制]
    A --> C[日志复制]
    A --> D[状态机应用]
    B --> E[RequestVote RPC]
    C --> F[AppendEntries RPC]
    E --> G[foorpc集成]
    F --> G
```

## 对于Raft层在我的项目中的表示

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant KV as KV服务
    participant Raft as Raft层
    participant SM as 状态机

    Client->>KV: SET key=value
    KV->>Raft: 提交日志（包含SET命令）
    Raft->>Raft: 日志复制（多数节点确认）
    Raft->>SM: 应用已提交的日志
    SM-->>KV: 更新内存数据
    KV-->>Client: 操作成功响应
```

### 写请求流程

1. 客户端请求首先进入KV服务
2. KV服务将操作封装为日志提交给Raft
3. Raft确保日志被多数节点复制
4. Raft层通知状态机应用该日志

所以说Raft 需要一个独立的状态机来应用日志

## 核心设计原则

- **聚焦核心机制**：保留选举、日志复制、状态机应用等核心流程
- **内存化存储**：用内存代替持久化存储，避免复杂IO操作
- **事件驱动架构**：通过通道传递状态变更事件，便于可视化捕获
- **模拟网络层**：用内存消息传递代替真实RPC，实现可视化动画
- **可观测性优先**：暴露关键状态指标，方便前端实时渲染
