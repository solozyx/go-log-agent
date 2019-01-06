package envinit

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"logagent/logengine/conf"
	"logagent/common/utils"
	"logagent/logengine/tailf"
)

var(
	G_envEngine *EnvEngine
)

type EnvEngine struct {

}

func InitLogEngine()(err error){
	G_envEngine = &EnvEngine{}
	// beego log
	if err = initBeegoLogs(); err != nil{
		return
	}
	// tail
	if err = initTailf(conf.G_engineConf.LogCollectConfList); err != nil{
		return
	}

	logs.Debug("LogEngine env init success :)")
	return
}

func initBeegoLogs() (err error){
	var(
		beegoLogsConfig map[string]interface{}
		confBytes []byte
	)
	beegoLogsConfig = make(map[string]interface{})
	beegoLogsConfig["filename"] = conf.G_engineConf.BeegoLogConf.BeegoLogPath
	beegoLogsConfig["level"] = utils.BeegoConvertLogLevel(conf.G_engineConf.BeegoLogConf.BeegoLogLevel)
	if confBytes,err = json.Marshal(beegoLogsConfig); err != nil {
		// TODO beego 日志文件错误
		logs.Error("LogEngine json Marshal beegoLogsConfig err: %v",err)
		panic(err)
	} else {
		logs.Info("LogEngine json Marshal beegoLogsConfig success ")
	}
	logs.SetLogger(logs.AdapterFile,string(confBytes))
	// 日志默认输出调用的文件名和文件行号,如果你期望不输出调用的文件名和文件行号
	// 开启传入参数 true,关闭传入参数 false,默认是关闭的.
	logs.EnableFuncCallDepth(true)
	return
}

func initTailf(logCollectConfList []*tailf.LogCollectConf)(err error){
	if err = tailf.InitTailfMgr(logCollectConfList); err != nil{
		// TODO tail 库初始化错误
		panic(err)
	}
	return
}