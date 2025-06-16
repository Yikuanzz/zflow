package core

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	v1 "zflow/api/registry"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// registry 注册中心
type registry struct {
	v1.UnimplementedRegistryServer
	mu       sync.RWMutex
	services map[string]map[string]*serviceEntry // name -> id -> entry
}

// serviceEntry 服务实例
type serviceEntry struct {
	inst   *v1.ServiceInstance
	expire time.Time
}

// newRegistry 创建注册中心
func NewRegistry() *registry {
	r := &registry{
		services: make(map[string]map[string]*serviceEntry),
	}
	// 清理协程
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			r.sweep()
		}
	}()
	return r
}

// Register 注册服务
func (r *registry) Register(ctx context.Context, in *v1.ServiceInstance) (*v1.Lease, error) {
	if in.TtlSec <= 0 {
		in.TtlSec = 10
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	grp, ok := r.services[in.Name]
	if !ok {
		grp = make(map[string]*serviceEntry)
		r.services[in.Name] = grp
	}
	grp[in.Id] = &serviceEntry{inst: in, expire: time.Now().Add(time.Duration(in.TtlSec) * time.Second)}
	return lease(in), nil
}

// Deregister 注销服务
func (r *registry) Deregister(ctx context.Context, l *v1.Lease) (*emptypb.Empty, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if grp, ok := r.services[l.Name]; ok {
		delete(grp, l.Id)
	}
	return &emptypb.Empty{}, nil
}

// KeepAlive 续租
func (r *registry) KeepAlive(ctx context.Context, l *v1.Lease) (*v1.Lease, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if grp, ok := r.services[l.Name]; ok {
		if e, ok := grp[l.Id]; ok {
			e.expire = time.Now().Add(time.Duration(e.inst.TtlSec) * time.Second)
			return lease(e.inst), nil
		}
	}
	return nil, status.Error(codes.NotFound, "instance not found")
}

// Discover 查询服务
func (r *registry) Discover(ctx context.Context, q *v1.Query) (*v1.Services, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return &v1.Services{Instances: r.clone(q.Name)}, nil
}

// Watch 监听服务
func (r *registry) Watch(q *v1.Query, stream v1.Registry_WatchServer) error {
	// 轮询推送
	// TODO：可以优化为 notify chan
	ticker := time.NewTicker(time.Second * 5)
	last := ""
	for {
		select {
		case <-ticker.C:
			list := r.clone(q.Name)
			cur, _ := json.Marshal(list)
			if string(cur) != last { // 变更才推送
				if err := stream.Send(&v1.Services{Instances: list}); err != nil {
					return err
				}
				last = string(cur)
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

// sweep 定时剔除过期实例
func (r *registry) sweep() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	for n, grp := range r.services {
		for id, e := range grp {
			if e.expire.Before(now) {
				delete(grp, id)
			}
		}
		if len(grp) == 0 {
			delete(r.services, n)
		}
	}

}

// 复制一份快照
func (r *registry) clone(name string) []*v1.ServiceInstance {
	var out []*v1.ServiceInstance
	if name == "" {
		for _, grp := range r.services {
			for _, e := range grp {
				out = append(out, e.inst)
			}
		}
	} else if grp, ok := r.services[name]; ok {
		for _, e := range grp {
			out = append(out, e.inst)
		}
	}
	return out
}
