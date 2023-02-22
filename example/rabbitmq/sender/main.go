package main

import (
	"encoding/json"
	"log"

	"github.com/zeromicro/go-queue/rabbitmq"
)

func main() {
	conf := rabbitmq.RabbitSenderConf{RabbitConf: rabbitmq.RabbitConf{
		Host:     "127.0.0.1",
		Port:     5672,
		Username: "guest",
		Password: "guest",
	}, ContentType: "application/json"}
	sender := rabbitmq.MustNewSender(conf)
	data := map[string]interface{}{
		"msg": "json test 111",
	}
	msg, _ := json.Marshal(data)
	err := sender.Send("mqtt-message-exchange", "go-queue", msg)
	if err != nil {
		log.Fatal(err)
	}
}
