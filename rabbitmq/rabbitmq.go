package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/streadway/amqp"
	"imooc-product/datamodels"
	"imooc-product/services"
	"log"
)

// uri  amqp://账号：密码@地址：端口/vhost
const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 连接信息
	Mqurl string
}

// 创建结构体实例的函数
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	if err != nil {
		golog.Error(err)
		return nil
	}
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		golog.Error(err)
		return nil
	}
	return rabbitmq
}

// 断开channel和connecion
func (r *RabbitMQ) Destory() {
	r.conn.Close()
	r.channel.Close()
}

// RabbitMQ的简单工作模式
func NewSimpleRabbitMQ(queueName string) *RabbitMQ {
	// 简单模式是：默认交换机，空key
	return NewRabbitMQ(queueName, "", "")
}

// 简单模式下的消息生产代码
func (r *RabbitMQ) PublishSimple(message string) error {
	// 申请队列,该操作是幂等的，无需担心重复申请带来的副作用
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称name
		false,       // 是否持久化 durable
		false,       // 如果未使用，是否自动删除autoDelete
		false,       // 是否具有排他性，其他人不能访问队列exclusive
		false,       // 是否阻塞 noWait
		nil,         //额外属性

	)
	if err != nil {
		golog.Error(err)
		return err
	}
	r.channel.Publish(
		r.Exchange,
		r.QueueName, //routing key,一般是目标队列名
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return nil
}

//简单模式下消费消息代码
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	// 申请队列,该操作是幂等的，无需担心重复申请带来的副作用
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称name
		false,       // 是否持久化 durable
		false,       // 如果未使用，是否自动删除autoDelete
		false,       // 是否具有排他性，其他人不能访问队列exclusive
		false,       // 是否阻塞 noWait
		nil,         //额外属性
	)
	if err != nil {
		return
	}
	_ = r.channel.Qos(
		1,     //消费者1次消费的最大数量
		0,     //服务器传递的最大容量
		false, //是否全局可用
	)
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		false, //关闭自动应答，消费完一个再来第二个。
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	forever := make(chan bool)
	// 多线程消费消息
	go func() {
		for d := range msgs {
			log.Printf("received a message : %s",d.Body)
			message := new(datamodels.Message)
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}

			// 如果为true 表示确认所有未确定的消息
			_ = d.Ack(false)//告知rabbitmq 消息已被消除 将消息从消息队列中删除 防止消息被重复消费
		}
	}()
	golog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
