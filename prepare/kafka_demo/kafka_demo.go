package main

import (
	"fmt"
	"time"
	"github.com/Shopify/sarama"
)

func main() {
	// 初始化kafka配置
	config := sarama.NewConfig()
	// 使用kafka的生产者 向MQ投放数据 消费者消费MQ数据
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区负载均衡 随机的 比如kafka有8个分区随机打到8个分区上
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 生产者对象
	client, err := sarama.NewSyncProducer([]string{"192.168.234.142:9092"}, config)
	if err != nil {
		fmt.Println("kafka producer new err:", err)
		return
	}
	defer client.Close()

	// 数据
	for{
		msg := &sarama.ProducerMessage{}
		msg.Topic = "nginx_log"
		msg.Value = sarama.StringEncoder("this is a good test, my message is good")
		// 生成者向MQ写入数据
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			fmt.Println("kafka Producer send message failed,", err)
			return
		}
		fmt.Printf("pid:%v offset:%v\n", pid, offset)
		time.Sleep(10 * time.Millisecond)
	}
}