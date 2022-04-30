package taskMgt

import (
	pb "GORushShoping/protobuf"
	"GORushShoping/rabbitmq"
	"GORushShoping/rebuild/order"
	"context"
	"encoding/json"
	"log"
)

var rab *rabbitmq.RabbitMQ

type Task struct {
	//OrderChan chan *order.Order
	OrderChan chan *pb.Order
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewInstance(chanSize int) *Task {
	ctx, cancelFunc := context.WithCancel(context.Background())
	taskObj := &Task{
		OrderChan: make(chan *pb.Order, chanSize),
		ctx:       ctx,
		cancel:    cancelFunc,
	}
	rab = rabbitmq.NewRabbitMQSimple("product")
	go taskObj.run()
	return taskObj
}

func (t *Task) run() {
	go func() { //后台单一线程
		for {
			select {
			case <-t.ctx.Done():
				return
			case orderMsg := <-t.OrderChan: //接收订单
				//go func() { //线程数量取决于channel的数据量

				or := order.Order{ //封装订单信息
					ProductID:  int(orderMsg.ProductID),
					UserID:     int(orderMsg.UserID),
					ProductNum: int(orderMsg.ProductNum),
				}
				//json.Marshal(or) 序列化
				byteOrder, err := json.Marshal(or)
				if err != nil {
					log.Println(err, "json marshal error")
				}
				err = rab.PublishSimple(string(byteOrder)) //发送订单信息,调用RabbitMQ生产者
				if err != nil {
					log.Println(err, "rabbitmq publish error")
				}
				//}()
				//fmt.Println(orderMsg)
			}
		}
	}()
}
