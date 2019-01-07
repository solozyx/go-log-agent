# CentOS 7.x

# 安装JDK
yum install -y java-1.8.0-openjdk
java -version

# 安装kafka
# http://mirror.bit.edu.cn/apache/kafka/2.1.0/kafka_2.12-2.1.0.tgz
tar -zxf kafka_2.12-2.1.0.tgz

# 这里使用kafka自带的zookeeper，只启动1个zookeeper，做一个单节点部署
# zookeeper默认监听 2181 端口
./bin/zookeeper-server-start.sh config/zookeeper.properties
netstat -anp | grep 2181

# 启动kafka 默认监听 9092端口
./bin/kafka-server-start.sh config/server.properties
netstat -anp | grep 9092

# 测试kafka
./bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
./bin/kafka-topics.sh --list --zookeeper localhost:2181

# 防火墙放开 9092 端口
firewall-cmd --zone=public --add-port=9092/tcp --permanent
systemctl stop firewalld.service
systemctl start firewalld.service