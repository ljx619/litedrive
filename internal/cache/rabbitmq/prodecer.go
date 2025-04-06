package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

func InitChannel() bool {
	// 判断 channel 是否已经创建
	if rabbitConn != nil {
		return true
	}

	var err error // 提前声明 err 变量

	// 获得 RabbitMQ 的连接
	rabbitConn, err = amqp.Dial(RabbitURL) // 这里使用 "="
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return false
	}

	// 打开一个 channel 用于消息的发布与接收
	rabbitChannel, err = rabbitConn.Channel() // 这里也使用 "="
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return false
	}

	return true
}

// 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	// 判断channel是否正常
	if !InitChannel() {
		return false
	}
	//执行消息发布动作
	err := rabbitChannel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
		return false
	}
	return true
}
