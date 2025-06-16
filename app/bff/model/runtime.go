package model

import (
	"fmt"
)

// Node 是一个具体实例，引用某个 NodeType
type Node struct {
	ID     string `json:"id"`
	TypeID string `json:"node_type"` // 对应 NodeType.ID
	Label  string `json:"label"`
	// 运行期字段 ↓↓↓
	State string `json:"-"` // running / success / failed ...
	// 存储每个端口的输入输出数据
	Inputs  map[string][]byte `json:"-"` // 端口名 -> 输入数据
	Outputs map[string][]byte `json:"-"` // 端口名 -> 输出数据
}

// Connection 表示有向边
type Connection struct {
	ID     string   `json:"connection_id"`
	TypeID string   `json:"connection_type"` // 对应 ConnectionType.ID
	From   Endpoint `json:"from"`
	To     Endpoint `json:"to"`
}

// Dag 持有拓扑结构，专注"图"层面的校验、遍历、拓扑排序等。
type Dag struct {
	Nodes       map[string]*Node
	Connections []Connection
}

// Workflow 则在更高一层，打包元数据 + DAG + 运行时映射关系
type Workflow struct {
	ID  string
	Dag *Dag

	// 元数据查字典
	NodeTypes       map[string]NodeType
	ConnectionTypes map[string]ConnectionType
}

// RawWorkflow 是 Workflow 的原始数据结构
type RawWorkflow struct {
	Nodes []struct {
		ID       string            `json:"id"`
		NodeType string            `json:"node_type"`
		Label    string            `json:"label"`
		Inputs   map[string][]byte `json:"inputs,omitempty"`
	} `json:"nodes"`
	Connections []struct {
		ID             string   `json:"connection_id"`
		ConnectionType string   `json:"connection_type"`
		From           Endpoint `json:"from"`
		To             Endpoint `json:"to"`
	} `json:"connections"`
}

// NewWorkflow 从 JSON 配置创建新的工作流实例
func NewWorkflow(uid string, raw RawWorkflow) (*Workflow, error) {
	wf := &Workflow{
		ID:              uid,
		Dag:             &Dag{Nodes: make(map[string]*Node)},
		NodeTypes:       make(map[string]NodeType),
		ConnectionTypes: make(map[string]ConnectionType),
	}

	// 3. 节点
	for _, n := range raw.Nodes {
		node := &Node{
			ID:     n.ID,
			TypeID: n.NodeType,
			Label:  n.Label,
			Inputs: make(map[string][]byte), // 初始化 Inputs map
		}

		// 如果有输入数据，复制到节点的 Inputs
		if n.Inputs != nil {
			for k, v := range n.Inputs {
				node.Inputs[k] = v
			}
		}

		wf.Dag.Nodes[n.ID] = node
	}

	// 4. 连接
	for _, c := range raw.Connections {
		conn := Connection{
			ID:     c.ID,
			TypeID: c.ConnectionType,
			From:   c.From,
			To:     c.To,
		}
		wf.Dag.Connections = append(wf.Dag.Connections, conn)
	}

	return wf, nil
}

// Validate 验证工作流的节点和连接是否有效
func (wf *Workflow) Validate() error {
	// 1. 验证基本结构
	if len(wf.Dag.Nodes) == 0 || len(wf.Dag.Connections) == 0 {
		return fmt.Errorf("workflow is required")
	}

	// 2. 验证节点
	for nodeID, node := range wf.Dag.Nodes {
		// 2.1 验证节点类型是否在工作流中定义
		nodeType, exists := wf.NodeTypes[node.TypeID]
		if !exists {
			return fmt.Errorf("node %s references undefined node type %s", nodeID, node.TypeID)
		}

		// 2.2 验证输入端口
		inputPorts := nodeType.Properties["inputs"]
		for inputName := range node.Inputs {
			found := false
			for _, port := range inputPorts {
				if port.Name == inputName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("节点 %s 的输入端口 %s 不存在", nodeID, inputName)
			}
		}
	}

	// 3. 验证连接
	for _, conn := range wf.Dag.Connections {
		// 3.1 验证连接类型是否在工作流中定义
		if _, exists := wf.ConnectionTypes[conn.TypeID]; !exists {
			return fmt.Errorf("connection %s references undefined connection type %s", conn.ID, conn.TypeID)
		}

		// 3.2 验证源节点和目标节点是否存在
		if _, exists := wf.Dag.Nodes[conn.From.NodeID]; !exists {
			return fmt.Errorf("connection %s references unknown source node %s", conn.ID, conn.From.NodeID)
		}
		if _, exists := wf.Dag.Nodes[conn.To.NodeID]; !exists {
			return fmt.Errorf("connection %s references unknown target node %s", conn.ID, conn.To.NodeID)
		}

		// 3.3 验证端口是否存在
		sourceNode := wf.Dag.Nodes[conn.From.NodeID]
		targetNode := wf.Dag.Nodes[conn.To.NodeID]
		sourceNodeType := wf.NodeTypes[sourceNode.TypeID]
		targetNodeType := wf.NodeTypes[targetNode.TypeID]

		// 验证源节点输出端口
		found := false
		for _, port := range sourceNodeType.Properties["outputs"] {
			if port.Name == conn.From.PortName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("connection %s references unknown output port %s in node %s",
				conn.ID, conn.From.PortName, conn.From.NodeID)
		}

		// 验证目标节点输入端口
		found = false
		for _, port := range targetNodeType.Properties["inputs"] {
			if port.Name == conn.To.PortName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("connection %s references unknown input port %s in node %s",
				conn.ID, conn.To.PortName, conn.To.NodeID)
		}
	}

	return nil
}

func (wf *Workflow) InjectOperations() {
}

// TopologicalSort 对 DAG 进行拓扑排序
func (wf *Workflow) TopologicalSort() ([]string, error) {
	// 构建邻接表
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// 初始化入度
	for nodeID := range wf.Dag.Nodes {
		inDegree[nodeID] = 0
	}

	// 构建图
	for _, conn := range wf.Dag.Connections {
		graph[conn.From.NodeID] = append(graph[conn.From.NodeID], conn.To.NodeID)
		inDegree[conn.To.NodeID]++
	}

	// 拓扑排序
	var result []string
	var queue []string

	// 将入度为 0 的节点加入队列
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
		}
	}

	// 处理队列
	for len(queue) > 0 {
		nodeID := queue[0]
		queue = queue[1:]
		result = append(result, nodeID)

		// 更新相邻节点的入度
		for _, neighbor := range graph[nodeID] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 检查是否有环
	if len(result) != len(wf.Dag.Nodes) {
		return nil, fmt.Errorf("workflow contains cycles")
	}

	return result, nil
}
