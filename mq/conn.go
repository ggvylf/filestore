package mq

import (
	"fmt"
	"log"

	"github.com/ggvylf/filestore/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

// UpdateRabbitHost : 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

// Init : 初始化MQ连接信息
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}

// 初始化channel
func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	// 获取连接
	conn, err := amqp.Dial(rabbitHost)
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
