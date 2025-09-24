package main

import (
	"os"

	// "github.com/puoxiu/discron/common/pkg/config"
	"github.com/puoxiu/discron/admin/internal/handler"
	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/server"
)

const  (
	ServerName="admin"
)

func main() {
	// c, err := config.LoadConfig("admin", "testing", "main")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(111)
	// }
	// fmt.Println(c)

	srv, err := server.NewApiServer(ServerName)
	if err!= nil {
		logger.Errorf("new api server error:%s", err.Error())
		os.Exit(111)
	}
	logger.Infof("logger init success")

	srv.RegisterRouters(handler.RegisterRouters)
	

	err = srv.ListenAndServe()
	if err!= nil {
		logger.Errorf("api-server:listen and serve error:%s", err.Error())
		os.Exit(111)
	}

	logger.Infof("api-server:listen and serve success")
}
