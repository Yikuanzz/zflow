package main

import (
	"context"
	"log"
	"strings"
	"time"
	v1 "zflow/api/base"
	"zflow/api/registry"
	registryCore "zflow/app/registry/core"

	"zflow/utils/cache"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// cache 节点类型缓存
var my_cache = cache.NewCache()

func main() {
	// 连接注册中心
	conn, err := grpc.NewClient(string(registryCore.SERVICE_REGISTRY_ADDR), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接注册中心失败: %v", err)
	}
	defer conn.Close()

	cli := registry.NewRegistryClient(conn)

	// 监听所有服务
	go watchAllServices(cli)

	// 等待服务发现完成
	time.Sleep(time.Second * 5)

	// 定期打印所有节点类型
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for range ticker.C {
			printAllNodeTypes(my_cache)
		}
	}()

	// 保持程序运行
	select {}
}

// watchAllServices 监听所有服务
func watchAllServices(cli registry.RegistryClient) {
	stream, err := cli.Watch(context.Background(), &registry.Query{})
	if err != nil {
		log.Fatalf("监听服务失败: %v", err)
	}

	for {
		srvList, err := stream.Recv()
		if err != nil {
			log.Printf("接收服务列表失败: %v", err)
			return
		}

		// 处理每个服务实例
		for _, inst := range srvList.Instances {
			go fetchServiceTypes(inst)
		}
	}
}

// fetchServiceTypes 获取服务的节点类型和连接类型
func fetchServiceTypes(inst *registry.ServiceInstance) {
	// 连接服务
	conn, err := grpc.NewClient(inst.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("连接服务 %s 失败: %v", inst.Name, err)
		return
	}
	defer conn.Close()

	// 创建客户端
	cli := v1.NewBaseServiceClient(conn)

	// 获取节点类型
	nodeTypes, err := cli.GetNodeTypes(context.Background(), &v1.GetNodeTypesRequest{})
	if err != nil {
		log.Printf("获取服务 %s 的节点类型失败: %v", inst.Name, err)
		return
	}

	// 缓存节点类型
	for _, nt := range nodeTypes.NodeTypes {
		my_cache.AddNodeType(inst.Name, nt)
	}

	// 获取连接类型
	connTypes, err := cli.GetConnTypes(context.Background(), &v1.GetConnTypesRequest{})
	if err != nil {
		log.Printf("获取服务 %s 的连接类型失败: %v", inst.Name, err)
		return
	}

	// 缓存连接类型
	for _, ct := range connTypes.ConnectionTypes {
		my_cache.AddConnType(inst.Name, ct)
	}

	log.Printf("服务 %s 的节点类型和连接类型已更新", inst.Name)
}

// printAllNodeTypes 打印所有节点类型
func printAllNodeTypes(cache *cache.Cache) {
	nodeTypes := cache.GetNodeTypes()
	connTypes := cache.GetConnTypes()

	log.Println("\n=== 节点类型 ===")
	for service, types := range nodeTypes {
		log.Printf("\n服务: %s", service)
		for _, nt := range types {
			log.Printf("- %s (%s): %s", nt.Uid, nt.Category, nt.Note)
		}
	}

	log.Println("\n=== 连接类型 ===")
	for service, types := range connTypes {
		log.Printf("\n服务: %s", service)
		for _, ct := range types {
			log.Printf("- %s: %s", ct.Uid, ct.Description)
		}
	}
	log.Println("\n" + strings.Repeat("-", 50)) // 添加分隔线
}
