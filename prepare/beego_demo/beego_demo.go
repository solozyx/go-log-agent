package main

import (
	"fmt"
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

func main() {
	conf := make(map[string]interface{})
	// 本项目的日志输入的日志文件路径
	conf["filename"] = "./logs/logagent.log"
	// 日志级别 DEBUG 在开发使用日志量大 正式上线关闭DEBUG
	conf["level"] = logs.LevelDebug

	data, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(data))
	// 日志默认输出调用的文件名和文件行号,如果你期望不输出调用的文件名和文件行号
	// 开启传入参数 true,关闭传入参数 false,默认是关闭的.
	logs.EnableFuncCallDepth(true)

	logs.Debug("this is a test, my name is %s", "stu01")
	logs.Trace("this is a trace, my name is %s", "stu02")
	logs.Warn("this is a warn, my name is %s", "stu03")
}

func InitConfig(){
	// 配置文件格式 配置文件路径
	conf, err := config.NewConfig("ini", "./conf/logagent.conf")
	if err != nil {
		fmt.Println("beego read config file logagent.conf err:", err)
		return
	}
	// 读取配置 section::item
	port, err := conf.Int("server::listen_port")
	if err != nil {
		fmt.Println("beego read config server::listen_port  err:", err)
		return
	}
	fmt.Println("port:", port)

	log_level := conf.String("logs::log_level")
	if len(log_level) == 0 {
		fmt.Println("beego read config logs::log_level err len = 0")
		log_level = "debug"
	}
	fmt.Println("log_level:", log_level)

	log_path := conf.String("logs::log_path")
	fmt.Println("log_path:", log_path)
}