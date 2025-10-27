package main

import (
	"fmt"
	"os"

	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/dbclient"
	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/notify"
	"github.com/puoxiu/discron/common/pkg/server"
	"github.com/puoxiu/discron/common/pkg/utils/event"
	// "github.com/puoxiu/discron/common/pkg/metrics"
	"github.com/puoxiu/discron/node/internal/service"
)

const ServerName = "node"

func main() {
	var config *models.Config
	var err error
	if config, err = server.InitNodeServer(ServerName); err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("node server config:", config)

	// metrics.StartMetricsServer("9100")

	nodeServer, err := service.NewNodeServer()
	if err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	service.RegisterTables(dbclient.GetMysqlDB())
	
	if err = nodeServer.Register(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("register node into etcd error:%s", err.Error()))
		os.Exit(1)
	}
	if err = nodeServer.Run(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node run error:%s", err.Error()))
		os.Exit(1)
	}

	go notify.Serve()

	logger.GetLogger().Info(fmt.Sprintf("crony node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String()))
	// 注册退出事件
	event.OnEvent(event.EXIT, nodeServer.Stop /*,stopwatcher()*/)
	// 监听退出信号
	event.WaitEvent()
	// 处理退出事件
	event.EmitEvent(event.EXIT, nil)
	logger.GetLogger().Info("exit success")
}
