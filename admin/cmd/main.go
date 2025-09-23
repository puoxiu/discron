package main

import (
	"fmt"
	"github.com/puoxiu/discron/common/pkg/config"
	"os"
)

func main() {
	err := config.LoadConfig("admin")
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}
}
