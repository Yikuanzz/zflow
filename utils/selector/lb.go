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
	instances []*ServiceInstance
	index     uint64 // 用于轮询
	mu        sync.RWMutex
}

// NewLocalLB 创建本地负载均衡器
func NewLocalLB() *LocalLB {
	return &LocalLB{
		instances: make([]*ServiceInstance, 0),
		index:     0,
	}
}

// GetNextInstance 获取下一个服务实例（轮询）
func (lb *LocalLB) GetNextInstance() *ServiceInstance {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.instances) == 0 {
		return nil
	}

	// 原子操作获取下一个索引
	index := atomic.AddUint64(&lb.index, 1) % uint64(len(lb.instances))
	return lb.instances[index]
}

// GetInstanceCount 获取当前实例数量
func (lb *LocalLB) GetInstanceCount() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.instances)
}

// GetAllInstances 获取所有实例
func (lb *LocalLB) GetAllInstances() []*ServiceInstance {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances := make([]*ServiceInstance, len(lb.instances))
	copy(instances, lb.instances)
	return instances
}

// SetInstances 设置实例
func (lb *LocalLB) SetInstances(instances []*ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.instances = instances
}

// AddInstance 添加实例
func (lb *LocalLB) AddInstance(instance *ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.instances = append(lb.instances, instance)
}
