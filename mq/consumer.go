package mq

import (
	"fmt"

	"github.com/ggvylf/filestore/config"
)

var done chan bool

// 消费消息
func StartConsume(queneName, consuseName string, callback func(msg []byte) bool) {

	if !initChannel(config.RabbitURL) {
		return
	}

	// 从channel中获取消息
	msgs, err := channel.Consume(queneName, consuseName, false, false, false, false, nil)
	if err != nil {
		fmt.Println("read mq channel err,err=", err.Error())
		return
	}

	done = make(chan bool)

	go func() {
		for msg := range msgs {
			suc := callback(msg.Body)
			if !suc {
				fmt.Println("消费失败，需要重新消费")
			} else {
				fmt.Println("消费成功")
				msg.Ack(false)
			}

		}
	}()

	// 阻塞goroutine
	<-done

	// 关闭channel
	channel.Close()

}

func stopConsume() {
	done <- true
}
