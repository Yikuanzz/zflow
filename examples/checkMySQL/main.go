package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"zflow/internal/model"
)

type CheckMySQLOperation struct{}

func (op *CheckMySQLOperation) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	ip := string(input)

	address := ip + ":3306"
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return []byte("未检测到 MySQL 服务"), nil
	}
	defer conn.Close()
	return []byte("存在 MySQL 服务"), nil
}

type echoMsgOperations struct{}

func (op *echoMsgOperations) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	return input, nil
}

func main() {
	// 1. 创建节点类型
	nodeType_MySQL := model.NodeType{
		ID:        1,
		Category:  "test",
		Note:      "检查机器是否有mysql服务",
		Operation: &CheckMySQLOperation{},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "in", Label: "机器IP地址", PortType: "connection"},
			},
			"outputs": {
				{Name: "out", Label: "有mysql服务", PortType: "connection"},
				{Name: "out", Label: "没有mysql服务", PortType: "connection"},
			},
		},
	}

	nodeType_Echo := model.NodeType{
		ID:        2,
		Category:  "test",
		Note:      "回显节点",
		Operation: &echoMsgOperations{},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "in", Label: "信息", PortType: "connection"},
			},
		},
	}

	// 2. 创建连接类型
	connType := model.ConnectionType{
		ID:               1,
		Name:             "test_conn",
		Description:      "测试连接",
		Color:            "#FF0000",
		AllowedPortTypes: []string{"connection"},
	}

	// 3. 创建工作流
	wf := &model.Workflow{
		ID: "test_workflow",
		Dag: &model.Dag{
			Nodes: map[string]*model.Node{
				"node1": {
					ID:     "node1",
					TypeID: 1,
					Label:  "节点1",
				},
				"node2": {
					ID:     "node2",
					TypeID: 2,
					Label:  "节点2",
				},
				"node3": {
					ID:     "node3",
					TypeID: 2,
					Label:  "节点3",
				},
			},
			Connections: []model.Connection{
				{
					ID:     "conn1",
					TypeID: 1,
					From: model.Endpoint{
						NodeID:   "node1",
						PortName: "out",
					},
					To: model.Endpoint{
						NodeID:   "node2",
						PortName: "in",
					},
				},
				{
					ID:     "conn2",
					TypeID: 1,
					From: model.Endpoint{
						NodeID:   "node1",
						PortName: "out",
					},
					To: model.Endpoint{
						NodeID:   "node3",
						PortName: "in",
					},
				},
			},
		},
		NodeTypes: map[int]model.NodeType{
			1: nodeType_MySQL,
			2: nodeType_Echo,
		},
		ConnectionTypes: map[int]model.ConnectionType{
			1: connType,
		},
	}

	// 4. 验证工作流
	if err := wf.Validate(); err != nil {
		log.Fatalf("工作流验证失败: %v", err)
	}

	// 5. 创建执行上下文
	ctx := &model.ExecutionContext{
		Workflow: wf,
		Logger: func(msg string) {
			fmt.Println(msg)
		},
		Vars: make(map[string]interface{}),
	}

	// 6. 执行工作流
	if err := wf.ExecuteWorkflow(ctx); err != nil {
		log.Fatalf("工作流执行失败: %v", err)
	}

	// 7. 打印工作流状态
	wfJSON, _ := json.MarshalIndent(wf, "", "  ")
	fmt.Printf("\n工作流状态:\n%s\n", string(wfJSON))
}
