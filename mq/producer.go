package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/ggvylf/filestore/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 初始化channel
func initChannle() bool {
	if channel != nil {
		return true
	}

	// 获取连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Println("conn mq failed,err=", err.Error())
		return false
	}
	// defer conn.Close()

	// 获取channel
	channel, err = conn.Channel()
	if err != nil {
		fmt.Println("mq channel failed,err=", err.Error())
		return false
	}
	// defer channel.Close()

	return true
}

// 投递消息到mq
func Publush(exchange, routingKey string, msg []byte) bool {
	if !initChannle() {
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
	fmt.Println("send msg=", msg)
	return true
}
