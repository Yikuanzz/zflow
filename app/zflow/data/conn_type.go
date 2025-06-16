package data

import "zflow/app/zflow/model"

// 定义常用的连接类型
var (
	// 数据流连接 - 用于传递普通数据
	DataFlowConn = model.ConnectionType{
		UID:              "1",
		Name:             "data_flow",
		Description:      "数据流连接，用于传递普通数据",
		Color:            "#4CAF50", // 绿色
		AllowedPortTypes: []string{"connection"},
	}

	// 文件流连接 - 用于传递文件路径
	FileFlowConn = model.ConnectionType{
		UID:              "2",
		Name:             "file_flow",
		Description:      "文件流连接，用于传递文件路径",
		Color:            "#2196F3", // 蓝色
		AllowedPortTypes: []string{"file"},
	}

	// 错误流连接 - 用于传递错误信息
	ErrorFlowConn = model.ConnectionType{
		UID:              "3",
		Name:             "error_flow",
		Description:      "错误流连接，用于传递错误信息",
		Color:            "#F44336", // 红色
		AllowedPortTypes: []string{"connection"},
	}

	// 控制流连接 - 用于控制流程
	ControlFlowConn = model.ConnectionType{
		UID:              "4",
		Name:             "control_flow",
		Description:      "控制流连接，用于控制流程",
		Color:            "#FF9800", // 橙色
		AllowedPortTypes: []string{"connection"},
	}
)

// GetDefaultConnectionTypes 返回所有默认的连接类型
func GetDefaultConnectionTypes() map[string]model.ConnectionType {
	return map[string]model.ConnectionType{
		DataFlowConn.UID: DataFlowConn,
		// FileFlowConn.UID:    FileFlowConn,
		// ErrorFlowConn.UID:   ErrorFlowConn,
		// ControlFlowConn.UID: ControlFlowConn,
	}
}
