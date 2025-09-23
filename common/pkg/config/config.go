package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/utils"
	"os"
	"path"
)

const (
	//ExtensionJson json配置后缀
	ExtensionJson = ".json"
	//ExtensionYaml yaml配置后缀
	ExtensionYaml = ".yaml"
	//ExtensionInI ini配置后缀
	ExtensionInI = ".ini"

	NameSpace = "conf"
)

var c models.Config
var (
	//本地Config自动载入顺序
	autoLoadLocalConfigs = []string{
		ExtensionJson,
		ExtensionYaml,
		ExtensionInI,
	}
)

var ConfOptions struct {
	flags.Options
	Environment string `short:"e" long:"env" description:"Use crony-server environment"`
	Version     bool   `short:"v" long:"version"  description:"Show crony-server version"`
	//todo
	EnablePProfile    bool   `short:"p" long:"enable-pprof"  description:"enable pprof"`
	EnableHealthCheck bool   `short:"a" long:"enable-health-check"  description:"enable health check"`
	ConfigFileName    string `short:"c" long:"config" description:"Use coa-server config file name" default:"main"`
	EnableDevMode     bool   `short:"m" long:"enable-dev-mode"  description:"enable dev mode"`
}

func LoadConfig(profile, envType string) error {
	// 1.解析命令行参数
	var parser = flags.NewParser(&ConfOptions, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		return fmt.Errorf("解析命令行参数失败: %w", err)
	}
	fmt.Println("confOptions is :", ConfOptions)

	// 2.拼接完整的环境配置目录（如 admin/conf/testing）
	envDir := path.Join(profile, NameSpace, envType)
	if !utils.Exists(envDir) {
		return fmt.Errorf("环境配置目录不存在: %s", envDir)
	}
	// 3.查找环境目录下的配置文件（尝试json/yaml/ini后缀）
	var confPath string
	configFileName := ConfOptions.ConfigFileName
	if configFileName == "" {
		configFileName = "main"
	}
	for _, registerExt := range autoLoadLocalConfigs {
		confPath = path.Join(envDir, configFileName+registerExt)
		if utils.Exists(confPath) {
			break	// 找到第一个存在的配置文件，直接使用
		}
	}
	if !utils.Exists(confPath) {
		return fmt.Errorf("在 %s 下未找到配置文件（尝试了后缀: %v）", envDir, autoLoadLocalConfigs)
	}
	fmt.Println("confPath is :", confPath)

	// 5. 用viper加载配置文件
	v := viper.New()
	v.SetConfigFile(confPath)
	// ext := utils.Ext(confPath)
	// v.SetConfigType(ext)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&c); err != nil {
			fmt.Println(err)
		}
	})

	// 6. 反序列化配置到结构体-解析配置到全局变量 c
	if err := v.Unmarshal(&c); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("config is :%#v", c)
	return nil
}
