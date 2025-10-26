# DistCron--分布式定时任务调度系统

> 一款轻量级分布式定时任务调度系统，支持任务的集中管理、节点分发执行与状态监控，无需依赖第三方平台（如 Discord），适用于中小规模服务的定时任务场景（如数据同步、日志清理、接口巡检等）目前支持：
* cmd 任务
* http 任务

## 🌟 核心功能
* 分布式架构：支持多执行节点（node）部署，任务自动分发，避免单点故障
* 灵活的任务配置：兼容标准 Cron 表达式，支持任务优先级、超时时间、重试策略配置
* 集中式管理：提供管理模块（admin），支持任务创建、编辑、启停、日志查询等
* 任务持久化：通过数据存储模块（需配置数据库 / 本地文件）确保任务重启不丢失


``` mermaid
graph TD
    %% 顶层：用户交互层
    subgraph Web
        Admin["Admin"]
    end

    %% 调度层
    subgraph "Dispatching center"
        API["API"]
        Etcd["Etcd"]
    end

    %% 节点集群
    subgraph "Node cluster"
        Node1["Node1"]
        Node2["Node2"]
        Node3["Node3"]
        Node4["Node4 Break"]
    end

    %% 数据库 & 告警
    SQL["SQL"]
    Alarm["Alarm Notification"]

    %% 关系连线
    Admin --> API --> Etcd

    %% 作业 CRUD 流程
    API -.->|Job CRUD...| SQL

    %% etcd 与节点通信
    Etcd -->|job| Node1
    Etcd -->|job| Node2
    Etcd -->|job| Node3
    Etcd -->|job| Node4

    Node1 -->|node1 state| Etcd
    Node2 -->|node2 state| Etcd
    Node3 -->|node3 state| Etcd
    Node4 -.->|node4 break| Etcd

    %% 节点之间 failover
    Node4 -->|FailOver / Least Priority| Node1
    Node4 --> Node2
    Node4 --> Node3

    %% Job 执行日志与失败
    Node1 -.->|Job Exec Log| SQL
    Node2 -.->|Job Exec Log| SQL
    Node3 -.->|Job Exec Log| SQL
    Node4 -.->|Job Exec Log| SQL

    SQL -->|Job Exec Fail| Alarm
    Node4 -.->|Node Break| Alarm
```


## 运行截图

**节点状态与负载信息**
![Discron Node Overview](./images/discorn_nodes.png)

**Grafana 可视化 Dashboard**
![Grafana Dashboard](./images/grafana.png)
