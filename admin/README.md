# admin 帮助文档

 ```mermaid
 graph TD
    A([main]) --> B[server.NewApiServer]
    B --> C1[解析命令行参数（pprof/健康检查）]
    B --> C2[加载环境与配置（系统/日志/MySQL/Etcd）]
    B --> C3[初始化依赖（日志→MySQL→Etcd）]
    B --> C4[创建ApiServer，配置Gin模式+信号监听]
    C4 --> D[业务注册]
    D --> D1[注册HTTP路由（/ping+/job接口）]
    D --> D2[启动NodeWatcher（加载Etcd节点+持续监听）]
    D --> E[srv.ListenAndServe]
    E --> E1[初始化Gin引擎+Panic恢复中间件]
    E --> E2[启动HTTP服务，处理外部请求]
    E --> E3[NodeWatcher持续更新节点列表]
    F[收到关闭信号（SIGINT/SIGHUP/SIGTERM）] --> G[ApiServer.Shutdown]
    G --> G1[执行关闭钩子]
    G --> G2[等待1秒释放资源]
    G --> G3[关闭HTTP服务]
    G --> G4[关闭日志句柄]
    G4 --> H[服务退出]
```



## admin 整体流程概览
> admin 是基于 Gin 框架构建的 HTTP 微服务，核心承担两大职责：
1. 对外提供业务 API 接口，支撑前端或其他服务的请求；
2. 监听 Etcd 中节点的动态变化（新增 / 删除），维护节点列表的实时性。


## admin 详细流程
1. ApiServer 初始化阶段
服务启动的基础准备，完成核心依赖与配置的初始化：
* 解析命令行参数（如运行环境、pprof / 健康检查开关）；
* 加载配置文件（包含服务端口、日志、MySQL、Etcd 连接信息）；
* 初始化公共组件（日志、MySQL 客户端、Etcd 客户端）；
* 设置信号监听（响应终止信号，为优雅关闭做准备）。

2. NodeWatcher 启动阶段
同步并监听 Etcd 节点数据，保障节点列表实时性：
* 从 Etcd 拉取已存在的节点数据，初始化本地节点列表；
* 启动异步监听，实时捕获 Etcd 节点的新增 / 删除事件，更新本地列表。

3. API 服务运行阶段
启动 HTTP 服务，接收并处理外部请求：
* 初始化 Gin 引擎，添加崩溃恢复中间件（避免服务因 panic 退出）；
* 绑定业务路由（如节点查询、任务管理接口）与全局中间件；
* 启动 HTTP 监听，阻塞等待并处理外部请求。

4. 优雅关闭阶段
* 收到终止信号（如 Ctrl+C、kill 命令）后，安全清理资源：
* 执行预注册的业务关闭逻辑（如关闭数据库连接）；
* 停止接收新请求，等待已接收的请求处理完成；
* 关闭 HTTP 服务与日志组件，确保资源不泄露。

