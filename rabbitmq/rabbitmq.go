package rabbitmq

import (
	"GORushShoping/services"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// MQURL 连接信息
const MQURL = "amqp://Bettry:748226@localhost:5672/GORush"

//RabbitMQ 结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind key 名称
	Key string
	//连接信息
	Mqurl string
	sync.Mutex
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	if err := r.channel.Close(); err != nil {
		fmt.Println(err)
	}
	if err := r.conn.Close(); err != nil {
		fmt.Println(err)
	}
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s : %s", message, err)
		//panic(fmt.Sprintf("%s : %s",message,err))
	}
}

// NewRabbitMQSimple 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitMQ := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq")
	//获取channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open channel")
	return rabbitMQ
}

// PublishSimple 简单模式队列生产
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	//1.申请队列,如果队列不存在会自动创建,存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//排他性
		false,
		//阻塞
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//调用channel,发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果true,根据则神exchange类型和routekey规则无法找到符合条件的队列会把消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有消费者,则会把消息返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return nil
}

// ConsumeSimple simple模式下消费者
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	//1.申请队列,如果队列不存在自动创建,存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//持久化
		false,
		//自动删除
		false,
		//排他性
		false,
		//阻塞处理
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//设置流量控制,防止数据库暴库
	r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量(以八位字节为单位)
		false, //如果设置为true,对channel可用
	)

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",    //用来区分多个消费者
		false, //自动应答,这里为false,我们使用手动处理
		false, //如果设置为true,表示不能将同一个connection中生产者发送的消息传递给这个connection中的消费者
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启用线程处理消息
	go func() {
		//消息处理逻辑处理,根据业务灵活调动
		for d := range msgs {
			//log.Printf("Received a message: %s\n", d.Body)
			/*message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}*/
			//插入订单
			/*_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}*/
			//扣除商品数量
			//err = productService.SubNumberOne(message.ProductID)
			err = productService.SubNumberOne(1)
			if err != nil {
				fmt.Println(err)
			}
			d.Ack(false) //如果为true表示确认所有未确认的消息,为false表示确认当前消息
			//log.Printf("Resolved the order : %s\n", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
