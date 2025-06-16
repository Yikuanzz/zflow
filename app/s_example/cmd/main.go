package main

import (
	"context"
	"log"
	"time"
	"zflow/api/registry"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// REGISTRY_SERVICE_ADDR 注册中心地址
var REGISTRY_SERVICE_ADDR = "127.0.0.1:50051"

// S_EXAMPLE_SERVICE_ADDR 服务地址
var S_EXAMPLE_SERVICE_ADDR = "127.0.0.1:9090"

func main() {
	conn, _ := grpc.NewClient(REGISTRY_SERVICE_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cli := registry.NewRegistryClient(conn)

	inst := &registry.ServiceInstance{
		Name: "s_example",
		Id:   uuid.New().String(),
		Addr: S_EXAMPLE_SERVICE_ADDR,
		Meta: map[string]string{
			"version": "v1.0.0", // 版本号
		},
		TtlSec: 10, // 10秒后过期
	}
	lease, err := cli.Register(context.Background(), inst)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	// 心跳协程
	go func() {
		tk := time.NewTicker(time.Second * 5)
		for range tk.C {
			if _, err := cli.KeepAlive(context.Background(), lease); err != nil {
				log.Printf("keepalive failed: %v, re-registering...", err)
				lease, _ = cli.Register(context.Background(), inst)
			}
		}
	}()

	// gRPC 服务启动
	// lis, err := net.Listen("tcp", S_EXAMPLE_SERVICE_ADDR)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()

}
