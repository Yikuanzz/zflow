package core

import (
	"fmt"
	"zflow/app/zflow/model"
)

// 默认连接类型
var DefaultConnType = model.ConnectionType{
	UID:              fmt.Sprintf("%s.default.v1", ServiceName),
	Name:             "默认连接",
	Description:      "用于节点间数据传输的默认连接类型",
	Color:            "#666666",
	AllowedPortTypes: []string{"connection"},
}
