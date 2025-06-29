package core

import (
	"fmt"

	"zflow/app/zflow/model"
)

// 节点UID
var (
	AddNodeTypeUID  = "add"
	MulNodeTypeUID  = "mul"
	EchoNodeTypeUID = "echo"
)

// NodeTypes 定义节点类型
var NodeTypes = map[string]*model.NodeType{
	"add":  &AddNodeType,
	"mul":  &MulNodeType,
	"echo": &EchoNodeType,
}

// AddNodeType 加法节点
var AddNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.add", ServiceName),
	Category:  "math",
	Note:      "两个数字相加，输出结果",
	Operation: AddOperationInst, // 这里你可以填入具体的加法 Operation 实例
	Properties: map[string][]model.Port{
		"inputs": {
			{Name: "a", Label: "加数A", PortType: "connection"},
			{Name: "b", Label: "加数B", PortType: "connection"},
		},
		"outputs": {
			{Name: "sum", Label: "和", PortType: "connection"},
		},
	},
}

// MulNodeType 乘法节点
var MulNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.mul", ServiceName),
	Category:  "math",
	Note:      "两个数字相乘，输出结果",
	Operation: MulOperationInst, // 这里你可以填入具体的乘法 Operation 实例
	Properties: map[string][]model.Port{
		"inputs": {
			{Name: "a", Label: "乘数A", PortType: "connection"},
			{Name: "b", Label: "乘数B", PortType: "connection"},
		},
		"outputs": {
			{Name: "product", Label: "积", PortType: "connection"},
		},
	},
}

// EchoNodeType 回显节点
var EchoNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.echo", ServiceName),
	Category:  "util",
	Note:      "回显输入内容，常用于调试或展示节点计算结果",
	Operation: EchoOperationInst, // 这里你可以填入具体的 Echo Operation 实例
	Properties: map[string][]model.Port{
		"inputs": {
			{Name: "input", Label: "输入内容", PortType: "connection"},
		},
		"outputs": {
			{Name: "output", Label: "输出内容", PortType: "connection"},
		},
	},
}
