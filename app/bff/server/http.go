package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"zflow/api/registry"
	"zflow/app/bff/global"
	"zflow/app/bff/model"
	"zflow/utils/selector"

	v1 "zflow/api/base"
	registryCore "zflow/app/registry/core"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer() *http.Server {
	// 连接注册中心
	conn, err := grpc.NewClient(registryCore.SERVICE_REGISTRY_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接注册中心失败: %v", err)
	}
	defer conn.Close()

	cli := registry.NewRegistryClient(conn)

	// 监听所有服务
	go watchAllServices(cli)

	router := gin.Default()

	// 获取所有节点类型
	router.GET("/node_types", func(c *gin.Context) {
		c.JSON(http.StatusOK, global.Cache.GetNodeTypes())
	})

	// 获取所有连接类型
	router.GET("/connection_types", func(c *gin.Context) {
		c.JSON(http.StatusOK, global.Cache.GetConnTypes())
	})

	// 执行工作流
	router.POST("/workflows", func(c *gin.Context) {
		var req struct {
			UID      string            `json:"uid"`
			Workflow model.RawWorkflow `json:"workflow"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.UID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
			return
		}
		if req.Workflow.Nodes == nil || req.Workflow.Connections == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "workflow is required"})
			return
		}

		// 1、创建工作流
		wf, err := model.NewWorkflow(req.UID, req.Workflow)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2、注入节点和连接类型的操作
		// if err := wf.InjectOperations(); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// 3、创建执行上下文
		ctx := &model.ExecutionContext{
			Workflow: wf,
			Logger: func(msg string) {
				fmt.Println(msg)
			},
			Vars: make(map[string]interface{}),
		}

		// 4、执行工作流
		if err := wf.ExecuteWorkflow(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 5、收集工作流执行结果
		result := wf.CollectWorkflowResults()

		c.JSON(http.StatusOK, result)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
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

		// 更新负载均衡器中的服务实例
		updateLoadBalancer(srvList.Instances)

		// 处理每个服务实例
		for _, inst := range srvList.Instances {
			go fetchServiceTypes(inst)
		}
	}
}

// updateLoadBalancer 更新负载均衡器中的服务实例
func updateLoadBalancer(instances []*registry.ServiceInstance) {
	// 按服务名分组
	serviceGroups := make(map[string][]*selector.ServiceInstance)
	for _, inst := range instances {
		serviceInstance := &selector.ServiceInstance{
			ID:   inst.Id,
			Addr: inst.Addr,
			Meta: inst.Meta,
		}
		serviceGroups[inst.Name] = append(serviceGroups[inst.Name], serviceInstance)
	}

	// 更新负载均衡器
	for serviceName, instances := range serviceGroups {
		global.LoadBalance.SetInstances(serviceName, instances)
		log.Printf("服务 %s 的实例已更新到负载均衡器，共 %d 个实例", serviceName, len(instances))
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
		global.Cache.AddNodeType(inst.Name, nt)
	}

	// 获取连接类型
	connTypes, err := cli.GetConnTypes(context.Background(), &v1.GetConnTypesRequest{})
	if err != nil {
		log.Printf("获取服务 %s 的连接类型失败: %v", inst.Name, err)
		return
	}

	// 缓存连接类型
	for _, ct := range connTypes.ConnectionTypes {
		global.Cache.AddConnType(inst.Name, ct)
	}

	log.Printf("服务 %s 的节点类型和连接类型已更新", inst.Name)
}
