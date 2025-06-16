package data

import (
	"fmt"
	"zflow/app/zflow/model"
)

// InjectOperations 验证并注入 NodeType 和 ConnectionType 的操作
func InjectOperations(wf *model.Workflow) error {
	// 获取默认的类型映射
	nodeTypeMap := GetDefaultNodeTypes()
	connTypeMap := GetDefaultConnectionTypes()

	// 验证并注入 NodeType
	for _, node := range wf.Dag.Nodes {
		// 检查节点类型是否存在
		nt, ok := nodeTypeMap[node.TypeID]
		if !ok {
			return fmt.Errorf("未知的节点类型: %s", node.TypeID)
		}

		// 注入完整的节点类型定义
		wf.NodeTypes[node.TypeID] = nt
	}

	// 验证并注入 ConnectionType
	for _, conn := range wf.Dag.Connections {
		// 检查连接类型是否存在
		ct, ok := connTypeMap[conn.TypeID]
		if !ok {
			return fmt.Errorf("未知的连接类型: %s", conn.TypeID)
		}

		// 注入完整的连接类型定义
		wf.ConnectionTypes[conn.TypeID] = ct
	}

	return nil
}
