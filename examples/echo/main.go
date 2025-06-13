package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"zflow/internal/data"
// 	"zflow/internal/model"
// )

// func main() {
// 	// 1. 创建节点类型
// 	nodeType := model.NodeType{
// 		ID:        1,
// 		Category:  "test",
// 		Note:      "测试节点",
// 		Operation: &data.EchoOperation{Message: "Hello from Node!"},
// 		Properties: map[string][]model.Port{
// 			"inputs": {
// 				{Name: "in", Label: "输入", PortType: "connection"},
// 			},
// 			"outputs": {
// 				{Name: "out", Label: "输出", PortType: "connection"},
// 			},
// 		},
// 	}

// 	// 2. 创建连接类型
// 	connType := model.ConnectionType{
// 		ID:               1,
// 		Name:             "test_conn",
// 		Description:      "测试连接",
// 		Color:            "#FF0000",
// 		AllowedPortTypes: []string{"connection"},
// 	}

// 	// 3. 创建工作流
// 	wf := &model.Workflow{
// 		ID: "test_workflow",
// 		Dag: &model.Dag{
// 			Nodes: map[string]*model.Node{
// 				"node1": {
// 					ID:     "node1",
// 					TypeID: 1,
// 					Label:  "节点1",
// 				},
// 				"node2": {
// 					ID:     "node2",
// 					TypeID: 1,
// 					Label:  "节点2",
// 				},
// 			},
// 			Connections: []model.Connection{
// 				{
// 					ID:     "conn1",
// 					TypeID: 1,
// 					From: model.Endpoint{
// 						NodeID:   "node1",
// 						PortName: "out",
// 					},
// 					To: model.Endpoint{
// 						NodeID:   "node2",
// 						PortName: "in",
// 					},
// 				},
// 			},
// 		},
// 		NodeTypes: map[int]model.NodeType{
// 			1: nodeType,
// 		},
// 		ConnectionTypes: map[int]model.ConnectionType{
// 			1: connType,
// 		},
// 	}

// 	// 4. 验证工作流
// 	if err := wf.Validate(); err != nil {
// 		log.Fatalf("工作流验证失败: %v", err)
// 	}

// 	// 5. 创建执行上下文
// 	ctx := &model.ExecutionContext{
// 		Workflow: wf,
// 		Logger: func(msg string) {
// 			fmt.Println(msg)
// 		},
// 		Vars: make(map[string]interface{}),
// 	}

// 	// 6. 执行工作流
// 	if err := wf.ExecuteWorkflow(ctx); err != nil {
// 		log.Fatalf("工作流执行失败: %v", err)
// 	}

// 	// 7. 打印工作流状态
// 	wfJSON, _ := json.MarshalIndent(wf, "", "  ")
// 	fmt.Printf("\n工作流状态:\n%s\n", string(wfJSON))
// }
