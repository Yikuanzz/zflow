package main

import (
	registryCore "zflow/app/registry/core"
	"zflow/app/service_example/core"
	"zflow/utils/micro"
)

func main() {
	// 创建微服务
	micro := micro.NewMicro(
		registryCore.SERVICE_REGISTRY_ADDR, // 服务注册中心地址
		core.ServiceName,                   // 服务名称
		core.ServiceAddr,                   // 服务地址
		core.NodeTypes,                     // 节点类型
		core.ConnTypes,                     // 连接类型
	)

	// 运行微服务
	micro.Run()
}
