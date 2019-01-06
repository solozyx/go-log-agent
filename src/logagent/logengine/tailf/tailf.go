package tailf

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var(
	G_tailfMgr *TailfMgr
)

type LogCollectConf struct {
	LogCollectPath string
	Topic string
}

type Tailf struct {
	aTail *tail.Tail
	logCollectConf *LogCollectConf
}

type TailfMgr struct {
	tailfList []*Tailf
}

func InitTailfMgr(logCollectConfList []*LogCollectConf)(err error){
	var(
		aLogCollectConf *LogCollectConf
		aTailf *Tailf
		aTail *tail.Tail
	)

	if len(logCollectConfList) == 0 {
		err = fmt.Errorf("LogEngine InitTailf invalid logCollectConfList")
		return
	}

	G_tailfMgr = &TailfMgr{}

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
	}

	return
}