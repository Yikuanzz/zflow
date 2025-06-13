package model

import (
	"fmt"
)

// Node 是一个具体实例，引用某个 NodeType
type Node struct {
	ID     string `json:"id"`
	TypeID int    `json:"node_type"` // 对应 NodeType.ID
	Label  string `json:"label"`
	// 运行期字段 ↓↓↓
	State string `json:"-"` // running / success / failed ...
	// 存储每个端口的输入输出数据
	Inputs  map[string][]byte `json:"-"` // 端口名 -> 输入数据
	Outputs map[string][]byte `json:"-"` // 端口名 -> 输出数据
}

// NewNode 创建一个新的节点实例
func NewNode(id string, typeID int, label string) *Node {
	return &Node{
		ID:      id,
		TypeID:  typeID,
		Label:   label,
		State:   "pending",
		Inputs:  make(map[string][]byte),
		Outputs: make(map[string][]byte),
	}
}

// Connection 表示有向边
type Connection struct {
	ID     string   `json:"connection_id"`
	TypeID int      `json:"connection_type"` // 对应 ConnectionType.ID
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
	NodeTypes       map[int]NodeType
	ConnectionTypes map[int]ConnectionType
}

// NewWorkflow 从 JSON 配置创建新的工作流实例
func NewWorkflow(id string, config []byte) (*Workflow, error) {
	wf := &Workflow{
		ID:              id,
		Dag:             &Dag{Nodes: make(map[string]*Node)},
		NodeTypes:       make(map[int]NodeType),
		ConnectionTypes: make(map[int]ConnectionType),
	}

	// TODO: 解析 JSON 配置
	return wf, nil
}

// Validate 验证工作流的节点和连接是否有效
func (wf *Workflow) Validate() error {
	// 验证节点
	for _, node := range wf.Dag.Nodes {
		if _, exists := wf.NodeTypes[node.TypeID]; !exists {
			return fmt.Errorf("node %s references unknown node type %d", node.ID, node.TypeID)
		}
	}

	// 验证连接
	for _, conn := range wf.Dag.Connections {
		if _, exists := wf.ConnectionTypes[conn.TypeID]; !exists {
			return fmt.Errorf("connection %s references unknown connection type %d", conn.ID, conn.TypeID)
		}

		// 验证源节点和目标节点是否存在
		if _, exists := wf.Dag.Nodes[conn.From.NodeID]; !exists {
			return fmt.Errorf("connection %s references unknown source node %s", conn.ID, conn.From.NodeID)
		}
		if _, exists := wf.Dag.Nodes[conn.To.NodeID]; !exists {
			return fmt.Errorf("connection %s references unknown target node %s", conn.ID, conn.To.NodeID)
		}
	}

	return nil
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
