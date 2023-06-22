### 下载镜像

```bash
docker search zookeeper  #查看zookeeper版本
docker search kafka  #查看kafka版本

docker pull wurstmeister/kafka #下载wurstmeister/kafka  也可以根据下面截图下载stars量最多的
docker pull wurstmeister/zookeeper ##下载wurstmeister/zookeeper  
```

![image-20230523160342634](./kafka-img\image-20230523160342634.png)

![image-20230523160409524](./kafka-img\image-20230523160409524.png)

![image-20230523160945937](./kafka-img\image-20230523160945937.png)

#### 查看镜像

```bash
docker images
```

![image-20230523163136228](./kafka-img\image-20230523163136228.png)

### 打开zookeeper和kafka

#### 启动zookeeper

```bash
docker run -d --name zk -p 2181:2181 -t wurstmeister/zookeeper
```

- `-d`：以后台模式（detached mode）运行容器。
- `--name zookeeper`：指定容器的名称为 "zk"。
- `-p 2181:2181`：映射容器内部的 2181 端口到主机的 2181 端口。
- `-t bitnami/zookeeper`：使用 wurstmeister提供的 ZooKeeper 镜像作为容器的基础镜像。

![image-20230523163308588](./kafka-img\image-20230523163308588.png)

#### 启动kafka

```bash
docker run -d --name kafkatest -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=10.10.38.253:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://10.10.38.253:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -t wurstmeister/kafka

```



- `-d`：以后台模式（detached mode）运行容器。
- `--name kafkatest`：指定容器的名称为 "kafkatest"。
- `-p 9092:9092`：映射容器内部的 9092 端口到主机的 9092 端口。
- `-e KAFKA_BROKER_ID=0`：指定 Kafka Broker 的 ID。
- `-e KAFKA_ZOOKEEPER_CONNECT=10.10.38.253:2181`：指定 ZooKeeper 的连接地址和端口号。
- `-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://10.10.38.253:9092`：指定 Kafka 监听器的公共 IP 地址和端口号。
- `-e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092`：指定 Kafka 监听器的本地地址和端口号。
- `-t wurstmeister/kafka`：使用 wurstmeister 提供的 Kafka 镜像作为容器的基础镜像。

#### 查看运行状态

```
docker ps
```

![image-20230523163946695](./kafka-img\image-20230523163946695.png)

#### 测试连接

进入kafka容器

```
docker exec -it kafka bash
```

![image-20230523164139546](./kafka-img\image-20230523164139546.png)

进入kafka/bin目录

```bash
cd opt/kafka_2.13-2.8.1/bin  #具体用ls查看目录名
```

![image-20230523164328102](./kafka-img\image-20230523164328102.png)

#### 创建主题（topic）

```bash
kafka-topics.sh --create --zookeeper 10.10.38.253:2181 --replication-factor 1 --partitions 1 --topic test1
```

- `--create`：表示创建一个新的主题。
- `--zookeeper 10.10.38.253:2181`：指定连接到的 ZooKeeper 的地址。
- `--replication-factor 1`：指定主题的副本因子，即每个分区的备份数量。
- `--partitions 1`：指定主题的分区数。
- `--topic test`：指定要创建的主题名称为 "test1"。

![image-20230523164644878](./kafka-img\image-20230523164644878.png)

1. ##### 在 `test1` 主题中发送一条消息：

   ```bash
   kafka-console-producer.sh --broker-list localhost:9092 --topic test1
   ```

   - `--broker-list localhost:9092`：指定要连接的 Kafka Broker 的地址。
   - `--topic test1`：指定要发送消息的主题名称为 "test1"。

   ![image-20230523164904171](./kafka-img\image-20230523164904171.png)

   - 

2. ##### 从 `test1` 主题中消费一条消息：

   ```bash
   kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test1 --from-beginning
   ```

- `--bootstrap-server localhost:9092`：指定要连接的 Kafka Broker 的地址。
- `--topic test1`：指定要消费消息的主题名称为 "test1"。
- `--from-beginning`：表示从该主题的开头开始消费消息，而不是从最新的消息开始消费。

![image-20230523164959375](./kafka-img\image-20230523164959375.png)

如果能够成功发送和接收消息，则说明 Kafka Broker 已经成功启动并且正常工作。



### golang代码测试

consumer.go

```golang
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	// 定义Kafka消费者配置
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 创建Kafka消费者
	consumer, err := sarama.NewConsumer([]string{"10.10.38.253:9092"}, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer consumer.Close()

	// 订阅主题
	topic := "ljxtopic"
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln(err)
	}

	// 创建消息通道
	messages := make(chan *sarama.ConsumerMessage, 1024)

	// 并发消费分区消息
	for _, partition := range partitionList {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Println(err)
			continue
		}
		defer partitionConsumer.Close()

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				messages <- message
			}
		}(partitionConsumer)
	}

	// 启动信号监听器
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// 处理消费的消息
	for {
		select {
		case message := <-messages:
			fmt.Printf("Message received! Topic: %s, Partition: %d, Offset: %d\n", message.Topic, message.Partition, message.Offset)
			fmt.Println("Message content:", string(message.Value))
		case <-signalChan:
			log.Println("Interrupt signal received, stop consuming messages!")
			return
		}
	}
}

```

producer.go

```golang
package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

func main() {
	// 定义Kafka生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	// 创建Kafka生产者
	producer, err := sarama.NewSyncProducer([]string{"10.10.38.253:9092"}, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer producer.Close()

	// 构造消息
	message := &sarama.ProducerMessage{
		Topic: "ljxtopic",
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	// 发送消息
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
	} else {
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	}
}

```

先启动consumer，再启动producer，然后consumer就会接收到hello kafka

![image-20230523165741346](./kafka-img/image-20230523165741346.png)

### 主题topic操作

1. ##### 创建一个新的主题

   您可以使用以下命令来创建一个名为 "test" 的新主题：

   ```
   kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
   ```

2. ##### 列出所有的主题

   您可以使用以下命令列出当前 Kafka 集群中的所有主题：

   ```
   kafka-topics.sh --list --zookeeper localhost:2181
   ```

3. ##### 查看主题的详细信息

   您可以使用以下命令查看特定主题的详细信息：

   ```
   kafka-topics.sh --describe --zookeeper localhost:2181 --topic test
   ```

4. ##### 修改主题的副本因子和分区数

   您可以使用以下命令修改特定主题的副本因子和分区数：

   ```
   kafka-topics.sh --alter --zookeeper localhost:2181 --topic test --partitions 2
   kafka-topics.sh --alter --zookeeper localhost:2181 --topic test --replication-factor 2
   ```

5. 删除主题

   ##### 您可以使用以下命令删除特定的主题：

   ```
   kafka-topics.sh --delete --zookeeper localhost:2181 --topic test
   ```

请注意，在生产环境中，应该谨慎地使用主题操作，并遵循最佳实践来管理和维护 Kafka 主题，以确保其稳定性和性能。