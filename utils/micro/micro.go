package micro

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zflow/api/registry"
	"zflow/app/zflow/model"
	"zflow/utils/service"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "zflow/api/base"
)

// Micro 微服务
type Micro struct {
	registryServiceAddr string
	baseService         *service.BaseService
	serviceInstance     *registry.ServiceInstance
	registryClient      registry.RegistryClient
	grpcConn            *grpc.ClientConn
}

// NewMicro 创建微服务
func NewMicro(registryServiceAddr, serviceName, serviceAddr string, nodeTypes map[string]*model.NodeType, connTypes map[string]*model.ConnectionType) *Micro {
	// 创建基础服务
	baseService := &service.BaseService{
		Name:      serviceName,
		Addr:      serviceAddr,
		NodeTypes: nodeTypes,
		ConnTypes: connTypes,
	}

	return &Micro{
		registryServiceAddr: registryServiceAddr,
		baseService:         baseService,
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

	// 关闭 gRPC 连接
	if m.grpcConn != nil {
		m.grpcConn.Close()
	}

	// 优雅关闭 gRPC 服务器
	s.GracefulStop()
	log.Println("Server stopped")
}

// registerService 注册服务到注册中心
func (m *Micro) registerService() {
	// 创建 gRPC 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, m.registryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to registry: %v", err)
		return
	}
	m.grpcConn = conn
	m.registryClient = registry.NewRegistryClient(conn)

	// 创建服务实例
	m.serviceInstance = &registry.ServiceInstance{
		Name: m.baseService.Name,
		Id:   uuid.New().String(),
		Addr: m.baseService.Addr,
		Meta: map[string]string{
			"version": "v1.0.0",
		},
		TtlSec: 10,
	}

	// 注册服务
	lease, err := m.registryClient.Register(context.Background(), m.serviceInstance)
	if err != nil {
		log.Printf("注册失败: %v", err)
		return
	}

	// 心跳协程
	go func() {
		tk := time.NewTicker(time.Second * 5)
		defer tk.Stop()

		for {
			select {
			case <-tk.C:
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				_, err := m.registryClient.KeepAlive(ctx, lease)
				cancel()

				if err != nil {
					log.Printf("keepalive failed: %v, re-registering...", err)
					// 重新注册
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					newLease, err := m.registryClient.Register(ctx, m.serviceInstance)
					cancel()

					if err != nil {
						log.Printf("re-register failed: %v", err)
						continue
					}
					lease = newLease
				}
			}
		}
	}()
}

// unregisterService 注销服务
func (m *Micro) unregisterService() error {
	if m.registryClient == nil || m.serviceInstance == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 注销服务
	_, err := m.registryClient.Deregister(ctx, &registry.Lease{
		Name: m.baseService.Name,
		Id:   m.serviceInstance.Id,
	})
	return err
}
