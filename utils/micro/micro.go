package thing

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zflow/api/registry"
	"zflow/app/service_example/core"
	"zflow/app/zflow/model"
	"zflow/utils/service"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "zflow/api/base"
)

// Micro 微服务
type Micro struct {
	baseService     *service.BaseService
	serviceInstance *registry.ServiceInstance
	registryClient  registry.RegistryClient
}

// NewMicro 创建微服务
func NewMicro(serviceName string, serviceAddr string, nodeTypes map[string]*model.NodeType, connTypes map[string]*model.ConnectionType) *Micro {
	// 创建基础服务
	baseService := &service.BaseService{
		Name:      serviceName,
		Addr:      serviceAddr,
		NodeTypes: nodeTypes,
		ConnTypes: connTypes,
	}

	return &Micro{
		baseService: baseService,
	}
}

// Run 运行微服务
func (m *Micro) Run() {
	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", m.baseService.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	v1.RegisterBaseServiceServer(s, m.baseService)

	// 启动服务注册
	go m.registerService()

	// 设置优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 gRPC 服务
	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("Shutting down server...")

	// 注销服务
	if err := m.unregisterService(); err != nil {
		log.Printf("Error unregistering service: %v", err)
	}

	// 优雅关闭 gRPC 服务器
	s.GracefulStop()
	log.Println("Server stopped")
}

// registerService 注册服务到注册中心
func (m *Micro) registerService() {
	conn, err := grpc.NewClient(string(core.RegistryServiceAddr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to registry: %v", err)
	}
	defer conn.Close()

	m.registryClient = registry.NewRegistryClient(conn)

	m.serviceInstance = &registry.ServiceInstance{
		Name: string(core.ServiceName),
		Id:   uuid.New().String(),
		Addr: string(core.ServiceAddr),
		Meta: map[string]string{
			"version": "v1.0.0", // 版本号
		},
		TtlSec: 10, // 10秒后过期
	}

	// 注册服务
	lease, err := m.registryClient.Register(context.Background(), m.serviceInstance)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	// 心跳协程
	tk := time.NewTicker(time.Second * 5)
	for range tk.C {
		if _, err := m.registryClient.KeepAlive(context.Background(), lease); err != nil {
			log.Printf("keepalive failed: %v, re-registering...", err)
			lease, _ = m.registryClient.Register(context.Background(), m.serviceInstance)
		}
	}
}

// unregisterService 注销服务
func (m *Micro) unregisterService() error {
	if m.registryClient == nil || m.serviceInstance == nil {
		return nil
	}

	// 注销服务
	_, err := m.registryClient.Deregister(context.Background(), &registry.Lease{
		Name: string(core.ServiceName),
		Id:   m.serviceInstance.Id,
	})
	if err != nil {
		return err
	}

	return nil
}
