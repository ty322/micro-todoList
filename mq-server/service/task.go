package service

import "mq-server/model"

// 从rabbitmq中接受信息, 写入数据库
func CreateTask() {
	ch, err := model.MQ.Channel()
	if err != nil {
		panic(err)
	}
	q, _ := ch.QueueDeclare("task_queue", true, false, false, false, nil)
}
