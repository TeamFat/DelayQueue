package config

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
)

// Config 配置文件结构体
type Config struct {
	Name string
}

// Init 配置启动函数
func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	c.initLog()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")                  // 设置配置文件格式为YAML
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

func (c *Config) initLog() {
	os.MkdirAll("./log/", 0777)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	filename := "./log/dq.log"
	logs.SetLogger("file", `{"filename":"`+filename+`","maxdays":7}`)
	logs.SetLevel(logs.LevelDebug)
}
