package conf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	"logagent/common/conf/constant"
	"logagent/logengine/tailf"
)

var(
	G_engineConf *LogEngineConf
)

type LogEngineConf struct {
	// beego框架日志配置
	BeegoLogConf *BeegoLogConf
	LogCollectConfList []*tailf.LogCollectConf
	// chan size
	ChanSizeLogMsg int
}

type BeegoLogConf struct {
	BeegoLogPath string
	BeegoLogLevel string
}



func InitLogEngineConfig(confType,confPath string)(err error) {
	var (
		beegoConf config.Configer
	)
	if beegoConf,err = config.NewConfig(confType,confPath); err != nil{
		err = fmt.Errorf("LogEngine beego config init err = %v",err)
		panic(err)
	}
	G_engineConf = &LogEngineConf{}

	initBeegoLogger(beegoConf)
	if err = initLogCollectConfList(beegoConf); err != nil{
		return
	}
	initChanSize(beegoConf)

	return
}

func initBeegoLogger(beegoConf config.Configer){
	var (
		beegoLogPath string
		beegoLogLevel string
	)
	// beego log conf
	beegoLogPath = beegoConf.String("logs::log_path")
	beegoLogLevel = beegoConf.String("logs::log_level")
	if len(beegoLogPath) == 0 {
		beegoLogPath = constant.BeegoLogDefaultPath
	}
	if len(beegoLogLevel) == 0 {
		beegoLogLevel = constant.BeegoLogDefaultLevel
	}
	logs.Info("LogEngine beego log init path = %s",beegoLogPath)
	logs.Info("LogEngine beego log init level = %s",beegoLogLevel)
	G_engineConf.BeegoLogConf = &BeegoLogConf{
		BeegoLogPath:beegoLogPath,
		BeegoLogLevel:beegoLogLevel,
	}
}

func initLogCollectConfList(beegoConf config.Configer)(err error){
	var(
		logCollectPath string
		topic string
		collectConf *tailf.LogCollectConf
	)
	logCollectPath = beegoConf.String("collect::log_collect_path")
	if len(logCollectPath) == 0 {
		err = fmt.Errorf("LogEngine logCollectPath conf is null")
		return
	}
	topic = beegoConf.String("collect::topic")
	if len(topic) == 0 {
		err = fmt.Errorf("LogEngine topic conf is null")
		return
	}
	collectConf = &tailf.LogCollectConf{
		LogCollectPath:logCollectPath,
		Topic:topic,
	}
	G_engineConf.LogCollectConfList = append(G_engineConf.LogCollectConfList,collectConf)
	return
}

func initChanSize(beegoConf config.Configer){
	var(
		chanSizeLogMsg int
		err error
	)
	if chanSizeLogMsg,err = beegoConf.Int("channel::chan_size_log_msg"); err != nil{
		chanSizeLogMsg = constant.ChanSizeLogMsg
	}
	G_engineConf.ChanSizeLogMsg = chanSizeLogMsg
	return
}