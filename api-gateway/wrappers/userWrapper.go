package wrappers

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type userWrapper struct {
	client.Client
}

func (u *userWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	config := hystrix.CommandConfig{
		Timeout:                30000,
		RequestVolumeThreshold: 20,   // 熔断器请求阈值, 默认20, 意思是有20个请求才能有错误百分比判断
		ErrorPercentThreshold:  50,   // 错误百分比, 当错误超过百分比时, 直接进行降级处理
		SleepWindow:            5000, // 过多长时间, 熔断器再次检查是否开启, 单位毫秒(默认5秒)
	}
	hystrix.ConfigureCommand(cmdName, config)
	return hystrix.Do(cmdName, func() error {
		return u.Client.Call(ctx, req, resp)
	}, func(err error) error {
		return err
	})
}

// 初始化wrapper
func NewUserWrapper(c client.Client) client.Client {
	return &userWrapper{c}
}
