package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"logagent/logengine/conf"
	"logagent/logengine/envinit"
)

func main(){
	var(
		err error
	)
    // 初始化日志
	if err = conf.InitLogEngineConfig("ini","./conf/logagent.conf"); err != nil {
		logs.Error("LogEngine read config err: %v",err)
		return
	}
	// 初始化服务组件 beego
	if err = envinit.InitLogEngine(); err != nil {
		logs.Error("LogEngine env init err: %v",err)
		return
	}
	// run
	beego.Run()
}