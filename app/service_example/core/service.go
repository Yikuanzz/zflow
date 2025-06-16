package core

import (
	"context"
	"fmt"
	"log"
	v1 "zflow/api/base"
	registryCore "zflow/app/registry/core"
	"zflow/app/zflow/model"
	"zflow/utils/tool"
)

var (
	// 注册中心服务地址
	RegistryServiceAddr = registryCore.SERVICE_REGISTRY_ADDR

	// 服务名称
	ServiceName = registryCore.SERVICE_SERVICE_A
	// 服务地址
	ServiceAddr = registryCore.SERVICE_SERVICE_A_ADDR
)

// service 服务
type Service struct {
	v1.UnimplementedBaseServiceServer
}

// GetNodeTypes 获取节点类型
func (s *Service) GetNodeTypes(ctx context.Context, req *v1.GetNodeTypesRequest) (*v1.GetNodeTypesResponse, error) {
	// 获取所有节点类型
	nodeTypes := []*v1.NodeType{
		tool.ConvertNodeType(AddNodeType),
		tool.ConvertNodeType(MulNodeType),
		tool.ConvertNodeType(EchoNodeType),
	}

	return &v1.GetNodeTypesResponse{
		NodeTypes: nodeTypes,
	}, nil
}

// GetConnTypes 获取连接类型
func (s *Service) GetConnTypes(ctx context.Context, req *v1.GetConnTypesRequest) (*v1.GetConnTypesResponse, error) {
	// 获取所有连接类型
	connTypes := []*v1.ConnectionType{
		tool.ConvertConnType(DefaultConnType),
	}

	return &v1.GetConnTypesResponse{
		ConnectionTypes: connTypes,
	}, nil
}

// RunNode 运行节点
func (s *Service) RunNode(ctx context.Context, req *v1.RunNodeRequest) (*v1.RunNodeResponse, error) {
	// 根据节点ID找到对应的节点类型
	var nodeType model.NodeType
	switch req.NodeId {
	case fmt.Sprintf("%s.add.v1", ServiceName):
		nodeType = AddNodeType
	case fmt.Sprintf("%s.mul.v1", ServiceName):
		nodeType = MulNodeType
	case fmt.Sprintf("%s.echo.v1", ServiceName):
		nodeType = EchoNodeType
	default:
		return &v1.RunNodeResponse{
			State: "failed",
			Error: "未知的节点类型",
		}, nil
	}

	// 创建执行上下文
	execCtx := &model.ExecutionContext{
		Logger: func(msg string) {
			log.Printf("[%s] %s", req.NodeId, msg)
		},
		Vars: make(map[string]interface{}),
	}

	// 将请求中的变量复制到上下文
	for k, v := range req.Vars {
		execCtx.Vars[k] = v
	}

	// 执行节点操作
	outputs, err := nodeType.Operation.Execute(execCtx, req.Inputs, execCtx.Vars)
	if err != nil {
		return &v1.RunNodeResponse{
			State: "failed",
			Error: err.Error(),
		}, nil
	}

	return &v1.RunNodeResponse{
		Outputs: outputs,
		State:   "success",
	}, nil
}
