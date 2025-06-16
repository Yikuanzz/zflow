package main

import (
	"context"
	"log"
	"time"
	"zflow/api/registry"
	v1 "zflow/api/s_example"
	registryCore "zflow/app/registry/core"
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
	// 启动服务发现
	go watch()

	// 等待服务发现完成
	time.Sleep(time.Second * 10)

	// 测试服务调用
	testService()

	// select {}
}

func watch() {
	conn, err := grpc.NewClient(string(registryCore.SERVICE_REGISTRY_ADDR), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接注册中心失败: %v", err)
	}
	defer conn.Close()

	cli := registry.NewRegistryClient(conn)

	// 监听服务
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

// testService 测试服务调用
func testService() {
	// 获取服务实
	instance := lb.GetNextInstance()
	if instance == nil {
		log.Fatal("没有可用的服务实例")
	}

	// 连接服务
	conn, err := grpc.NewClient(instance.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务失败: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	cli := v1.NewS_ExampleServiceClient(conn)

	// 测试 GetNodeTypes
	log.Println("测试 GetNodeTypes...")
	nodeTypes, err := cli.GetNodeTypes(context.Background(), &v1.GetNodeTypesRequest{})
	if err != nil {
		log.Fatalf("获取节点类型失败: %v", err)
		return
	}
	log.Printf("获取到 %d 个节点类型:", len(nodeTypes.NodeTypes))
	for _, nt := range nodeTypes.NodeTypes {
		log.Printf("- %s (%s): %s", nt.Uid, nt.Category, nt.Note)
	}

	// 测试 GetConnTypes
	log.Println("\n测试 GetConnTypes...")
	connTypes, err := cli.GetConnTypes(context.Background(), &v1.GetConnTypesRequest{})
	if err != nil {
		log.Fatalf("获取连接类型失败: %v", err)
		return
	}
	log.Printf("获取到 %d 个连接类型:", len(connTypes.ConnectionTypes))
	for _, ct := range connTypes.ConnectionTypes {
		log.Printf("- %s: %s", ct.Uid, ct.Description)
	}

	// 测试 RunNode (加法节点)
	log.Println("\n测试 RunNode (加法)...")
	addResp, err := cli.RunNode(context.Background(), &v1.RunNodeRequest{
		NodeId: "s_example.add.v1",
		Inputs: map[string][]byte{
			"a": []byte("10"),
			"b": []byte("20"),
		},
	})
	if err != nil {
		log.Fatalf("运行加法节点失败: %v", err)
		return
	}
	log.Printf("加法结果: %s", string(addResp.Outputs["sum"]))

	// 测试 RunNode (乘法节点)
	log.Println("\n测试 RunNode (乘法)...")
	mulResp, err := cli.RunNode(context.Background(), &v1.RunNodeRequest{
		NodeId: "s_example.mul.v1",
		Inputs: map[string][]byte{
			"a": []byte("10"),
			"b": []byte("20"),
		},
	})
	if err != nil {
		log.Fatalf("运行乘法节点失败: %v", err)
		return
	}
	log.Printf("乘法结果: %s", string(mulResp.Outputs["product"]))

	// 测试 RunNode (回显节点)
	log.Println("\n测试 RunNode (回显)...")
	echoResp, err := cli.RunNode(context.Background(), &v1.RunNodeRequest{
		NodeId: "s_example.echo.v1",
		Inputs: map[string][]byte{
			"input": []byte("Hello, World!"),
		},
	})
	if err != nil {
		log.Fatalf("运行回显节点失败: %v", err)
		return
	}
	log.Printf("回显结果: %s", string(echoResp.Outputs["output"]))

	conn.Close()
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
