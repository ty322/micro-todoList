package wrappers

import (
	"api-gateway/services"
	"context"
	"strconv"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

func NewTask(id uint64, name string) *services.TaskModel {
	return &services.TaskModel{
		Id:         id,
		Title:      name,
		Content:    "响应超时",
		StartTime:  1000,
		EndTime:    1000,
		Status:     0,
		CreateTime: 1000,
		UpdateTime: 1000,
	}
}

// 降级函数
func DefaultTasks(resp interface{}) {
	models := make([]*services.TaskModel, 0)
	var i uint64
	for i = 0; i < 10; i++ {
		models = append(models, NewTask(i, "降级备忘录"+strconv.Itoa(20+int(i))))
	}
	result := resp.(services.TaskListResponse)
	result.TaskList = models
}

type TaskWrappr struct {
	client.Client
}

func (t *TaskWrappr) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	config := hystrix.CommandConfig{
		Timeout:                30000,
		RequestVolumeThreshold: 20,   // 熔断器请求阈值, 默认20, 意思是有20个请求才能有错误百分比判断
		ErrorPercentThreshold:  50,   // 错误百分比, 当错误超过百分比时, 直接进行降级处理
		SleepWindow:            5000, // 过多长时间, 熔断器再次检查是否开启, 单位毫秒(默认5秒)
	}
	hystrix.ConfigureCommand(cmdName, config)
	return hystrix.Do(cmdName, func() error {
		return t.Client.Call(ctx, req, resp)
	}, func(err error) error {
		return err
	})
}

// 初始化wrapper
func NewTaskWrapper(c client.Client) client.Client {
	return &TaskWrappr{c}
}
