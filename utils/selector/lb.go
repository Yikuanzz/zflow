package selector

import (
	"sync"
	"sync/atomic"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID   string
	Addr string
	Meta map[string]string
}

// LocalLB 本地负载均衡器
type LocalLB struct {
	instances map[string][]*ServiceInstance
	index     uint64 // 用于轮询
	mu        sync.RWMutex
}

// NewLocalLB 创建本地负载均衡器
func NewLocalLB() *LocalLB {
	return &LocalLB{
		instances: make(map[string][]*ServiceInstance),
		index:     0,
	}
}

// GetNextInstance 获取下一个服务实例（轮询）
func (lb *LocalLB) GetNextInstance(serviceName string) *ServiceInstance {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances := lb.instances[serviceName]
	if len(instances) == 0 {
		return nil
	}

	// 原子操作获取下一个索引
	index := atomic.AddUint64(&lb.index, 1) % uint64(len(instances))
	return instances[index]
}

// GetInstanceCount 获取当前实例数量
func (lb *LocalLB) GetInstanceCount(serviceName string) int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.instances[serviceName])
}

// GetAllInstances 获取所有实例
func (lb *LocalLB) GetAllInstances(serviceName string) []*ServiceInstance {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances := lb.instances[serviceName]
	if len(instances) == 0 {
		return nil
	}

	result := make([]*ServiceInstance, len(instances))
	copy(result, instances)
	return result
}

// SetInstances 设置实例
func (lb *LocalLB) SetInstances(serviceName string, instances []*ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if instances == nil {
		instances = make([]*ServiceInstance, 0)
	}
	lb.instances[serviceName] = instances
}

// AddInstance 添加实例
func (lb *LocalLB) AddInstance(serviceName string, instance *ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if instance == nil {
		return
	}

	instances := lb.instances[serviceName]
	lb.instances[serviceName] = append(instances, instance)
}
