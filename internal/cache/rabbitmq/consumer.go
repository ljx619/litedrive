package rabbitmq

import "log"

var done = make(chan bool)

// 开始监听队列,获取消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	//通过 channel.consume 获得消息信道
	msgs, err := rabbitChannel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		//循环获取队列中的新消息
		for msg := range msgs {
			//调用callback方法来处理新的消息
			procssSuc := callback(msg.Body)
			if !procssSuc {
				//TODO 将任务写到另一个队列中，用于一场情况的重试
			}
		}
	}()

	// done 没有信息 就会进入阻塞状态
	<-done

	//如果收到信息,则会执行
	rabbitChannel.Close()
}
