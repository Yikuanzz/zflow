package global

import (
	"zflow/utils/cache"
	"zflow/utils/selector"
)

// Cache 缓存存储节点类型和连接类型
var Cache *cache.Cache

// LoadBalance 负载均衡
var LoadBalance *selector.LocalLB

func init() {
	// Cache 初始化缓存
	Cache = cache.NewCache()

	// LoadBalance 初始化负载均衡
	LoadBalance = selector.NewLocalLB()
}
