package main

import (
	"fmt"
	"os"
	"time"

	// "github.com/puoxiu/discron/common/pkg/config"
	"github.com/puoxiu/discron/admin/internal/handler"
	"github.com/puoxiu/discron/admin/internal/service"
	"github.com/puoxiu/discron/common/pkg/config"
	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/notify"
	"github.com/puoxiu/discron/common/pkg/server"
)

const (
	ServerName = "admin"
)

func main() {
	srv, err := server.NewApiServer(ServerName)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("new api server error:%s", err.Error()))
		os.Exit(1)
	}
	//注册API路由业务
	srv.RegisterRouters(handler.RegisterRouters)
	service.DefaultNodeWatcher = service.NewNodeWatcherService()
	err = service.DefaultNodeWatcher.Watch()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("resolver  error:%#v", err))
	}
	//初始化邮件配置
	go notify.Serve()
	var closeChan chan struct{}
	period := config.GetConfigModels().System.LogCleanPeriod
	if period > 0 {
		closeChan = service.RunLogCleaner(time.Duration(period)*time.Minute, config.GetConfigModels().System.LogCleanExpiration)
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("startup api server error:%v", err.Error()))
		close(closeChan)
		os.Exit(1)
	}
	os.Exit(0)
}
