package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	etcd_client "github.com/coreos/etcd/clientv3"
)

const (
	EtcdKey = "/solozyx/backend/logagent/config/127.0.0.1"
)

type LogConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// connect etcd success
// /logagent/conf/ : sample_value
func main() {
	cli,_ := connEtcd()
	etcdExmaple(cli)
	// setLogConfToEtcd(cli)
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

func etcdExmaple(cli *etcd_client.Client) {
	defer cli.Close()

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/logagent/conf/", "sample_value")
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/logagent/conf/")
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

// connect etcd success
// kv.Key = /solozyx/backend/logagent/config/127.0.0.1 : kv.Value =
// [{"path":"D:/project/nginx/logs/access.log","topic":"nginx_log"},
// {"path":"D:/project/nginx/logs/error.log","topic":"nginx_log_err"}]
func setLogConfToEtcd(cli *etcd_client.Client) {
	defer cli.Close()

	var err error

	var logConfArr []LogConf
	logConfArr = append(
		logConfArr,
		LogConf{
			Path:  "D:/project/nginx/logs/access.log",
			Topic: "nginx_log",
		},
	)
	logConfArr = append(
		logConfArr,
		LogConf{
			Path:  "D:/project/nginx/logs/error.log",
			Topic: "nginx_log_err",
		},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	// 访问etcd用context做超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, string(data))
	// etcd put 操作完成取消该 ctx
	cancel()
	if err != nil {
		fmt.Println("etcd put value err = ", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("etcd get key err = ", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("kv.Key = %s : kv.Value = %s\n", kv.Key, kv.Value)
	}
}