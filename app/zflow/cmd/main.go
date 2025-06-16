package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zflow/api/registry"
	registryCore "zflow/app/registry/core"
	"zflow/app/zflow/server"
	"zflow/utils/selector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID   string
	Addr string
	Meta map[string]string
}

// lb 本地负载均衡器
var lb = selector.NewLocalLB()

func main() {
	// 新建服务
	server := server.NewServer()

	// 监听服务
	go watch()

	// 启动服务器
	go func() {
		log.Printf("服务器正在启动，监听端口 %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭失败:", err)
	}

	log.Println("服务器已关闭")
}

// watch 监听服务
func watch() {
	conn, err := grpc.NewClient(string(registryCore.SERVICE_REGISTRY_ADDR), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接注册中心失败: %v", err)
	}
	defer conn.Close()

	cli := registry.NewRegistryClient(conn)

	stream, err := cli.Watch(context.Background(), &registry.Query{Name: string(registryCore.SERVICE_S_EXAMPLE)})
	if err != nil {
		log.Fatalf("监听服务失败: %v", err)
	}

	for {
		srvList, err := stream.Recv()
		if err != nil {
			log.Printf("接收服务列表失败: %v", err)
			return
		}
		refreshLocalCache(srvList.Instances)
		log.Printf("服务列表更新: %d 个实例", len(srvList.Instances))
	}
}

// refreshLocalCache 更新本地缓存
func refreshLocalCache(instances []*registry.ServiceInstance) {
	// 清空现有实例
	lb.SetInstances(make([]*selector.ServiceInstance, 0, len(instances)))

	// 添加新实例
	for _, inst := range instances {
		lb.AddInstance(&selector.ServiceInstance{
			ID:   inst.Id,
			Addr: inst.Addr,
			Meta: inst.Meta,
		})
	}
}
