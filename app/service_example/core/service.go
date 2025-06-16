package core

import "fmt"

var (
	ServiceName = "service_example" // 服务名称
	ServiceAddr = "127.0.0.1:9090"  // 服务地址
)

// WrapUID 包装节点UID
func WrapUID(uid string, version string) string {
	return fmt.Sprintf("%s.%s.%s", ServiceName, uid, version)
}
