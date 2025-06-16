package main

import (
	"context"
	"log"
	"time"

	"zflow/api/registry"
	registryCore "zflow/app/registry/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接注册中心
	conn, err := grpc.NewClient(string(registryCore.SERVICE_REGISTRY_ADDR), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接注册中心失败: %v", err)
	}
	defer conn.Close()

	cli := registry.NewRegistryClient(conn)

	// 测试服务注册
	testRegister(cli)

	// 测试服务发现
	testWatch(cli)

	// 保持程序运行
	select {}
}

// testRegister 测试服务注册
func testRegister(cli registry.RegistryClient) {
	// 注册 s_example 服务
	inst := &registry.ServiceInstance{
		Name: "service_zz",
		Id:   "test-instance-1",
		Addr: "127.0.0.1:50051",
		Meta: map[string]string{
			"version": "v1.0.0",
		},
		TtlSec: 10,
	}

	lease, err := cli.Register(context.Background(), inst)
	if err != nil {
		log.Fatalf("注册服务失败: %v", err)
	}
	log.Printf("服务注册成功，租约ID: %s", lease.Id)

	// 启动心跳
	go func() {
		tk := time.NewTicker(time.Second * 5)
		defer tk.Stop()
		for range tk.C {
			if _, err := cli.KeepAlive(context.Background(), lease); err != nil {
				log.Printf("心跳失败: %v", err)
				return
			}
			log.Println("心跳成功")
		}
	}()
}

// testWatch 测试服务发现
func testWatch(cli registry.RegistryClient) {
	stream, err := cli.Watch(context.Background(), &registry.Query{
		Name: "service_zz",
	})
	if err != nil {
		log.Fatalf("监听服务失败: %v", err)
	}

	go func() {
		for {
			srvList, err := stream.Recv()
			if err != nil {
				log.Printf("接收服务列表失败: %v", err)
				return
			}
			log.Printf("服务列表更新: %d 个实例", len(srvList.Instances))
			for _, inst := range srvList.Instances {
				log.Printf("- 实例ID: %s, 地址: %s", inst.Id, inst.Addr)
			}
		}
	}()
}
