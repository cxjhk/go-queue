package rabbitmq

import (
	"github.com/zeromicro/go-zero/core/threading"
	"log"

	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/queue"
)

type (
	ConsumeHandle func(message string) error

	ConsumeHandler interface {
		Consume(message amqp.Delivery) error
	}

	RabbitListener struct {
		conn    *amqp.Connection
		channel *amqp.Channel
		forever chan bool
		handler ConsumeHandler
		queues  RabbitListenerConf
	}
)

func MustNewListener(listenerConf RabbitListenerConf, handler ConsumeHandler) queue.MessageQueue {
	listener := RabbitListener{queues: listenerConf, handler: handler, forever: make(chan bool)}
	conn, err := amqp.Dial(getRabbitURL(listenerConf.RabbitConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %v", err)
	}

	listener.conn = conn
	channel, err := listener.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	listener.channel = channel
	return listener
}

func (q RabbitListener) Start() {
	for _, que := range q.queues.ListenerQueues {
		for i := 0; i < que.Num; i++ {
			threading.GoSafe(func() {
				msg, err := q.channel.Consume(
					que.Name,
					"",
					que.AutoAck,
					que.Exclusive,
					que.NoLocal,
					que.NoWait,
					nil,
				)
				if err != nil {
					log.Fatalf("failed to listener, error: %v", err)
				}
				for d := range msg {
					if err := q.handler.Consume(d); err != nil {
						logx.Errorf("Error on consuming: %s, error: %v", string(d.Body), err)
					}
				}
			})
		}
	}

	<-q.forever
}

func (q RabbitListener) Stop() {
	q.channel.Close()
	q.conn.Close()
	close(q.forever)
}
