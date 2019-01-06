package tailf

import (
	"fmt"
	"time"
	"github.com/hpcloud/tail"
	"github.com/astaxie/beego/logs"
)

var(
	G_tailfMgr *TailfMgr
)

// 需要收集日志的配置项
type LogCollectConf struct {
	// 收集日志文件路径
	LogCollectPath string
	// 该日志对应kafka Topic
	Topic string
}

// 日志收集器
// 每个 Tailf 轮询收集logCollectConf.LogCollectPath 下的日志
type Tailf struct {
	aTail *tail.Tail
	logCollectConf *LogCollectConf
}

// 写入kafka的日志消息数据结构
type TextMsg struct {
	Msg string
	Topic string
}

type TailfMgr struct {
	tailfList []*Tailf
	msgChan chan *TextMsg
}

// 待收集日志配置列表 支持多业务系统日志收集
func InitTailfMgr(logCollectConfList []*LogCollectConf,chanSize int)(err error){
	var(
		aLogCollectConf *LogCollectConf
		aTailf *Tailf
		aTail *tail.Tail
	)

	if len(logCollectConfList) == 0 {
		err = fmt.Errorf("LogEngine InitTailfMgr invalid logCollectConfList")
		return
	}

	G_tailfMgr = &TailfMgr{
		msgChan:make(chan *TextMsg,chanSize),
	}

	// 读取所有业务系统日志
	for _,aLogCollectConf = range logCollectConfList{
		aTailf = &Tailf{
			logCollectConf:aLogCollectConf,
		}
		if aTail,err = tail.TailFile(aLogCollectConf.LogCollectPath,tail.Config{
			ReOpen:true,
			Follow:true,
			// Location:&tail.SeekInfo{Offset:0,Whence:2},
			MustExist:true,
			Poll:true,
		}); err != nil {
			err = fmt.Errorf("LogEngine init tail err = %v",err)
			return
		}
		aTailf.aTail = aTail
		G_tailfMgr.tailfList = append(G_tailfMgr.tailfList,aTailf)

		// TODO start goroutine
		// 收集1个业务系统日志 就启动1个goroutine 收集8个业务系统就起8个goroutine
		go G_tailfMgr.readLogsFromTail(aTailf)
	}

	return
}

func (tailfMgr *TailfMgr)readLogsFromTail(aTailf *Tailf){
	var(
		line *tail.Line
		ok bool
		textMsg *TextMsg
	)
	// TODO for true
	// 做一个配置 优雅退出for
	for true{
		if line,ok = <- aTailf.aTail.Lines; !ok {
			// channel被关闭 是比较重要的日志 Warn
			logs.Warn("LogEngine collecting log path chan close,tail file close reopen, file = %s",aTailf.aTail.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// 读取到1行日志
		textMsg = &TextMsg{
			Msg:line.Text,
			Topic:aTailf.logCollectConf.Topic,
		}
		tailfMgr.msgChan <- textMsg
	}
}