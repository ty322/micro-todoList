package main

import (
	"task/conf"
	"task/core"
	services "task/service"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	conf.Init("")
	// etcd 注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcTaskService"), // 微服务名字
		micro.Address("127.0.0.1:8083"),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数, 初始化
	microService.Init()

	// 服务注册
	_ = services.RegisterTaskServiceHandler(microService.Server(), new(core.TaskSerivce))

	// 启动微服务
	_ = microService.Run()
}
