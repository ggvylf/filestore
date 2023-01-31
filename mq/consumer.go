package mq

import (
	"fmt"
)

var done chan bool

// 消费消息
func StartConsume(queneName, consuseName string, callback func(msg []byte) bool) {

	// 从channel中获取消息
	msgs, err := channel.Consume(queneName, consuseName, true, false, false, false, nil)
	if err != nil {
		fmt.Println("read mq channel err,err=", err.Error())
		return
	}

	done = make(chan bool)
	go func() {
		for msg := range msgs {
			suc := callback(msg.Body)
			if !suc {
				// TODO: 重试消费
			}

		}
	}()

	<-done

	// 关闭channel
	channel.Close()

}

func stopConsume() {
	done <- true
}
