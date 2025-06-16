package core

import (
	"fmt"
	"zflow/app/zflow/model"
)

// NoteTypeList 节点类型列表
var NoteTypeList = []model.NodeType{
	AddNodeType,
	MulNodeType,
	EchoNodeType,
}

// 加法节点
var AddNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.add.v1", ServiceName),
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

// 乘法节点
var MulNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.mul.v1", ServiceName),
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

// 回显节点
var EchoNodeType = model.NodeType{
	UID:       fmt.Sprintf("%s.echo.v1", ServiceName),
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
