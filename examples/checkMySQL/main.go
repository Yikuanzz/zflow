package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"zflow/internal/model"
)

type CheckMySQLOperation struct{}

func (op *CheckMySQLOperation) Execute(ctx model.Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error) {
	// 获取输入参数
	ip, exists := inputs["ip"]
	if !exists {
		return nil, fmt.Errorf("未找到输入IP地址")
	}

	port, exists := inputs["port"]
	if !exists {
		port = []byte("3306") // 默认端口
	}

	timeout, exists := inputs["timeout"]
	if !exists {
		timeout = []byte("3") // 默认超时时间
	}

	// 解析超时时间
	timeoutSec, err := strconv.Atoi(string(timeout))
	if err != nil {
		return nil, fmt.Errorf("无效的超时时间: %v", err)
	}

	// 构建地址
	address := fmt.Sprintf("%s:%s", string(ip), string(port))
	ctx.Log(fmt.Sprintf("正在检查 MySQL 服务: %s (超时: %d秒)", address, timeoutSec))

	// 尝试连接
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeoutSec)*time.Second)
	if err != nil {
		ctx.Log(fmt.Sprintf("MySQL 服务检查失败: %v", err))
		return map[string][]byte{
			"status":  []byte("failed"),
			"message": []byte("未检测到 MySQL 服务"),
			"error":   []byte(err.Error()),
		}, nil
	}
	defer conn.Close()

	ctx.Log("MySQL 服务检查成功")
	return map[string][]byte{
		"status":  []byte("success"),
		"message": []byte("存在 MySQL 服务"),
		"address": []byte(address),
	}, nil
}

type echoMsgOperations struct{}

func (op *echoMsgOperations) Execute(ctx model.Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error) {
	// 如果没有输入，返回空输出
	if len(inputs) == 0 {
		return map[string][]byte{}, nil
	}

	// 回显所有输入
	outputs := make(map[string][]byte)
	for name, data := range inputs {
		ctx.Log(fmt.Sprintf("回显消息 [%s]: %s", name, string(data)))
		outputs[name] = data
	}
	return outputs, nil
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
				{Name: "ip", Label: "机器IP地址", PortType: "connection"},
				{Name: "port", Label: "端口号", PortType: "connection"},
				{Name: "timeout", Label: "超时时间(秒)", PortType: "connection"},
			},
			"outputs": {
				{Name: "status", Label: "状态", PortType: "connection"},
				{Name: "message", Label: "消息", PortType: "connection"},
				{Name: "error", Label: "错误信息", PortType: "connection"},
				{Name: "address", Label: "地址", PortType: "connection"},
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
			"outputs": {
				{Name: "out", Label: "输出", PortType: "connection"},
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
				"node1": model.NewNode("node1", 1, "MySQL检查节点"),
				"node2": model.NewNode("node2", 2, "成功回显节点"),
				"node3": model.NewNode("node3", 2, "失败回显节点"),
			},
			Connections: []model.Connection{
				{
					ID:     "conn1",
					TypeID: 1,
					From: model.Endpoint{
						NodeID:   "node1",
						PortName: "message",
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
						PortName: "error",
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

	// 设置初始输入数据
	wf.Dag.Nodes["node1"].Inputs["ip"] = []byte("127.0.0.1")
	wf.Dag.Nodes["node1"].Inputs["port"] = []byte("3306")
	wf.Dag.Nodes["node1"].Inputs["timeout"] = []byte("3")

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
	fmt.Println("\n工作流执行完成！")
	for nodeID, node := range wf.Dag.Nodes {
		fmt.Printf("节点 %s 状态: %s\n", nodeID, node.State)
		if len(node.Outputs) > 0 {
			fmt.Printf("输出数据:\n")
			for port, data := range node.Outputs {
				fmt.Printf("  %s: %s\n", port, string(data))
			}
		}
	}
}
