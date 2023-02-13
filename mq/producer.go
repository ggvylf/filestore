package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/ggvylf/filestore/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// 投递消息到mq
func Publush(exchange, routingKey string, msg []byte) bool {
	if !initChannel(config.RabbitURL) {
		return false
	}

	// 定义投递超时上下文 5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 投递消息
	err := channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		fmt.Println("msg publish err,err=", err.Error())
		return false
	}
	fmt.Println("send msg=", string(msg))
	return true
}
