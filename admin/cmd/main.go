package main

import (
	"fmt"
	"github.com/puoxiu/discron/common/pkg/config"
	"os"
)

func main() {
	c, err := config.LoadConfig("admin", "testing", "main")
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}
	fmt.Println(c)
}
