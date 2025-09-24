package main

import (
	"fmt"
	"os"

	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/server"
	"github.com/puoxiu/discron/common/pkg/utils/event"
)
const ServerName = "node"

func main() {
	// err := server.InitNodeServer(ServerName)
	nodeServer, err := server.NewNodeServer(ServerName)
	if err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
	logger.Debugf("nodeServer:%#v", *nodeServer)
	logger.Debugf("node:%#v", *nodeServer.Node)

	if err = nodeServer.Register(); err != nil {
		logger.Errorf("nodeServer register error:%v", err.Error())
		os.Exit(1)
	}

	logger.Infof("cronix node %s service started, Ctrl+C or send kill sign to exit", nodeServer.String())
	// 注册退出事件
	event.OnEvent(event.EXIT, nodeServer.Stop /*,stopwatcher()*/)
	// 监听退出信号
	event.WaitEvent()
	// 处理退出事件
	event.EmitEvent(event.EXIT, nil)
	logger.Infof("exit success")
}