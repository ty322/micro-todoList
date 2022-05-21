package core

import (
	"context"
	"encoding/json"
	"errors"
	"task/model"
	services "task/service"

	"github.com/streadway/amqp"
)

// 创建备忘录, 将备忘录信息放到rabbitmq消息队列中
func (*TaskSerivce) CreateTask(ctx context.Context, req *services.TaskRequest, resp *services.TaskDetailResponse) error {
	ch, err := model.MQ.Channel()
	if err != nil {
		err = errors.New("rabbitmq err: " + err.Error())
	}
	q, _ := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	body, _ := json.Marshal(req) // title, content
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		err = errors.New("rabbitmq publish err: " + err.Error())
	}
	return nil
}

// 实现备忘录服务接口 获取备忘录列表
func (*TaskSerivce) GetTasksList(ctx context.Context, req *services.TaskRequest, resp *services.TaskListResponse) error {
	if req.Limit == 0 {
		req.Limit = 10
	}

	var taskData []model.Task
	var count uint32
	// 查找备忘录
	err := model.DB.Offset(req.Start).Limit(req.Limit).Where("uid=?", req.Uid).Find(&taskData).Error
	if err != nil {
		return errors.New("mysql find: err: " + err.Error())
	}
	// 统计数量
	model.DB.Model(&model.Task{}).Where("uid=?", req.Uid).Count(&count)
	// 返回proto里面定义的类型
	var taskRes []*services.TaskModel
	for _, item := range taskData {
		taskRes = append(taskRes, BuildTask(item))
	}

	resp.TaskList = taskRes
	resp.Count = count
	return nil
}

// 获取备忘录的详细信息
func (*TaskSerivce) GetTask(ctx context.Context, req *services.TaskRequest, resp *services.TaskDetailResponse) error {
	taskData := model.Task{}
	model.DB.First(&taskData, req.Id)
	taskRes := BuildTask(taskData)
	resp.TaskDetail = taskRes
	return nil
}

// 修改备忘录
func (*TaskSerivce) UpdateTask(Taskctx context.Context, req *services.TaskRequest, resp *services.TaskDetailResponse) error {
	taskData := model.Task{}
	// 查找该用户的备忘录
	model.DB.Model(&model.Task{}).Where("id= ? AND uid= ?", req.Id, req.Uid).First(&taskData)
	taskData.Title = req.Title
	taskData.Status = int(req.Status)
	taskData.Content = req.Content
	model.DB.Save(&taskData)
	resp.TaskDetail = BuildTask(taskData)
	return nil
}

// 删除备忘录
func (*TaskSerivce) DeleteTask(ctx context.Context, req *services.TaskRequest, resp *services.TaskDetailResponse) error {
	model.DB.Model(&model.Task{}).Where("id=? AND uid=?", req.Id, req.Uid).Delete(&model.Task{})
	return nil
}
