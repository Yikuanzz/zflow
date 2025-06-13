package data

import "zflow/internal/model"

// 定义常用的节点类型
var (
	// 基础节点类型
	EchoNode = model.NodeType{
		ID:       1,
		Category: "basic",
		Note:     "回显节点，用于测试",
		Operation: &EchoOperation{
			Message: "Hello from Echo Node!",
		},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "in", Label: "输入", PortType: "connection"},
			},
			"outputs": {
				{Name: "out", Label: "输出", PortType: "connection"},
			},
		},
	}

	// 文件操作节点类型
	FileReadNode = model.NodeType{
		ID:       2,
		Category: "file",
		Note:     "文件读取节点",
		Operation: &FileReadOperation{
			Encoding: "utf-8",
		},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "file_path", Label: "文件路径", PortType: "file"},
			},
			"outputs": {
				{Name: "content", Label: "文件内容", PortType: "connection"},
			},
		},
	}

	FileWriteNode = model.NodeType{
		ID:       3,
		Category: "file",
		Note:     "文件写入节点",
		Operation: &FileWriteOperation{
			Encoding: "utf-8",
		},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "content", Label: "文件内容", PortType: "connection"},
				{Name: "file_path", Label: "文件路径", PortType: "file"},
			},
			"outputs": {
				{Name: "success", Label: "写入结果", PortType: "connection"},
			},
		},
	}

	// 数据转换节点类型
	TransformNode = model.NodeType{
		ID:       4,
		Category: "transform",
		Note:     "数据转换节点",
		Operation: &TransformOperation{
			Template: "{{.input}}",
		},
		Properties: map[string][]model.Port{
			"inputs": {
				{Name: "in", Label: "输入数据", PortType: "connection"},
			},
			"outputs": {
				{Name: "out", Label: "转换结果", PortType: "connection"},
			},
		},
	}
)

// GetDefaultNodeTypes 返回所有默认的节点类型
func GetDefaultNodeTypes() map[int]model.NodeType {
	return map[int]model.NodeType{
		EchoNode.ID:      EchoNode,
		FileReadNode.ID:  FileReadNode,
		FileWriteNode.ID: FileWriteNode,
		TransformNode.ID: TransformNode,
	}
}
