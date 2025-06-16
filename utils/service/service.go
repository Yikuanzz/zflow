package service

import (
	"context"
	"log"
	v1 "zflow/api/base"
	"zflow/app/zflow/model"
	"zflow/utils/tool"
)

// BaseService 基础服务
type BaseService struct {
	v1.UnimplementedBaseServiceServer
	Name      string
	Addr      string
	NodeTypes map[string]*model.NodeType
	ConnTypes map[string]*model.ConnectionType
}

// GetNodeTypes 获取节点类型
func (s *BaseService) GetNodeTypes(ctx context.Context, req *v1.GetNodeTypesRequest) (*v1.GetNodeTypesResponse, error) {
	var nodeTypes []*v1.NodeType
	for _, nodeType := range s.NodeTypes {
		nodeTypes = append(nodeTypes, tool.ConvertNodeType(*nodeType))
	}

	return &v1.GetNodeTypesResponse{
		NodeTypes: nodeTypes,
	}, nil
}

// GetConnTypes 获取连接类型
func (s *BaseService) GetConnTypes(ctx context.Context, req *v1.GetConnTypesRequest) (*v1.GetConnTypesResponse, error) {
	var connTypes []*v1.ConnectionType
	for _, connType := range s.ConnTypes {
		connTypes = append(connTypes, tool.ConvertConnType(*connType))
	}

	return &v1.GetConnTypesResponse{
		ConnectionTypes: connTypes,
	}, nil
}

// RunNode 运行节点
func (s *BaseService) RunNode(ctx context.Context, req *v1.RunNodeRequest) (*v1.RunNodeResponse, error) {
	// 根据节点ID找到对应的节点类型
	var nodeType *model.NodeType
	for _, nt := range s.NodeTypes {
		if nt.UID == req.NodeId {
			nodeType = nt
			break
		}
	}

	if nodeType == nil {
		return &v1.RunNodeResponse{
			State: "failed",
			Error: "未找到节点类型",
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
