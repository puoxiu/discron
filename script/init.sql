-- 初始化所有表（完全匹配 models 包结构，表名与 CronixXXXTableName 常量对应）
USE cronix; -- 切换到业务数据库（需与 docker-compose 中 MYSQL_DATABASE 一致）

-- 1. 用户表（对应 models.User，表名常量 CronixUserTableName = "cronix_user"）
CREATE TABLE IF NOT EXISTS `cronix_user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '用户ID（主键）',
  `username` VARCHAR(50) NOT NULL COMMENT '用户名',
  `password` VARCHAR(100) NOT NULL COMMENT '密码',
  `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
  `role` INT DEFAULT NULL COMMENT '角色（1=普通用户，2=管理员，对应 models.RoleNormal/RoleAdmin）',
  `status` INT DEFAULT NULL COMMENT '用户状态',
  `created` BIGINT NOT NULL COMMENT '创建时间戳（秒）',
  `updated` BIGINT NOT NULL COMMENT '更新时间戳（秒）',
  PRIMARY KEY (`id`) -- 匹配模型 gorm:"column:id;primaryKey"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 2. 分组表（对应 models.Group，表名常量 CronixGroupTableName = "cronix_group"）
CREATE TABLE IF NOT EXISTS `cronix_group` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '分组ID（主键）',
  `name` VARCHAR(50) NOT NULL COMMENT '分组名称（必填，对应 binding:"required"）',
  `created` BIGINT NOT NULL COMMENT '创建时间戳（秒）',
  `updated` BIGINT NOT NULL COMMENT '更新时间戳（秒）',
  PRIMARY KEY (`id`) -- 匹配模型 gorm:"id;primaryKey"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分组表';

-- 3. 节点-分组关联表（对应 models.NodeGroup，表名常量 CronixNodeGroupTableName = "cronix_node_group"）
CREATE TABLE IF NOT EXISTS `cronix_node_group` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '关联ID（主键）',
  `node_uuid` VARCHAR(64) NOT NULL COMMENT '节点UUID（必填，对应 binding:"required"）',
  `group_id` INT NOT NULL COMMENT '分组ID（必填，对应 binding:"required"）',
  PRIMARY KEY (`id`) -- 匹配模型 gorm:"id;primaryKey"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点-分组关联表';

-- 4. 用户-分组关联表（对应 models.UserGroup，表名常量 CronixUserGroupTableName = "cronix_user_group"）
CREATE TABLE IF NOT EXISTS `cronix_user_group` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '关联ID（主键）',
  `user_id` INT NOT NULL COMMENT '用户ID（必填，对应 binding:"required"）',
  `group_id` INT NOT NULL COMMENT '分组ID（必填，对应 binding:"required"）',
  PRIMARY KEY (`id`) -- 匹配模型 gorm:"id;primaryKey"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户-分组关联表';

-- 5. 任务表（对应 models.Job，表名常量 CronixJobTableName = "cronix_job"）
CREATE TABLE IF NOT EXISTS `cronix_job` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '任务ID（主键）',
  `name` VARCHAR(100) NOT NULL COMMENT '任务名称（必填，对应 binding:"required"）',
  `command` TEXT NOT NULL COMMENT '任务命令/HTTP地址（必填，对应 binding:"required"）',
  `cmd_user` VARCHAR(50) DEFAULT NULL COMMENT '执行命令的用户（对应 models.Job.CmdUser）',
  `timeout` BIGINT DEFAULT NULL COMMENT '任务超时时间（秒，对应 models.Job.Timeout）',
  `retry_times` INT DEFAULT NULL COMMENT '失败重试次数（对应 models.Job.RetryTimes）',
  `retry_interval` BIGINT DEFAULT NULL COMMENT '重试间隔（秒，对应 models.Job.RetryInterval）',
  `kind` INT DEFAULT NULL COMMENT '任务类型（1=单机任务，2=分组任务，对应 models.JobKindAlone/JobKindGroup）',
  `type` INT NOT NULL COMMENT '执行类型（1=CMD，2=HTTP，必填，对应 binding:"required"，models.JobTypeCmd/JobTypeHttp）',
  `http_method` INT DEFAULT NULL COMMENT 'HTTP方法（1=GET，2=POST，对应 models.HTTPMethodGet/HTTPMethodPost）',
  `notify_status` TINYINT(1) DEFAULT NULL COMMENT '是否发送失败通知（对应 models.Job.NotifyStatus）',
  `notify_type` INT DEFAULT NULL COMMENT '通知类型（对应 models.Job.NotifyType）',
  `status` INT DEFAULT NULL COMMENT '任务状态（对应 models.Job.Status）',
  `notify_to` MEDIUMTEXT DEFAULT NULL COMMENT '通知接收人（二进制存储，对应 models.Job.NotifyTo）',
  `notify_to_type` INT DEFAULT NULL COMMENT '通知接收人类型（对应 models.Job.NotifyToType）',
  `spec` VARCHAR(50) NOT NULL COMMENT 'Cron表达式（对应 models.Job.Spec）',
  `created` BIGINT NOT NULL COMMENT '创建时间戳（秒）',
  `updated` BIGINT NOT NULL COMMENT '更新时间戳（秒）',
  `run_on` VARCHAR(64) DEFAULT NULL COMMENT '执行节点UUID（对应 models.Job.RunOn）',
  PRIMARY KEY (`id`) -- 匹配模型 gorm:"column:id;primaryKey"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- 6. 任务日志表（对应 models.JobLog，表名常量 CronixJobLogTableName = "cronix_job_log"）
CREATE TABLE IF NOT EXISTS `cronix_job_log` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `name` VARCHAR(100) NOT NULL COMMENT '任务名称（对应 models.JobLog.Name）',
  `group_id` INT NOT NULL COMMENT '分组ID（对应 models.JobLog.GroupId）',
  `job_id` INT NOT NULL COMMENT '任务ID（对应 models.JobLog.JobId）',
  `command` TEXT NOT NULL COMMENT '执行的命令（对应 models.JobLog.Command）',
  `ip` VARCHAR(50) DEFAULT NULL COMMENT '节点IP（对应 models.JobLog.IP）',
  `hostname` VARCHAR(100) DEFAULT NULL COMMENT '节点主机名（对应 models.JobLog.Hostname）',
  `node_uuid` VARCHAR(64) NOT NULL COMMENT '节点UUID（对应 models.JobLog.NodeUUID）',
  `success` TINYINT(1) NOT NULL COMMENT '是否成功（0=失败，1=成功，对应 models.JobLog.Success）',
  `output` MEDIUMTEXT DEFAULT NULL COMMENT '执行输出日志（对应 models.JobLog.Output）',
  `spec` VARCHAR(50) NOT NULL COMMENT 'Cron表达式（对应 models.JobLog.Spec）',
  `retry_times` INT DEFAULT NULL COMMENT '重试次数（对应 models.JobLog.RetryTimes）',
  `start_time` BIGINT NOT NULL COMMENT '开始时间戳（秒，对应 models.JobLog.StartTime）',
  `end_time` BIGINT NOT NULL COMMENT '结束时间戳（秒，对应 models.JobLog.EndTime）',
  PRIMARY KEY (`id`) -- 模型未显式定义主键，默认以 id 为主键
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务日志表';

-- 7. 节点表（对应 models.Node，表名常量 CronixNodeTableName = "cronix_node"）
CREATE TABLE IF NOT EXISTS `cronix_node` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '节点ID',
  `pid` VARCHAR(50) DEFAULT NULL COMMENT '节点进程PID（对应 models.Node.PID）',
  `ip` VARCHAR(50) NOT NULL COMMENT '节点IP（对应 models.Node.IP）',
  `hostname` VARCHAR(100) NOT NULL COMMENT '节点主机名（对应 models.Node.Hostname）',
  `uuid` VARCHAR(64) NOT NULL COMMENT '节点UUID（全局唯一，对应 models.Node.UUID）',
  `version` VARCHAR(20) DEFAULT NULL COMMENT '节点版本（对应 models.Node.Version）',
  `up` BIGINT NOT NULL COMMENT '节点启动时间戳（秒，对应 models.Node.UpTime）',
  `down` BIGINT DEFAULT NULL COMMENT '节点下线时间戳（秒，对应 models.Node.DownTime）',
  `status` INT DEFAULT NULL COMMENT '节点状态（1=在线，2=离线，对应 models.NodeConnSuccess/NodeConnFail）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_node_uuid` (`uuid`) -- 节点UUID唯一，匹配 models.Node.FindByUUID() 查询场景
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点表';