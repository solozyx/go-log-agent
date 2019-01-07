package serve

import (
	"time"
	"github.com/astaxie/beego/logs"
	"logagent/logengine/tailf"
	"logagent/logengine/kafka"
)

var (
	G_serve *Serve
)

type Serve struct {

}

func InitServe(){
	G_serve = &Serve{}
	go G_serve.Run()
}

func (serve *Serve)Run() {
	var(
		msg *tailf.TextMsg
		err error
	)
	for{
		msg = tailf.G_tailfMgr.GetOneTextMsg()
		if err = kafka.G_kafkaClient.SendToKafka(msg); err != nil{
			logs.Error("LogEngine send to kafka err = %v",err)
			time.Sleep(time.Second)
			continue
		}
	}
}


