# admin 帮助文档

## run
```bash
# normal
go run admin/cmd/main.go

# with pprof
go run admin/cmd/main.go -p

# with health check
go run admin/cmd/main.go -h

go run admin/cmd/main.go -h -p

```

 ```mermaid
graph TD
    %% 启动初始化阶段
    A([main函数启动]) --> B[server.NewApiServer]
    B --> C1[解析命令行参数]
    C1 --> C11{参数校验}
    C11 -->|--v版本| C11a[打印版本→os.Exit]
    C11 -->|--help| C11b[打印帮助→os.Exit]
    C11 -->|合法参数| C12[启动pprof协程]
    C12 --> C13[启动健康检查协程]
    C13 --> C2[环境校验+加载配置]
    C2 --> C3[初始化依赖]
    C3 --> C31[logger.Init]
    C31 --> C32[notify.Init]
    C32 --> C33[MySQL建库+dbclient初始化]
    C33 --> C34[etcdclient.Init]
    C34 --> C4[创建ApiServer+启动信号监听]
    C4 --> C41[配置Gin模式]

    %% 业务初始化阶段
    C41 --> D[业务初始化]
    D --> D1[注册HTTP路由]
    D1 --> D2[初始化NodeWatcher→加载Etcd节点+启动监听协程]
    D2 --> D3[service.RegisterTables（表迁移+初始数据）]
    D3 --> D4[go notify.Serve]
    D4 --> D5{日志清理周期>0?}
    D5 -->|是| D51[启动日志清理协程→获取closeChan]
    D5 -->|否| E[srv.ListenAndServe]
    D51 --> E[srv.ListenAndServe]

    %% 服务运行阶段
    E --> E1[初始化Gin引擎+Panic恢复中间件]
    E1 --> E2[启动HTTP服务→处理外部请求]
    E2 --> E3[NodeWatcher持续运行]
    E3 --> E31[PUT事件→新增节点+分配任务]
    E3 --> E32[DELETE事件→故障转移+报警邮件]
    E2 --> E4[日志清理协程定时执行]
    E2 --> E5[通知消费协程处理报警]

    %% 优雅关闭阶段
    F[收到关闭信号] --> G[ApiServer.Shutdown]
    G --> G1[执行关闭钩子（预留）]
    G1 --> G2[等待1秒释放轻量资源]
    G2 --> G3[close→终止日志清理]
    G3 --> G4[关闭HTTP服务（15秒超时）]
    G4 --> G5[logger.Shutdown→关闭日志句柄]
    G5 --> H[os.Exit→服务退出]
```



## 一、admin 整体流程概览
admin 是 discron 调度系统的管理端，负责节点管理、任务调度、用户鉴权与报警通知，运行流程分 4 个核心阶段：
* 启动初始化：解析参数→启动辅助服务→加载配置→初始化依赖（日志 / MySQL/Etcd）；
* 业务初始化：注册路由→启动节点监听→初始化数据库→启动通知 / 日志清理协程；
* 服务运行：启动 HTTP 服务→处理请求→维护节点与任务→定时清理日志；
* 优雅关闭：接收信号→释放资源→关闭服务→退出进程。

## 二、admin 详细流程
1. 启动初始化（server.NewApiServer）
完成服务启动前的基础准备，确保依赖可用。

1.1 命令行参数解析
核心参数：--env（环境）、--enable-pprof（性能分析）、--enable-health-check（健康检查）、--config（配置文件）；
特殊处理：-v 打印版本、--help 打印帮助，触发后直接退出；
辅助服务：pprof（默认 8188 端口）、健康检查（默认 8186 端口），启用后通过独立协程运行。

1.2 配置与依赖初始化
环境校验：校验 --env 合法性（如 production/testing），加载对应配置（系统 / 日志 / MySQL/Etcd/Email）；
依赖初始化顺序：
日志：按配置创建实例，支持文件 / 控制台输出；
notify：初始化邮件 / WebHook 通知队列（缓冲 64 条）；
MySQL：先建库，再初始化连接池；
Etcd：创建客户端，用于节点监听与任务存储。

1.3 ApiServer 准备
创建实例：绑定 HTTP 监听地址（从配置读取）；
信号监听：监听 SIGINT/SIGHUP/SIGTERM，触发时启动优雅关闭；
Gin 模式：生产环境用 ReleaseMode（禁日志），测试环境用 DebugMode（打请求日志）。

2. 业务初始化（main 函数）
完成业务层核心组件初始化。

2.1 注册 HTTP 路由
全局中间件：Cors 处理跨域；
路由分组（按鉴权拆分）：
路由前缀	鉴权方式	功能	关键接口示例
/ping	无	基础探测	GET→pong
根路径	无	公开接口（注册 / 登录）	/register//login
/statis	JWT	统计查询	/today//system
/job	JWT	任务管理	/add//kill
/user	JWT	用户管理	/update//change_pw
/node	JWT	节点管理	/search//del

2.2 NodeWatcher 启动（节点核心管理）
初始化：关联 Etcd 客户端，加载现有节点到内存列表；
监听逻辑：
新增节点（PUT）：解析节点信息→延迟 5 秒（等节点就绪）→分配未分配任务；
节点下线（DELETE）：移除内存节点→迁移节点任务（故障转移）→发送失活报警→成功则删数据库记录。

2.3 数据库与后台协程
表初始化：用 GORM 迁移 User/Node/Job/JobLog 表，注入默认管理员（root/123456）；
通知消费：启动协程处理通知队列，邮件同步发、WebHook 异步发；
日志清理：配置周期 > 0 时启动定时协程，删除过期任务日志。

3. 服务运行（srv.ListenAndServe）
启动服务并处理请求与后台任务。
3.1 Gin 引擎与 HTTP 服务
中间件：apiRecoveryMiddleware 捕获 Panic，打印脱敏日志（Authorization 隐去），返回 500；
HTTP 服务：绑定 Gin 引擎，设 20 秒读写超时，启动后处理外部请求。
3.2 后台任务运行
NodeWatcher：持续维护节点与任务；
日志清理：按周期执行；
通知消费：实时处理报警。

4. 优雅关闭（收到关闭信号）
流程：执行关闭钩子→等 1 秒释放资源→关闭日志清理协程→15 秒内关闭 HTTP 服务→关闭日志句柄→退出。