package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"logagent/logengine/tailf"
)

var(
	G_kafkaClient *KafkaClient
)

type KafkaClient struct {
	// 接口 不能使用指针 struct实现了接口 直接用接口赋值 接口类型变量 不需要指针标识
	kafkaProducerClient sarama.SyncProducer
}

func InitKafkaClient(addr string)(err error){
	// 初始化kafka配置
	config := sarama.NewConfig()
	// 使用kafka的生产者 向MQ投放数据 消费者消费MQ数据
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区负载均衡 随机的 比如kafka有8个分区随机打到8个分区上
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	// 生产者对象
	client, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("LogEngine new kafka producer client err:", err)
		return
	}
	logs.Debug("LogEngine new kafka producer client success")
	// 不能关闭连接 defer client.Close()
	G_kafkaClient = &KafkaClient{
		kafkaProducerClient:client,
	}
	return
}

func (kc *KafkaClient)SendToKafka(msg *tailf.TextMsg)(err error){
	var(
		kafkaMsg *sarama.ProducerMessage
		partitionId int32
		offset int64
	)
	logs.Debug("LogEngine KafkaProducerClient read msg = %s , topic = %s \n",msg.Msg,msg.Topic)
	kafkaMsg = &sarama.ProducerMessage{}
	kafkaMsg.Topic = msg.Topic
	kafkaMsg.Value = sarama.StringEncoder(msg.Msg)
	// 生成者向MQ写入数据
	partitionId, offset, err = kc.kafkaProducerClient.SendMessage(kafkaMsg)
	if err != nil {
		err = fmt.Errorf("LogEngine kafkaProducerClient send msg to kafka server err = %v,Msg = %s,Topic = %s",
			err,msg.Msg,msg.Topic)
		return
	}
	logs.Debug("LogEngine kafkaProducerClient send msg to kafka server success, partitionId=%d offset=%d Topic=%s\n",
		partitionId, offset,msg.Topic)

	return
}