package model

import (
	"encoding/json"
	"fmt"
)

// Operation 与 Node 一一对应，真正执行工作。
type Operation interface {
	Execute(context Context, input []byte, vars map[string]interface{}) ([]byte, error)
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

	// 2. 按顺序执行节点
	for _, nodeID := range order {
		node := wf.Dag.Nodes[nodeID]
		nodeType := wf.NodeTypes[node.TypeID]

		// 收集输入
		input, err := wf.collectNodeInputs(nodeID)
		if err != nil {
			return fmt.Errorf("failed to collect inputs for node %s: %v", nodeID, err)
		}

		// 执行操作
		_, err = nodeType.Operation.Execute(ctx, input, ctx.Vars)
		if err != nil {
			return fmt.Errorf("node %s execution failed: %v", nodeID, err)
		}

		// 更新节点状态
		node.State = "success"
	}

	return nil
}

// collectNodeInputs 收集节点的所有输入
func (wf *Workflow) collectNodeInputs(nodeID string) ([]byte, error) {
	node := wf.Dag.Nodes[nodeID]
	nodeType := wf.NodeTypes[node.TypeID]

	// 获取输入端口列表
	inputPorts, exists := nodeType.Properties["inputs"]
	if !exists {
		return nil, fmt.Errorf("node type %d has no input ports defined", node.TypeID)
	}

	// 收集每个输入端口的数据
	inputs := make(map[string]interface{})
	for _, port := range inputPorts {
		// 查找连接到该端口的连接
		for _, conn := range wf.Dag.Connections {
			if conn.To.NodeID == nodeID && conn.To.PortName == port.Name {
				// TODO: 从源节点获取输出数据
				// 这里需要实现从上游节点获取数据的逻辑
				inputs[port.Name] = nil // 临时占位
			}
		}
	}

	// 将输入数据序列化为 JSON
	inputJSON, err := json.Marshal(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inputs: %v", err)
	}

	return inputJSON, nil
}

// EchoOperation 是一个简单的回显操作，用于测试
type EchoOperation struct {
	Message string
}

func (op *EchoOperation) Execute(ctx Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	ctx.Log(fmt.Sprintf("Echo: %s", op.Message))
	return []byte(op.Message), nil
}
