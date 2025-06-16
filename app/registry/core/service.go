package core

// ServiceName 服务名称
type ServiceName string

var (
	// 示例服务
	SERVICE_S_EXAMPLE = ServiceName("s_example")
	// 服务A
	SERVICE_SERVICE_A = ServiceName("service_a")
)

// ServiceAddr 服务地址
type ServiceAddr string

var (
	// 注册中心地址
	SERVICE_REGISTRY_ADDR = ServiceAddr("127.0.0.1:50051")
	// 示例服务地址
	SERVICE_S_EXAMPLE_ADDR = ServiceAddr("127.0.0.1:9090")
	// 服务A地址
	SERVICE_SERVICE_A_ADDR = ServiceAddr("127.0.0.1:9091")
)
