package main

import (
	"context"
	"fmt"
	"time"
	etcd_client "github.com/coreos/etcd/clientv3"
	"strconv"
)

func main() {
	cli, _ := connEtcd()
	defer cli.Close()

	go func() {
		time.Sleep(2*time.Second)
		for i := 0; i < 5; i++ {
			// 没有超时
			cli.Put(context.Background(), "/logagent/conf/", strconv.Itoa(i))
		}
	}()

	for {
		rch := cli.Watch(context.Background(), "/logagent/conf/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("etcd watch %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

func connEtcd()(client *etcd_client.Client,err error){
	client, err = etcd_client.New(etcd_client.Config{
		// 实际应用是 3个实例 5个实例 根据业务实际压力情况 自动做健康检查容错处理
		Endpoints:   []string{"127.0.0.1:2379"},
		// 连接超时
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect etcd err = ", err)
		return
	}
	fmt.Println("connect etcd success")
	return
}