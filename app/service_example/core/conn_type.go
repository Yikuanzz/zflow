package core

import "zflow/app/zflow/model"

// ConnTypes 定义常用的连接类型
var ConnTypes = map[string]*model.ConnectionType{
	"data_flow": &DataFlowConn,
}

// DataFlowConn 数据流连接 - 用于传递普通数据
var DataFlowConn = model.ConnectionType{
	UID:              "1",
	Name:             "data_flow",
	Description:      "数据流连接，用于传递普通数据",
	Color:            "#4CAF50", // 绿色
	AllowedPortTypes: []string{"connection"},
}
