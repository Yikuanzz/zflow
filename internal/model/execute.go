package model

import (
	"fmt"
)

// Operation 与 Node 一一对应，真正执行工作。
type Operation interface {
	Execute(context Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error)
}

// Context 是运行期给 Operation 的最少上下文
type Context interface {
	Log(msg string)
}

// ExecutionContext 实现 Context 接口，提供完整的执行上下文
type ExecutionContext struct {
	Workflow *Workflow
	Logger   func(msg string)
	Vars     map[string]interface{}
}

func (ctx *ExecutionContext) Log(msg string) {
	if ctx.Logger != nil {
		ctx.Logger(msg)
	}
}

// ExecuteWorkflow 执行整个工作流
func (wf *Workflow) ExecuteWorkflow(ctx *ExecutionContext) error {
	// 1. 获取拓扑排序
	order, err := wf.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to sort workflow: %v", err)
	}

	ctx.Log(fmt.Sprintf("工作流执行顺序: %v", order))

	// 2. 按顺序执行节点
	for _, nodeID := range order {
		node := wf.Dag.Nodes[nodeID]
		nodeType := wf.NodeTypes[node.TypeID]

		ctx.Log(fmt.Sprintf("开始执行节点 %s (%s)", nodeID, node.Label))

		// 收集输入
		err := wf.collectNodeInputs(nodeID)
		if err != nil {
			node.State = "failed"
			return fmt.Errorf("failed to collect inputs for node %s: %v", nodeID, err)
		}

		// 执行操作
		outputs, err := nodeType.Operation.Execute(ctx, node.Inputs, ctx.Vars)
		if err != nil {
			node.State = "failed"
			return fmt.Errorf("node %s execution failed: %v", nodeID, err)
		}

		// 存储输出
		node.Outputs = outputs
		node.State = "success"
		ctx.Log(fmt.Sprintf("节点 %s 执行完成，状态: %s", nodeID, node.State))
	}

	return nil
}

// collectNodeInputs 收集节点的所有输入
func (wf *Workflow) collectNodeInputs(nodeID string) error {
	node := wf.Dag.Nodes[nodeID]
	nodeType := wf.NodeTypes[node.TypeID]

	// 获取输入端口列表
	inputPorts, exists := nodeType.Properties["inputs"]
	if !exists || len(inputPorts) == 0 {
		// 如果没有输入端口，直接返回
		return nil
	}

	// 收集每个输入端口的数据
	for _, port := range inputPorts {
		// 如果节点已经有预设的输入数据，则跳过
		if _, exists := node.Inputs[port.Name]; exists {
			continue
		}

		// 查找连接到该端口的连接
		found := false
		for _, conn := range wf.Dag.Connections {
			if conn.To.NodeID == nodeID && conn.To.PortName == port.Name {
				// 从源节点获取输出数据
				sourceNode := wf.Dag.Nodes[conn.From.NodeID]
				if sourceNode.State == "success" {
					if output, exists := sourceNode.Outputs[conn.From.PortName]; exists && output != nil {
						node.Inputs[port.Name] = output
						found = true
						break
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("node %s 的输入端口 %s 没有找到对应的连接或源节点未执行完成", nodeID, port.Name)
		}
	}

	return nil
}

// CollectWorkflowResults 收集工作流执行结果
func (wf *Workflow) CollectWorkflowResults() map[string]interface{} {
	result := map[string]interface{}{
		"workflow_id": wf.ID,
		"status":      "success",
		"nodes":       make(map[string]map[string]interface{}),
	}

	// 遍历所有节点，收集执行结果
	for nodeID, node := range wf.Dag.Nodes {
		nodeResult := map[string]interface{}{
			"id":    nodeID,
			"label": node.Label,
			"state": node.State,
		}

		// 收集输入数据
		if len(node.Inputs) > 0 {
			inputs := make(map[string]string)
			for port, data := range node.Inputs {
				inputs[port] = string(data)
			}
			nodeResult["inputs"] = inputs
		}

		// 收集输出数据
		if len(node.Outputs) > 0 {
			outputs := make(map[string]string)
			for port, data := range node.Outputs {
				outputs[port] = string(data)
			}
			nodeResult["outputs"] = outputs
		}

		result["nodes"].(map[string]map[string]interface{})[nodeID] = nodeResult
	}

	return result
}
