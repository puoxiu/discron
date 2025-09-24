package main

import (
	"fmt"
	"os"

	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/server"
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
}